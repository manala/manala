package backwalk

import (
	"io/fs"
	"os"
	"path/filepath"
)

// backwalk ascends path, calling backwalkFn.
func backwalk(path string, file os.DirEntry, walkDirFunc fs.WalkDirFunc) error {
	if err := walkDirFunc(path, file, nil); err != nil || !file.IsDir() {
		if err == filepath.SkipDir && file.IsDir() {
			// Successfully skipped directory
			err = nil
		}
		return err
	}

	pathAbs, err := filepath.Abs(path)
	if err != nil {
		// Second call, to report Abs error
		err = walkDirFunc(path, file, err)
		if err != nil {
			return err
		}
	}

	// Get parent dir
	parentPath := filepath.Join(path, "..")
	parentPathAbs, err := filepath.Abs(parentPath)
	if err != nil {
		// Second call, to report parent Abs error
		err = walkDirFunc(parentPath, file, err)
		if err != nil {
			return err
		}
	}

	// If absolute parent path equals to absolute path,
	// we have reached the filesystem root
	if parentPathAbs == pathAbs {
		return nil
	}

	info, err := os.Lstat(parentPath)
	if err != nil {
		// Second call, to report Stat error.
		err = walkDirFunc(parentPath, file, err)
		if err != nil {
			return err
		}
	}

	if err := backwalk(parentPath, fs.FileInfoToDirEntry(info), walkDirFunc); err != nil {
		if err == filepath.SkipDir {
			// Successfully skipped directory
			err = nil
		}
		return err
	}

	return nil
}

// Backwalk backwalks the file tree at path, calling fn for each
// directory in the tree, including root.
// Greatly inspired by filepath.WalkDir
func Backwalk(path string, walkDirFunc fs.WalkDirFunc) error {
	info, err := os.Lstat(path)
	if err != nil {
		err = walkDirFunc(path, nil, err)
	} else {
		err = backwalk(path, fs.FileInfoToDirEntry(info), walkDirFunc)
	}
	if err == filepath.SkipDir {
		return nil
	}
	return err
}
