package syncer

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"github.com/apex/log"
	"io"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	internalTemplate "manala/internal/template"
	"os"
	"path/filepath"
	"regexp"
)

type Syncer struct {
	Log *internalLog.Logger
}

// Sync a source with a destination
func (syncer *Syncer) Sync(
	srcDir string,
	src string,
	dstDir string,
	dst string,
	templateProvider internalTemplate.ProviderInterface,
) error {
	node, err := newNode(srcDir, src, dstDir, dst, templateProvider)
	if err != nil {
		return err
	}

	if err := syncer.syncNode(node); err != nil {
		return err
	}

	return nil
}

func (syncer *Syncer) syncNode(node *node) error {
	relSrcPath, _ := filepath.Rel(node.Src.Dir, node.Src.Path)
	relDstPath, _ := filepath.Rel(node.Dst.Dir, node.Dst.Path)
	if node.Src.IsDir {

		syncer.Log.WithFields(log.Fields{
			"src": relSrcPath,
			"dst": relDstPath,
		}).Debug("sync dir")

		// Destination is a file; remove
		if node.Dst.IsExist && !node.Dst.IsDir {
			if err := os.Remove(node.Dst.Path); err != nil {
				return internalOs.FileSystemError(err)
			}
			node.Dst.IsExist = false
		}

		// Destination does not exist; create
		if !node.Dst.IsExist {
			if err := os.MkdirAll(node.Dst.Path, 0755); err != nil {
				return internalOs.FileSystemError(err)
			}

			syncer.Log.WithField(
				"path", relDstPath,
			).Info("dir synced")
		}

		// Iterate over source files
		// Make a map of destination files map for quick lookup; used in deletion below
		dstMap := make(map[string]bool)
		for _, file := range node.Src.Files {
			fileNode, err := newNode(
				node.Src.Dir,
				filepath.Join(relSrcPath, file),
				node.Dst.Dir,
				filepath.Join(relDstPath, file),
				node.TemplateProvider,
			)
			if err != nil {
				return err
			}

			dstMap[filepath.Base(fileNode.Dst.Path)] = true

			if err := syncer.syncNode(fileNode); err != nil {
				return err
			}
		}

		// Delete not synced destination files
		files, err := os.ReadDir(node.Dst.Path)
		if err != nil {
			return internalOs.FileSystemError(err)
		}

		for _, file := range files {
			if !dstMap[file.Name()] {
				if err := os.RemoveAll(filepath.Join(node.Dst.Path, file.Name())); err != nil {
					return internalOs.FileSystemError(err)
				}
			}
		}

		return nil

	} else {

		syncer.Log.WithFields(log.Fields{
			"src": relSrcPath,
			"dst": relDstPath,
		}).Debug("sync file")

		if node.Dst.IsExist {
			// Destination is a directory; remove
			if node.Dst.IsDir {
				if err := os.RemoveAll(node.Dst.Path); err != nil {
					return internalOs.FileSystemError(err)
				}
				node.Dst.IsExist = false
				node.Dst.IsDir = false
			}
			// Node is a dist and destination already exists (or was a directory); exit
			if node.IsDist {
				return nil
			}
		} else {
			// Ensure destination parents directories exists
			if dir := filepath.Dir(node.Dst.Path); dir != "." {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return internalOs.FileSystemError(err)
				}
			}
		}

		equal := false

		var srcReader io.Reader

		if node.IsTmpl {
			// Write template
			var buffer bytes.Buffer
			if err := node.TemplateProvider.Template().WithFile(node.Src.Path).Write(&buffer); err != nil {
				return err
			}

			srcReader = bytes.NewReader(buffer.Bytes())

			if node.Dst.IsExist {
				// Get template hash
				hash := sha1.New()
				if _, err := io.Copy(hash, &buffer); err != nil {
					return err
				}
				equal = bytes.Equal(hash.Sum(nil), node.Dst.Hash)
			}
		} else {
			// Node is not a template, let's go buffering \o/
			srcFile, err := os.Open(node.Src.Path)
			if err != nil {
				return internalOs.FileSystemError(err)
			}
			defer srcFile.Close()

			if node.Dst.IsExist {
				// Get source hash
				hash := sha1.New()
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
				return internalOs.FileSystemError(err)
			}
			defer dstFile.Close()

			// Copy from source to destination
			_, err = io.Copy(dstFile, srcReader)
			if err != nil {
				return err
			}

			syncer.Log.WithField(
				"path", relDstPath,
			).Info("file synced")
		} else {
			dstMode := node.Dst.Mode &^ 0111
			if node.Src.IsExecutable {
				dstMode = node.Dst.Mode | 0111
			}

			if dstMode != node.Dst.Mode {
				if err := os.Chmod(node.Dst.Path, dstMode); err != nil {
					return internalOs.FileSystemError(err)
				}
			}
		}

		return nil
	}
}

type node struct {
	Src struct {
		Dir          string
		Path         string
		IsDir        bool
		Files        []string
		IsExecutable bool
	}
	IsDist bool
	IsTmpl bool
	Dst    struct {
		Dir     string
		Path    string
		Mode    os.FileMode
		Hash    []byte
		IsExist bool
		IsDir   bool
		Files   []string
	}
	TemplateProvider internalTemplate.ProviderInterface
}

var distRegex = regexp.MustCompile(`(\.dist)(?:$|\.tmpl$)`)
var tmplRegex = regexp.MustCompile(`(\.tmpl)(?:$|\.dist$)`)

func newNode(srcDir string, src string, dstDir string, dst string, templateProvider internalTemplate.ProviderInterface) (*node, error) {
	node := &node{}
	node.Src.Dir = srcDir
	node.Dst.Dir = dstDir
	node.TemplateProvider = templateProvider

	srcPath := filepath.Join(node.Src.Dir, src)

	// Source stat
	srcStat, err := os.Stat(srcPath)
	if err != nil {
		// Source does not exist
		if errors.Is(err, os.ErrNotExist) {
			return nil, SourceNotExistError(srcPath)
		} else {
			return nil, internalOs.FileSystemError(err)
		}
	}
	node.Src.IsDir = srcStat.IsDir()

	if node.Src.IsDir {
		files, err := os.ReadDir(srcPath)
		if err != nil {
			return nil, internalOs.FileSystemError(err)
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

	dstPath := filepath.Join(node.Dst.Dir, dst)

	// Destination stat
	dstStat, err := os.Stat(dstPath)
	if err != nil {
		// Error other than not existing destination
		if !errors.Is(err, os.ErrNotExist) {
			return nil, internalOs.FileSystemError(err)
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
			files, err := os.ReadDir(dstPath)
			if err != nil {
				return nil, internalOs.FileSystemError(err)
			}
			for _, file := range files {
				node.Dst.Files = append(node.Dst.Files, file.Name())
			}
		} else {
			// Get destination hash
			file, err := os.Open(dstPath)
			if err != nil {
				return nil, internalOs.FileSystemError(err)
			}
			defer file.Close()

			hash := sha1.New()
			if _, err := io.Copy(hash, file); err != nil {
				return nil, err
			}

			node.Dst.Hash = hash.Sum(nil)
		}
	}

	node.Src.Path = srcPath
	node.Dst.Path = dstPath

	return node, nil
}
