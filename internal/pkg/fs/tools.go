package fs

import (
	"os"
	"osdsvr/pkg/zlog"

	"go.uber.org/zap"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// PathExists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			zlog.Error("failed to make dir", zap.Error(err))
		} else {
			return true, nil
		}
	}
	return false, err
}

func CreatePathIfNotExists(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			zlog.Error("failed to make dir", zap.Error(err))
			return err
		} else {
			return nil
		}
	}
	return nil
}
