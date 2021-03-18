package syncer

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"manala/fs"
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

func New(log *logger.Logger, fsManager models.FsManagerInterface, templateManager models.TemplateManagerInterface) *Syncer {
	return &Syncer{
		log:             log,
		fsManager:       fsManager,
		templateManager: templateManager,
	}
}

type Syncer struct {
	log             *logger.Logger
	fsManager       models.FsManagerInterface
	templateManager models.TemplateManagerInterface
}

// Sync a project from its recipe
func (snc *Syncer) SyncProject(prj models.ProjectInterface) error {
	// Recipe template
	recTmpl, err := snc.templateManager.NewRecipeTemplate(prj.Recipe())
	if err != nil {
		return err
	}

	// Loop over sync nodes
	for _, node := range prj.Recipe().Sync() {
		if err := snc.Sync(
			snc.fsManager.NewModelFs(prj.Recipe()),
			node.Source,
			recTmpl,
			node.Destination,
			snc.fsManager.NewModelFs(prj),
			prj.Vars(),
		); err != nil {
			return err
		}
	}

	return nil
}

// Sync a source with a destination
func (snc *Syncer) Sync(srcFs fs.ReadInterface, src string, srcTmpl template.Interface, dst string, dstFs fs.ReadWriteInterface, vars map[string]interface{}) error {
	node, err := newNode(srcFs, src, dstFs, dst, srcTmpl, vars)
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
			if err := node.Dst.Fs.Remove(node.Dst.Path); err != nil {
				return err
			}
			node.Dst.IsExist = false
		}

		// Destination does not exists; create
		if !node.Dst.IsExist {
			if err := node.Dst.Fs.MkdirAll(node.Dst.Path, 0755); err != nil {
				return err
			}

			snc.log.InfoWithField("Synced directory", "path", node.Dst.Path)
		}

		// Iterate over source files
		// Make a map of destination files map for quick lookup; used in deletion below
		dstMap := make(map[string]bool)
		for _, file := range node.Src.Files {
			fileNode, err := newNode(
				node.Src.Fs,
				filepath.Join(node.Src.Path, file),
				node.Dst.Fs,
				filepath.Join(node.Dst.Path, file),
				node.Src.Template,
				node.Vars,
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
		files, err := node.Dst.Fs.ReadDir(node.Dst.Path)
		if err != nil {
			return err
		}

		for _, file := range files {
			if !dstMap[file.Name()] {
				if err := node.Dst.Fs.RemoveAll(filepath.Join(node.Dst.Path, file.Name())); err != nil {
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
			if err := node.Dst.Fs.RemoveAll(node.Dst.Path); err != nil {
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
			// Parse
			if err := node.Src.Template.ParseFile(node.Src.Path); err != nil {
				return err
			}

			// Execute
			var buffer bytes.Buffer
			if err := node.Src.Template.Execute(&buffer, node.Vars); err != nil {
				return fmt.Errorf("invalid template \"%s\" (%s)", node.Src.Path, err)
			}

			srcReader = bytes.NewReader(buffer.Bytes())

			if node.Dst.IsExist {
				// Get template hash
				hash := sha1.New()
				if _, err := io.Copy(hash, &buffer); err != nil {
					return err
				}
				equal = bytes.Compare(hash.Sum(nil), node.Dst.Hash) == 0
			}
		} else {
			// Node is not a template, let's go buffering \o/
			srcFile, err := node.Src.Fs.Open(node.Src.Path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			if node.Dst.IsExist {
				// Open a new source file to not interfere with previous one read position,
				// as fs.File interface does not provide a Seek method
				srcFileHash, err := node.Src.Fs.Open(node.Src.Path)
				if err != nil {
					return err
				}
				defer srcFileHash.Close()

				// Get source hash
				hash := sha1.New()
				if _, err := io.Copy(hash, srcFileHash); err != nil {
					return err
				}
				equal = bytes.Compare(hash.Sum(nil), node.Dst.Hash) == 0
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
			dstFile, err := node.Dst.Fs.OpenFile(node.Dst.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, dstMode)
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
				if err := node.Dst.Fs.Chmod(node.Dst.Path, dstMode); err != nil {
					return err
				}
			}
		}

		return nil
	}
}

type node struct {
	Src struct {
		Fs           fs.ReadInterface
		Template     template.Interface
		Path         string
		IsDir        bool
		Files        []string
		IsExecutable bool
	}
	IsDist bool
	IsTmpl bool
	Dst    struct {
		Fs      fs.ReadWriteInterface
		Path    string
		Mode    os.FileMode
		Hash    []byte
		IsExist bool
		IsDir   bool
		Files   []string
	}
	Vars map[string]interface{}
}

var distRegex = regexp.MustCompile(`(\.dist)(?:$|\.tmpl$)`)
var tmplRegex = regexp.MustCompile(`(\.tmpl)(?:$|\.dist$)`)

func newNode(srcFs fs.ReadInterface, src string, dstFs fs.ReadWriteInterface, dst string, srcTmpl template.Interface, vars map[string]interface{}) (*node, error) {
	node := &node{}
	node.Src.Fs = srcFs
	node.Src.Template = srcTmpl
	node.Dst.Fs = dstFs
	node.Vars = vars

	// Source stat
	srcStat, err := srcFs.Stat(src)
	if err != nil {
		// Source does not exist
		if errors.Is(err, os.ErrNotExist) {
			return nil, &SourceNotExistError{src}
		} else {
			return nil, err
		}
	}
	node.Src.IsDir = srcStat.IsDir()

	if node.Src.IsDir {
		files, err := srcFs.ReadDir(src)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			node.Src.Files = append(node.Src.Files, file.Name())
		}
	} else {
		node.Src.IsExecutable = (srcStat.Mode() & 0100) != 0

		if distRegex.MatchString(src) {
			node.IsDist = true
			dst = distRegex.ReplaceAllString(dst, "")
		}

		if tmplRegex.MatchString(src) {
			node.IsTmpl = true
			dst = tmplRegex.ReplaceAllString(dst, "")
		}
	}

	// Destination stat
	dstStat, err := dstFs.Stat(dst)
	if err != nil {
		// Error other than not existing destination
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		node.Dst.IsExist = false
	} else {
		node.Dst.IsExist = true
		node.Dst.IsDir = dstStat.IsDir()
	}

	if node.Dst.IsExist {
		// Mode
		node.Dst.Mode = dstStat.Mode()

		if node.Dst.IsDir {
			files, err := dstFs.ReadDir(dst)
			if err != nil {
				return nil, err
			}
			for _, file := range files {
				node.Dst.Files = append(node.Dst.Files, file.Name())
			}
		} else {
			// Get destination hash
			file, err := dstFs.Open(dst)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			hash := sha1.New()
			if _, err := io.Copy(hash, file); err != nil {
				return nil, err
			}

			node.Dst.Hash = hash.Sum(nil)
		}
	}

	node.Src.Path = src
	node.Dst.Path = dst

	return node, nil
}
