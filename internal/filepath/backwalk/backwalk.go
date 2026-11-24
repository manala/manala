package backwalk

import (
	"io/fs"
	"os"
	"path/filepath"
)

// walkDir ascends path, calling walkDirFunc.
func walkDir(path string, dir os.DirEntry, walkDirFunc fs.WalkDirFunc) error {
	if err := walkDirFunc(path, dir, nil); err != nil || !dir.IsDir() {
		if err == filepath.SkipAll && dir.IsDir() {
			// Successfully skipped directory
			err = nil
		}

		return err
	}

	pathAbs, err := filepath.Abs(path)
	if err != nil {
		// Report the error; if the callback does not abort, stop ascending
		// rather than continuing with an invalid path.
		if err := walkDirFunc(path, dir, err); err != nil {
			return err
		}

		return nil
	}

	// Get parent dir
	parent := filepath.Join(path, "..")
	parentAbs, err := filepath.Abs(parent)
	if err != nil {
		// Report the error; if the callback does not abort, stop ascending
		// rather than continuing with an invalid parent path.
		if err := walkDirFunc(parent, dir, err); err != nil {
			return err
		}

		return nil
	}

	// If absolute parent path equals to absolute path,
	// we have reached the filesystem root
	if parentAbs == pathAbs {
		return nil
	}

	info, err := os.Lstat(parent)
	if err != nil {
		// Report the error; if the callback does not abort, stop here rather
		// than recursing with a nil DirEntry (which would panic on IsDir()).
		if err := walkDirFunc(parent, dir, err); err != nil {
			return err
		}

		return nil
	}

	if err := walkDir(parent, fs.FileInfoToDirEntry(info), walkDirFunc); err != nil {
		if err == filepath.SkipAll {
			// Successfully skipped directory
			err = nil
		}

		return err
	}

	return nil
}

// WalkDir back walks the file tree at path, calling fn for each
// directory in the tree, including root.
// Greatly inspired by filepath.WalkDir.
func WalkDir(dir string, fn fs.WalkDirFunc) error {
	info, err := os.Lstat(dir)
	if err != nil {
		err = fn(dir, nil, err)
	} else {
		err = walkDir(dir, fs.FileInfoToDirEntry(info), fn)
	}

	if err == filepath.SkipAll {
		return nil
	}

	return err
}
