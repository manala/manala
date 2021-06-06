package fs

import (
	ioFs "io/fs"
	"os"
	"path"
)

/***********/
/* Manager */
/***********/

// NewManager returns a file system manager
func NewManager() *manager {
	return &manager{}
}

type ManagerInterface interface {
	NewDirFs(dir string) *Fs
}

type manager struct {
}

// NewDirFs returns a dir file system
func (manager *manager) NewDirFs(dir string) *Fs {
	return &Fs{
		dir: dir,
	}
}

/***************/
/* File System */
/***************/

type ReadInterface interface {
	ioFs.FS
	ioFs.StatFS
	ioFs.ReadFileFS
	ioFs.ReadDirFS
}

type WriteInterface interface {
	OpenFile(name string, flag int, perm ioFs.FileMode) (*os.File, error)
	Chmod(name string, mode ioFs.FileMode) error
	Remove(name string) error
	MkdirAll(path string, perm ioFs.FileMode) error
	RemoveAll(path string) error
}

type ReadWriteInterface interface {
	ReadInterface
	WriteInterface
}

type Fs struct {
	dir string
}

func (fs *Fs) path(name string) string {
	return path.Join(
		fs.dir,
		name,
	)
}

// Open the named file for reading
func (fs *Fs) Open(name string) (ioFs.File, error) {
	return os.Open(fs.path(name))
}

// Stat returns a FileInfo describing the named file
func (fs *Fs) Stat(name string) (ioFs.FileInfo, error) {
	return os.Stat(fs.path(name))
}

// ReadFile reads the named file and returns the contents
func (fs *Fs) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(fs.path(name))
}

// ReadDir reads the named directory
func (fs *Fs) ReadDir(name string) ([]ioFs.DirEntry, error) {
	return os.ReadDir(fs.path(name))
}

// OpenFile is the generalized open call
func (fs *Fs) OpenFile(name string, flag int, perm ioFs.FileMode) (*os.File, error) {
	return os.OpenFile(fs.path(name), flag, perm)
}

// Chmod changes the mode of the named file to mode
func (fs *Fs) Chmod(name string, mode ioFs.FileMode) error {
	return os.Chmod(fs.path(name), mode)
}

// Remove removes the named file or (empty) directory
func (fs *Fs) Remove(name string) error {
	return os.Remove(fs.path(name))
}

// MkdirAll creates a directory named path
func (fs *Fs) MkdirAll(path string, perm ioFs.FileMode) error {
	return os.MkdirAll(fs.path(path), perm)
}

// RemoveAll removes path and any children it contains
func (fs *Fs) RemoveAll(path string) error {
	return os.RemoveAll(fs.path(path))
}
