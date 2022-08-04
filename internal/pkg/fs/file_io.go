package fs

import (
	"errors"
	"os"
)

type FileIO struct {
	fd *os.File
}

func NewFileIO(path string) (*FileIO, error) {
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &FileIO{fd: fd}, nil
}

func (w *FileIO) WriteAt(b []byte, off int64) (n int, err error) {
	return w.fd.WriteAt(b, off)
}

func (w *FileIO) ReadAt(b []byte, off int64) (n int, err error) {
	return w.fd.ReadAt(b, off)
}

func (w *FileIO) Sync() (err error) {
	return w.fd.Sync()
}

func (w *FileIO) Close() (err error) {
	return w.fd.Close()
}

func (w *FileIO) ReadAllBytes() (data []byte, err error) {
	filesize, err := w.GetFileSize()
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, filesize)

	_, err = w.fd.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (w *FileIO) GetFileSize() (n int64, err error) {
	if w.fd == nil {
		return -1, errors.New("fd is nil")
	}
	fileinfo, err := w.fd.Stat()
	if err != nil {
		return -1, err
	}
	filesize := fileinfo.Size()
	return filesize, nil
}

func (w *FileIO) DeleteAllBytes() (err error) {
	return w.fd.Truncate(0)
}
