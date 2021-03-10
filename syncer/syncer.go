package syncer

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"manala/logger"
	"manala/models"
	"manala/template"
	"os"
	"path/filepath"
	"regexp"
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

/**********/
/* Syncer */
/**********/

func New(log *logger.Logger, tmpl *template.Template) *Syncer {
	return &Syncer{
		log:  log,
		tmpl: tmpl,
	}
}

type Syncer struct {
	log  *logger.Logger
	tmpl *template.Template
}

// Sync a project from its recipe
func (snc *Syncer) SyncProject(prj models.ProjectInterface) error {
	// Include template helpers if any
	helpers := filepath.Join(prj.Recipe().Dir(), "_helpers.tmpl")
	if _, err := os.Stat(helpers); err == nil {
		if err := snc.tmpl.ParseFiles(helpers); err == nil {
			return err
		}
	}

	for _, sync := range prj.Recipe().SyncUnits() {
		if err := snc.Sync(
			filepath.Join(prj.Recipe().Dir(), sync.Source),
			filepath.Join(prj.Dir(), sync.Destination),
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
func (snc *Syncer) Sync(src string, dst string, ctx interface{}) error {
	node, err := newNode(src, dst, snc.tmpl, ctx)
	if err != nil {
		return err
	}

	if err := snc.syncNode(node); err != nil {
		return err
	}

	return nil
}

func (snc *Syncer) syncNode(node *node) error {
	if node.Src.IsDir {

		snc.log.DebugWithFields("Syncing directory...", logger.Fields{
			"src": node.Src.Path,
			"dst": node.Dst.Path,
		})

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

			snc.log.InfoWithField("Synced directory", "path", node.Dst.Path)
		}

		// Iterate over source files
		// Make a map of destination files map for quick lookup; used in deletion below
		dstMap := make(map[string]bool)
		for _, file := range node.Src.Files {
			fileNode, err := newNode(
				filepath.Join(node.Src.Path, file),
				filepath.Join(node.Dst.Path, file),
				node.Template,
				&node.Context,
			)
			if err != nil {
				return err
			}

			dstMap[filepath.Base(fileNode.Dst.Path)] = true

			if err := snc.syncNode(fileNode); err != nil {
				return err
			}
		}

		// Delete not synced destination files
		files, err := os.ReadDir(node.Dst.Path)
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

		snc.log.DebugWithFields("Syncing file...", logger.Fields{
			"src": node.Src.Path,
			"dst": node.Dst.Path,
		})

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
			tmplContent, err := os.ReadFile(node.Src.Path)
			if err != nil {
				return err
			}
			// Parse
			if err := node.Template.Parse(string(tmplContent)); err != nil {
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

			snc.log.InfoWithField("Synced file", "path", node.Dst.Path)
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
		files, err := os.ReadDir(node.Src.Path)
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
			files, err := os.ReadDir(node.Dst.Path)
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
