package sync

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/errors/std"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/template/engine"
)

type Syncer struct {
	log *log.Log
}

func NewSyncer(log *log.Log) *Syncer {
	return &Syncer{
		log: log,
	}
}

// Sync a source with a destination.
func (syncer *Syncer) Sync(
	srcDir string,
	src string,
	dstDir string,
	dst string,
	templateExecutor *engine.Executor,
) error {
	node, err := newNode(srcDir, src, dstDir, dst, templateExecutor)
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
		// Log
		syncer.log.Debug("sync dir",
			"src", relSrcPath,
			"dst", relDstPath,
		)

		// Destination is a file; remove
		if node.Dst.IsExist && !node.Dst.IsDir {
			if err := os.Remove(node.Dst.Path); err != nil {
				return serror.New("file system error").
					With("file", node.Dst.Path).
					WithErr(std.From(err))
			}

			node.Dst.IsExist = false
		}

		// Destination does not exist; create
		if !node.Dst.IsExist {
			if err := os.MkdirAll(node.Dst.Path, 0o755); err != nil {
				return serror.New("file system error").
					With("dir", node.Dst.Path).
					WithErr(std.From(err))
			}

			// Log
			syncer.log.Info("dir synced",
				"path", relDstPath,
			)
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
				node.TemplateExecutor,
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
			return serror.New("file system error").
				With("dir", node.Dst.Path).
				WithErr(std.From(err))
		}

		for _, file := range files {
			if !dstMap[file.Name()] {
				path := filepath.Join(node.Dst.Path, file.Name())
				if err := os.RemoveAll(path); err != nil {
					return serror.New("file system error").
						With("file", path).
						WithErr(std.From(err))
				}
			}
		}

		return nil
	}

	// Log
	syncer.log.Debug("sync file",
		"src", relSrcPath,
		"dst", relDstPath,
	)

	if node.Dst.IsExist {
		// Destination is a directory; remove
		if node.Dst.IsDir {
			if err := os.RemoveAll(node.Dst.Path); err != nil {
				return serror.New("file system error").
					With("dir", node.Dst.Path).
					WithErr(std.From(err))
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
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return serror.New("file system error").
					With("dir", dir).
					WithErr(std.From(err))
			}
		}
	}

	equal := false

	var srcReader io.Reader

	if node.IsTmpl {
		// Write template
		buffer := &bytes.Buffer{}
		if err := node.TemplateExecutor.ExecuteTemplate(buffer, node.Src.Path); err != nil {
			return err
		}

		srcReader = bytes.NewReader(buffer.Bytes())

		if node.Dst.IsExist {
			// Get template hash
			hash := sha256.New()
			if _, err := io.Copy(hash, buffer); err != nil {
				return err
			}

			equal = bytes.Equal(hash.Sum(nil), node.Dst.Hash)
		}
	} else {
		// Node is not a template, let's go buffering \o/
		srcFile, err := os.Open(node.Src.Path)
		if err != nil {
			return serror.New("file system error").
				With("file", node.Src.Path).
				WithErr(std.From(err))
		}

		defer srcFile.Close()

		if node.Dst.IsExist {
			// Get source hash
			hash := sha256.New()
			if _, err := io.Copy(hash, srcFile); err != nil {
				return err
			}

			equal = bytes.Equal(hash.Sum(nil), node.Dst.Hash)

			if _, err := srcFile.Seek(0, io.SeekStart); err != nil {
				return err
			}
		}

		srcReader = srcFile
	}

	// Files are not equals or destination does not exist
	if !equal {
		// Write to a temporary file in the destination directory, then rename
		// it into place. Renaming is atomic, so an interrupted copy or a flush
		// error leaves any existing destination untouched instead of truncating
		// it in place.
		tmpFile, err := os.CreateTemp(filepath.Dir(node.Dst.Path), ".manala-sync-*")
		if err != nil {
			return serror.New("file system error").
				With("file", node.Dst.Path).
				WithErr(std.From(err))
		}

		tmpPath := tmpFile.Name()

		// Remove the temporary file unless it is successfully renamed into place.
		committed := false
		defer func() {
			if !committed {
				_ = tmpFile.Close()
				_ = os.Remove(tmpPath)
			}
		}()

		// Copy from source to the temporary file
		if _, err := io.Copy(tmpFile, srcReader); err != nil {
			return err
		}

		// Flush and close, surfacing a late write error (e.g. a full disk
		// reported only at close) instead of swallowing it.
		if err := tmpFile.Sync(); err != nil {
			return serror.New("file system error").
				With("file", node.Dst.Path).
				WithErr(std.From(err))
		}
		if err := tmpFile.Close(); err != nil {
			return serror.New("file system error").
				With("file", node.Dst.Path).
				WithErr(std.From(err))
		}

		// Destination file mode: preserve an existing file's permissions
		// (matching the previous in-place truncate), use a standard mode for a
		// new file.
		dstMode := node.Dst.Mode
		if !node.Dst.IsExist {
			dstMode = 0o644
			if node.Src.IsExecutable {
				dstMode = 0o755
			}
		}
		if err := os.Chmod(tmpPath, dstMode); err != nil {
			return serror.New("file system error").
				With("file", node.Dst.Path).
				WithErr(std.From(err))
		}

		// Atomically replace the destination
		if err := os.Rename(tmpPath, node.Dst.Path); err != nil {
			return serror.New("file system error").
				With("file", node.Dst.Path).
				WithErr(std.From(err))
		}
		committed = true

		// Log
		syncer.log.Info("file synced",
			"path", relDstPath,
		)
	} else {
		dstMode := node.Dst.Mode &^ 0o111
		if node.Src.IsExecutable {
			dstMode = node.Dst.Mode | 0o111
		}

		if dstMode != node.Dst.Mode {
			if err := os.Chmod(node.Dst.Path, dstMode); err != nil {
				return serror.New("file system error").
					With("file", node.Dst.Path).
					WithErr(std.From(err))
			}
		}
	}

	return nil
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
	TemplateExecutor *engine.Executor
}

var (
	distRegex = regexp.MustCompile(`(\.dist)(?:$|\.tmpl$)`)
	tmplRegex = regexp.MustCompile(`(\.tmpl)(?:$|\.dist$)`)
)

func newNode(srcDir, src, dstDir, dst string, templateExecutor *engine.Executor) (*node, error) {
	node := &node{}
	node.Src.Dir = srcDir
	node.Dst.Dir = dstDir
	node.TemplateExecutor = templateExecutor

	srcPath := filepath.Join(node.Src.Dir, src)

	// Source stat
	srcStat, err := os.Stat(srcPath)
	if err != nil {
		// Source does not exist
		if errors.Is(err, os.ErrNotExist) {
			return nil, serror.New("no source file or directory").
				With("path", srcPath)
		}

		return nil, serror.New("file system error").
			With("path", srcPath).
			WithErr(std.From(err))
	}

	node.Src.IsDir = srcStat.IsDir()

	if node.Src.IsDir {
		files, err := os.ReadDir(srcPath)
		if err != nil {
			return nil, serror.New("file system error").
				With("dir", srcPath).
				WithErr(std.From(err))
		}

		for _, file := range files {
			node.Src.Files = append(node.Src.Files, file.Name())
		}
	} else {
		node.Src.IsExecutable = (srcStat.Mode() & 0o100) != 0

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
			return nil, serror.New("file system error").
				With("path", dstPath).
				WithErr(std.From(err))
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
				return nil, serror.New("file system error").
					With("dir", dstPath).
					WithErr(std.From(err))
			}

			for _, file := range files {
				node.Dst.Files = append(node.Dst.Files, file.Name())
			}
		} else {
			// Get destination hash
			file, err := os.Open(dstPath)
			if err != nil {
				return nil, serror.New("file system error").
					With("file", dstPath).
					WithErr(std.From(err))
			}

			defer file.Close()

			hash := sha256.New()
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
