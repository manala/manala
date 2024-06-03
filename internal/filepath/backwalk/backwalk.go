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
		// Second call, to report Abs error
		err = walkDirFunc(path, dir, err)
		if err != nil {
			return err
		}
	}

	// Get parent dir
	parent := filepath.Join(path, "..")
	parentAbs, err := filepath.Abs(parent)

	if err != nil {
		// Second call, to report parent Abs error
		err = walkDirFunc(parent, dir, err)
		if err != nil {
			return err
		}
	}

	// If absolute parent path equals to absolute path,
	// we have reached the filesystem root
	if parentAbs == pathAbs {
		return nil
	}

	info, err := os.Lstat(parent)
	if err != nil {
		// Second call, to report Stat error.
		err = walkDirFunc(parent, dir, err)
		if err != nil {
			return err
		}
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
