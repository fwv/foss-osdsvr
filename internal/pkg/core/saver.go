package core

import (
	"errors"
	"os"
	"osdsvr/internal/pkg/fs"
	"osdsvr/pkg/zlog"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Entry struct {
	fid int64
	loc int64
}

type Saver interface {
	Store(data []byte) (*Entry, error)
	AsyncStore(data []byte, buildIdx func(oid int64, entry *Entry) error) error
	Load(*Entry) ([]byte, error)
}

type OSDSaver struct {
	mu sync.Mutex
	OSDInitializer
	activeFile       fs.File
	maxFileSizeLimit int64
	taskCh           chan (int64)
	maxFileID        int64
	path             string
}

func NewOSDSaver(path string, maxFileSizeLimit int64) *OSDSaver {
	s := &OSDSaver{
		maxFileSizeLimit: maxFileSizeLimit,
		taskCh:           make(chan int64),
		maxFileID:        0,
		path:             path,
	}
	err := s.LoadInitFile(path, "osdsaver.ini")
	if err != nil {
		zlog.Error("OSDSaver load init file failed", zap.Error(err))
	}
	return s
}

func (s *OSDSaver) LoadInitFile(path string, initFileName string) error {
	s.OSDInitializer.LoadInitFile(path, initFileName)
	// read init from init file
	if s.initFile == nil {
		zlog.Error("init file is nil", zap.String("path", path), zap.String("file name", initFileName))
		return errors.New("init file is nill")
	}
	// skip read for the first time
	if s.initFile.Size() == 0 {
		return nil
	}
	r, err := s.initFile.ReadAt(0)
	if err != nil {
		return err
	}
	init := &OSDSaverInit{}
	proto.Unmarshal(r.GetValue(), init)
	s.maxFileID = init.MaxFileID
	zlog.Info("load init file for OSDSaver successfully", zap.String("path", path), zap.Int64("maxFileID", init.MaxFileID))
	return nil
}

func (s *OSDSaver) SyncInitFile() error {
	if s.initFile != nil {
		init := &OSDSaverInit{
			MaxFileID: s.maxFileID,
		}
		initData, err := proto.Marshal(init)
		if err != nil {
			return err
		}
		os.Truncate(s.GetInitFilePath(), 0)
		r := fs.NewRecord(-1, initData, 1, time.Now().Unix())
		_, _, err = s.initFile.WriteAt(r, 0)
		if err != nil {
			return err
		}
		zlog.Info("sync init file sucessfully", zap.Any("init file", s.GetInitFilePath()), zap.Any("maxFileId", s.maxFileID))
	}
	return nil
}

func (s *OSDSaver) Store(data []byte) (*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// check if file is too big
	dataSize := len(data)
	if dataSize > int(s.maxFileSizeLimit) {
		// todo: big file
		zlog.Info("processing big file", zap.Int("data size", dataSize))
		return nil, nil
	}
	if s.activeFile == nil || s.activeFile.Size()+int64(dataSize) > s.maxFileSizeLimit {
		if err := s.switchActiveFile(); err != nil {
			return nil, err
		}
	}
	// write data into file
	start, err := s.activeFile.Append(fs.NewRecord(-1, data, 1, time.Now().Unix()))
	if err != nil {
		return nil, nil
	}
	entry := &Entry{
		fid: s.maxFileID,
		loc: int64(start),
	}
	return entry, nil
}

func (s *OSDSaver) AsyncStore(data []byte, buildIdx func(oid int64, entry *Entry) error) error {
	return nil
}

func (s *OSDSaver) Load(entry *Entry) ([]byte, error) {
	filePath, err := s.fetchFilePath(entry)
	if err != nil {
		return nil, err
	}
	f, err := fs.NewOSDFile(filePath)
	if err != nil {
		return nil, err
	}
	r, err := f.ReadAt(int(entry.loc))
	if err != nil {
		zlog.Error("failed to read file", zap.Int("off", int(entry.loc)))
		return nil, err
	}
	return r.GetValue(), nil
}

func (s *OSDSaver) switchActiveFile() error {
	// handle init data
	s.maxFileID++
	if err := s.SyncInitFile(); err != nil {
		zlog.Error("OSDSaver sync init file failed", zap.Error(err))
		return err
	}
	// hander active file
	newFilePath, err := s.fetchFilePath(&Entry{fid: s.maxFileID})
	if err != nil {
		return err
	}
	newfile, err := s.createFile(newFilePath)
	if err != nil {
		return err
	}
	// todo: close old file?
	if s.activeFile != nil {
		s.activeFile.Sync()
	}
	s.activeFile = newfile
	return nil
}

func (s *OSDSaver) createFile(filePath string) (fs.File, error) {
	if _, err := fs.PathExists(s.path); err != nil {
		return nil, err
	}
	return fs.NewOSDFile(filePath)
}

func (s *OSDSaver) fetchFilePath(entry *Entry) (string, error) {
	if entry.fid == 0 {
		return "", errors.New("entry fid can't be 0")
	}
	str := strings.Join([]string{s.path, strconv.FormatInt(entry.fid, 10), fs.FILE_SUFFIX_STR}, "")
	return str, nil
}
