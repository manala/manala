package syncer

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/apex/log"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"manala/models"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

/**********/
/* Errors */
/**********/

type SourceNotExistError struct {
	Source string
}

func (e *SourceNotExistError) Error() string {
	return "no source " + e.Source + " file or directory "
}

/********/
/* Sync */
/********/

// Sync a project from a recipe
func SyncProject(prj models.ProjectInterface) error {
	// Template
	tmpl := NewTemplate()

	// Include helpers if any
	helpers := path.Join(prj.Recipe().Dir(), "_helpers.tmpl")
	if _, err := os.Stat(helpers); err == nil {
		_, err = tmpl.ParseFiles(helpers)
		if err != nil {
			return err
		}
	}

	for _, sync := range prj.Recipe().SyncUnits() {
		if err := Sync(
			path.Join(prj.Recipe().Dir(), sync.Source),
			path.Join(prj.Dir(), sync.Destination),
			tmpl,
			map[string]interface{}{
				"Vars": prj.Vars(),
			},
		); err != nil {
			return err
		}
	}

	return nil
}

// Sync a source with a destination
func Sync(src string, dst string, tmpl *template.Template, ctx interface{}) error {
	node, err := newNode(src, dst, tmpl, ctx)
	if err != nil {
		return err
	}

	if err := syncNode(node); err != nil {
		return err
	}

	return nil
}

type node struct {
	Src struct {
		Path         string
		IsDir        bool
		Files        []string
		IsExecutable bool
	}
	IsDist bool
	IsTmpl bool
	Dst    struct {
		Path    string
		Mode    os.FileMode
		Hash    []byte
		IsExist bool
		IsDir   bool
		Files   []string
	}
	Template *template.Template
	Context  interface{}
}

var distRegex = regexp.MustCompile(`(\.dist)(?:$|\.tmpl$)`)
var tmplRegex = regexp.MustCompile(`(\.tmpl)(?:$|\.dist$)`)

func newNode(src string, dst string, tmpl *template.Template, cxt interface{}) (*node, error) {
	node := &node{}
	node.Src.Path = src
	node.Dst.Path = dst
	node.Template = tmpl
	node.Context = cxt

	// Source info
	stat, err := os.Stat(node.Src.Path)
	if err != nil {
		// Source does not exist
		if os.IsNotExist(err) {
			return nil, &SourceNotExistError{node.Src.Path}
		} else {
			return nil, err
		}
	}
	node.Src.IsDir = stat.IsDir()

	if node.Src.IsDir {
		files, err := ioutil.ReadDir(node.Src.Path)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			node.Src.Files = append(node.Src.Files, file.Name())
		}
	} else {
		node.Src.IsExecutable = (stat.Mode() & 0100) != 0

		if distRegex.MatchString(node.Src.Path) {
			node.IsDist = true
			node.Dst.Path = distRegex.ReplaceAllString(node.Dst.Path, "")
		}

		if tmplRegex.MatchString(node.Src.Path) {
			node.IsTmpl = true
			node.Dst.Path = tmplRegex.ReplaceAllString(node.Dst.Path, "")
		}
	}

	// Destination info
	stat, err = os.Stat(node.Dst.Path)
	if err != nil {
		// Error other than not existing destination
		if !os.IsNotExist(err) {
			return nil, err
		}
		node.Dst.IsExist = false
	} else {
		node.Dst.IsExist = true
		node.Dst.IsDir = stat.IsDir()
	}

	if node.Dst.IsExist {
		// Mode
		node.Dst.Mode = stat.Mode()

		if node.Dst.IsDir {
			files, err := ioutil.ReadDir(node.Dst.Path)
			if err != nil {
				return nil, err
			}
			for _, file := range files {
				node.Dst.Files = append(node.Dst.Files, file.Name())
			}
		} else {
			// Get destination hash
			file, err := os.Open(node.Dst.Path)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			hash := md5.New()
			if _, err := io.Copy(hash, file); err != nil {
				return nil, err
			}

			node.Dst.Hash = hash.Sum(nil)
		}
	}

	return node, nil
}

func syncNode(node *node) error {
	if node.Src.IsDir {

		log.WithFields(log.Fields{
			"src": node.Src.Path,
			"dst": node.Dst.Path,
		}).Debug("Syncing directory...")

		// Destination is a file; remove
		if node.Dst.IsExist && !node.Dst.IsDir {
			if err := os.Remove(node.Dst.Path); err != nil {
				return err
			}
			node.Dst.IsExist = false
		}

		// Destination does not exists; create
		if !node.Dst.IsExist {
			if err := os.MkdirAll(node.Dst.Path, 0755); err != nil {
				return err
			}

			log.WithFields(log.Fields{
				"path": node.Dst.Path,
			}).Info("Synced directory")
		}

		// Iterate over source files
		// Make a map of destination files map for quick lookup; used in deletion below
		dstMap := make(map[string]bool)
		for _, file := range node.Src.Files {
			fileNode, err := newNode(
				path.Join(node.Src.Path, file),
				path.Join(node.Dst.Path, file),
				node.Template,
				&node.Context,
			)
			if err != nil {
				return err
			}

			dstMap[filepath.Base(fileNode.Dst.Path)] = true

			if err := syncNode(fileNode); err != nil {
				return err
			}
		}

		// Delete not synced destination files
		files, err := ioutil.ReadDir(node.Dst.Path)
		if err != nil {
			return err
		}

		for _, file := range files {
			if !dstMap[file.Name()] {
				if err := os.RemoveAll(filepath.Join(node.Dst.Path, file.Name())); err != nil {
					return err
				}
			}
		}

		return nil

	} else {

		log.WithFields(log.Fields{
			"src": node.Src.Path,
			"dst": node.Dst.Path,
		}).Debug("Syncing file...")

		// Destination is a directory; remove
		if node.Dst.IsExist && node.Dst.IsDir {
			if err := os.RemoveAll(node.Dst.Path); err != nil {
				return err
			}
			node.Dst.IsExist = false
			node.Dst.IsDir = false
		}

		// Node is a dist and destination already exists (or was a directory); exit
		if node.Dst.IsExist && node.IsDist {
			return nil
		}

		equal := false

		var srcReader io.Reader

		if node.IsTmpl {
			// Read template content
			tmplContent, err := ioutil.ReadFile(node.Src.Path)
			if err != nil {
				return err
			}
			// Parse
			_, err = node.Template.Parse(string(tmplContent))
			if err != nil {
				return err
			}
			// Execute
			var buffer bytes.Buffer
			if err := node.Template.Execute(&buffer, node.Context); err != nil {
				return fmt.Errorf("invalid template \"%s\" (%s)", node.Src.Path, err)
			}

			srcReader = bytes.NewReader(buffer.Bytes())

			if node.Dst.IsExist {
				// Get template hash
				hash := md5.New()
				if _, err := io.Copy(hash, &buffer); err != nil {
					return err
				}
				equal = bytes.Compare(hash.Sum(nil), node.Dst.Hash) == 0
			}
		} else {
			// Node is not a template, let's go buffering \o/
			srcFile, err := os.Open(node.Src.Path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			if node.Dst.IsExist {
				// Get source hash
				hash := md5.New()
				if _, err := io.Copy(hash, srcFile); err != nil {
					return err
				}
				equal = bytes.Compare(hash.Sum(nil), node.Dst.Hash) == 0

				if _, err := srcFile.Seek(0, io.SeekStart); err != nil {
					return err
				}
			}

			srcReader = srcFile
		}

		// Files are not equals or destination does not exists
		if !equal {
			// Destination file mode
			var dstMode os.FileMode = 0666
			if node.Src.IsExecutable {
				dstMode = 0777
			}

			// Create or truncate destination file
			dstFile, err := os.OpenFile(node.Dst.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, dstMode)
			if err != nil {
				return err
			}
			defer dstFile.Close()

			// Copy from source to destination
			_, err = io.Copy(dstFile, srcReader)
			if err != nil {
				return err
			}

			log.WithFields(log.Fields{
				"path": node.Dst.Path,
			}).Info("Synced file")
		} else {
			dstMode := node.Dst.Mode &^ 0111
			if node.Src.IsExecutable {
				dstMode = node.Dst.Mode | 0111
			}

			if dstMode != node.Dst.Mode {
				if err := os.Chmod(node.Dst.Path, dstMode); err != nil {
					return err
				}
			}
		}

		return nil
	}
}

/************/
/* Template */
/************/

func NewTemplate() *template.Template {
	tmpl := template.New("")

	// Execution stops immediately with an error.
	tmpl.Option("missingkey=error")

	tmpl.Funcs(sprig.TxtFuncMap())
	tmpl.Funcs(template.FuncMap{
		"toYaml":  templateToYamlFunc(),
		"include": templateIncludeFunc(tmpl),
	})

	return tmpl
}

// As seen in helm
func templateToYamlFunc() func(value interface{}) string {
	return func(value interface{}) string {
		var buf bytes.Buffer

		enc := yaml.NewEncoder(&buf)

		if err := enc.Encode(value); err != nil {
			// Swallow errors inside of a template.
			return ""
		}

		return strings.TrimSuffix(buf.String(), "\n")
	}
}

// As seen in helm
func templateIncludeFunc(tmpl *template.Template) func(name string, data interface{}) (string, error) {
	includedNames := make([]string, 0)
	return func(name string, data interface{}) (string, error) {
		var buf strings.Builder
		includedCount := 0
		for _, n := range includedNames {
			if n == name {
				includedCount += 1
			}
		}
		if includedCount >= 16 {
			return "", fmt.Errorf("rendering template has reached the maximum nested reference name level: %s", name)
		}
		includedNames = append(includedNames, name)
		err := tmpl.ExecuteTemplate(&buf, name, data)
		includedNames = includedNames[:len(includedNames)-1]
		return buf.String(), err
	}
}
