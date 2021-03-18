package models

import (
	"manala/fs"
)

/***********/
/* Manager */
/***********/

// Create a model file system manager
func NewFsManager(manager fs.ManagerInterface) *fsManager {
	return &fsManager{
		fsManager: manager,
	}
}

type FsManagerInterface interface {
	NewModelFs(model model) fs.ReadWriteInterface
}

type fsManager struct {
	fsManager fs.ManagerInterface
}

/***********************/
/* File System - Model */
/***********************/

// Create a model file system
func (manager *fsManager) NewModelFs(model model) fs.ReadWriteInterface {
	return manager.fsManager.NewDirFs(model.getDir())
}
