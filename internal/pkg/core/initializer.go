package core

import (
	"osdsvr/internal/pkg/fs"
	"strings"
)

type Initializer interface {
	LoadInitFile(path string, initFileName string) error
	SyncInitFile() error
	GetInitFilePath() string
}

type OSDInitializer struct {
	initFile     *fs.FileIO
	path         string
	initFileName string
}

func (t *OSDInitializer) LoadInitFile(path string, initFileName string) error {
	// check init dir exist
	if _, err := fs.PathExists(path); err != nil {
		return err
	}
	filePath := strings.Join([]string{path, initFileName}, "")
	file, err := fs.NewFileIO(filePath)
	if err != nil {
		return err
	}
	t.path = path
	t.initFileName = initFileName
	t.initFile = file
	return nil
}

func (t *OSDInitializer) GetInitFilePath() string {
	return strings.Join([]string{t.path, t.initFileName}, "")
}

func (t *OSDInitializer) SyncInitFile() error {
	// abstract method
	return nil
}
