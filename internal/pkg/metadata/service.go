package metadata

import (
	"osdsvr/internal/pkg/config"
	"osdsvr/internal/pkg/fs"
	"strings"

	"google.golang.org/protobuf/proto"
)

var (
	metadataPostfix = ".meta"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) IfMetadataExist(bucketID string, objectName string) bool {
	path := s.GetMetadataPath(bucketID, objectName)
	return fs.Exists(path)
}

func (s *Service) GetMetadata(bucketID string, objectName string) (*MetaData, error) {
	dir := strings.Join([]string{*config.STORAGE_PATH, bucketID, "/metadata/"}, "")
	fullpath := strings.Join([]string{dir, objectName, metadataPostfix}, "")
	if !s.IfMetadataExist(bucketID, objectName) {
		return nil, nil
	}
	file, err := fs.NewFileIO(fullpath)
	if err != nil {
		return nil, err
	}
	data, err := file.ReadAllBytes()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	meta := &MetaData{}
	if err := proto.Unmarshal(data, meta); err != nil {
		return nil, err
	}
	return meta, nil
}

func (s *Service) AddVerisonRecord(bucketID string, objectName string, record *VersionRecord) error {
	meta, err := s.GetMetadata(bucketID, objectName)
	if err != nil {
		return err
	}
	if meta != nil {
		latestVersion := meta.LatestVersion
		if meta.Versions == nil {
			meta.Versions = make(map[int64]*VersionRecord)
		}
		record.Version = latestVersion + 1
		meta.LatestVersion = latestVersion + 1
		meta.Versions[record.Version] = record
	} else {
		meta = &MetaData{
			Name:          objectName,
			BucketId:      bucketID,
			Versions:      make(map[int64]*VersionRecord),
			LatestVersion: 1,
		}
		meta.Versions[1] = record
	}
	if err := s.Sync(meta); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetVerisonRecord(bucketID string, objectName string, version int64) (*VersionRecord, error) {
	var v int64
	meta, err := s.GetMetadata(bucketID, objectName)
	if err != nil {
		return nil, err
	}
	if meta != nil {
		if version == 0 && meta.LatestVersion != 0 {
			v = meta.LatestVersion
		} else {
			v = version
		}
		if meta.Versions != nil {
			return meta.Versions[v], nil
		}
	}
	return nil, nil
}

func (s *Service) GetLatestVerisonRecord(bucketID string, objectName string) (*VersionRecord, error) {
	meta, err := s.GetMetadata(bucketID, objectName)
	if err != nil {
		return nil, err
	}
	if meta != nil {
		if meta.Versions != nil && meta.LatestVersion != 0 {
			return meta.Versions[meta.LatestVersion], nil
		}
	}
	return nil, nil
}

func (s *Service) Sync(meta *MetaData) error {
	if meta == nil {
		return nil
	}
	dir := strings.Join([]string{*config.STORAGE_PATH, meta.BucketId, "/metadata/"}, "")
	fs.CreatePathIfNotExists(dir)
	fullpath := strings.Join([]string{dir, meta.Name, metadataPostfix}, "")
	file, err := fs.NewFileIO(fullpath)
	if err != nil {
		return err
	}
	if err := file.DeleteAllBytes(); err != nil {
		return err
	}
	data, err := proto.Marshal(meta)
	if err != nil {
		return err
	}
	_, err = file.WriteAt(data, 0)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetMetadataPath(bucketID string, objectName string) string {
	dir := strings.Join([]string{*config.STORAGE_PATH, bucketID, "/metadata/"}, "")
	path := strings.Join([]string{dir, objectName, metadataPostfix}, "")
	return path
}
