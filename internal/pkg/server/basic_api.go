package server

import (
	"errors"
	"osdsvr/internal/pkg/core"
	"osdsvr/internal/pkg/metadata"
	"osdsvr/pkg/proto/osdpb"
)

var (
	ErrInvalidMetaData = errors.New("request metadata is invalid, please check if parameter is right")
)

func (s *OsdServer) UploadFile(stream osdpb.OsdService_UploadFileServer) error {
	// handle metadata
	req, _ := stream.Recv()
	fileName := ""
	bucketID := ""
	if req.MetaData == nil || req.MetaData.BucketId == "" || req.MetaData.Name == "" {
		return ErrInvalidMetaData
	}
	fileName = req.MetaData.Name
	bucketID = req.MetaData.BucketId
	tmpPath := s.oService.GetObjectTmpPath(fileName)
	done := make(chan bool)
	task := &core.Task{
		DoTask: func(v ...any) error {
			path := v[0].(string)
			bucketId := v[1].(string)
			st := v[2].(osdpb.OsdService_UploadFileServer)
			ch := v[3].(chan bool)
			name := v[4].(string)
			if path == "" || bucketId == "" {
				ch <- true
				return ErrInvalidMetaData
			}
			// write tmp object
			hash, err := s.oService.WriteObject(path, st)
			if err != nil {
				ch <- true
				return err
			}
			// object check
			if err := s.oService.CheckObject(); err != nil {
				ch <- true
				return err
			}
			// object rename
			newp := s.oService.GetObjectPath(bucketId, hash)
			err = s.oService.RenameObject(path, newp)
			if err != nil {
				ch <- true
				return err
			}
			// add version record
			record := &metadata.VersionRecord{
				Hash:       hash,
				Size:       0,
				UploadTime: 0,
			}
			err = s.mService.AddVerisonRecord(bucketId, name, record)
			if err != nil {
				ch <- true
				return err
			}
			ch <- true
			return nil
		},
		Param: []any{tmpPath, bucketID, stream, done, fileName},
	}
	s.scheduler.AddTask(task)
	<-done
	return nil
}

func (s *OsdServer) DownloadFile(req *osdpb.FileDownloadRequest, stream osdpb.OsdService_DownloadFileServer) error {
	if req.MetaData == nil {
		return errors.New("metadata is nil")
	}
	objectName := ""
	bucketID := ""
	var version int64
	objectName = req.MetaData.Name
	bucketID = req.MetaData.BucketId
	version = req.MetaData.Version
	if objectName == "" || bucketID == "" {
		return ErrInvalidMetaData
	}
	// find metadata
	record, err := s.mService.GetVerisonRecord(bucketID, objectName, version)
	if err != nil {
		return err
	}
	hash := record.Hash
	// read object
	s.oService.GetObject(hash, bucketID, stream)
	return nil
}
