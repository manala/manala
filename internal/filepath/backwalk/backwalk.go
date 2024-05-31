package backwalk

import (
	"io/fs"
	"os"
	"path/filepath"
)

// backwalkDir ascends path, calling backwalkDirFunc.
func backwalkDir(path string, dir os.DirEntry, backwalkDirFunc fs.WalkDirFunc) error {
	if err := backwalkDirFunc(path, dir, nil); err != nil || !dir.IsDir() {
		if err == filepath.SkipAll && dir.IsDir() {
			// Successfully skipped directory
			err = nil
		}

		return err
	}

	pathAbs, err := filepath.Abs(path)
	if err != nil {
		// Second call, to report Abs error
		err = backwalkDirFunc(path, dir, err)
		if err != nil {
			return err
		}
	}

	// Get parent dir
	parent := filepath.Join(path, "..")
	parentAbs, err := filepath.Abs(parent)

	if err != nil {
		// Second call, to report parent Abs error
		err = backwalkDirFunc(parent, dir, err)
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
		err = backwalkDirFunc(parent, dir, err)
		if err != nil {
			return err
		}
	}

	if err := backwalkDir(parent, fs.FileInfoToDirEntry(info), backwalkDirFunc); err != nil {
		if err == filepath.SkipAll {
			// Successfully skipped directory
			err = nil
		}

		return err
	}

	return nil
}

// BackwalkDir backwalks the file tree at path, calling fn for each
// directory in the tree, including root.
// Greatly inspired by filepath.WalkDir.
func BackwalkDir(dir string, fn fs.WalkDirFunc) error {
	info, err := os.Lstat(dir)
	if err != nil {
		err = fn(dir, nil, err)
	} else {
		err = backwalkDir(dir, fs.FileInfoToDirEntry(info), fn)
	}

	if err == filepath.SkipAll {
		return nil
	}

	return err
}
