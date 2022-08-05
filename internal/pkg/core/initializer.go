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
	InitFile     *fs.FileIO
	Path         string
	InitFileName string
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
	t.Path = path
	t.InitFileName = initFileName
	t.InitFile = file
	return nil
}

func (t *OSDInitializer) GetInitFilePath() string {
	return strings.Join([]string{t.Path, t.InitFileName}, "")
}

func (t *OSDInitializer) SyncInitFile() error {
	// abstract method
	return nil
}
