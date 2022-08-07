package server

import (
	"context"
	"encoding/json"
	"errors"
	"osdsvr/internal/pkg/config"
	"osdsvr/internal/pkg/core"
	"osdsvr/internal/pkg/metadata"
	"osdsvr/pkg/proto/osdpb"
	"osdsvr/pkg/zlog"

	"go.uber.org/zap"
)

var (
	ErrInvalidMetaData       = errors.New("request metadata is invalid, please check if parameter is right")
	ErrVerisonRecordNotFound = errors.New("version record is not find")
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
		DoTask: func(v ...interface{}) error {
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
			// create redundant object data
			if *config.RS_CODE_MODE {
				err = s.rsService.CreateRedundantObject(bucketId, hash, s.scheduler)
				if err != nil {
					ch <- true
					return err
				}
			}
			ch <- true
			return nil
		},
		Param: []interface{}{tmpPath, bucketID, stream, done, fileName},
	}
	s.scheduler.AddTask(task)
	// if task != nil {
	// 	done <- true
	// }
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
	if record == nil {
		return ErrVerisonRecordNotFound
	}
	hash := record.Hash
	// read object
	if err := s.oService.GetObject(hash, bucketID, stream); err != nil {
		if *config.RS_CODE_MODE {
			zlog.Info("object data loss detected! try to reconstruct.", zap.Any("bucket id", bucketID), zap.Any("object name", hash))
			if err := s.rsService.ReconstructObject(bucketID, hash); err != nil {
				zlog.Error("reconstruct failed", zap.Error(err))
				return err
			}
			if err := s.oService.GetObject(hash, bucketID, stream); err != nil {
				zlog.Error("reconstruct failed", zap.Error(err))
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (s *OsdServer) QueryObjectVersionRecords(ctx context.Context, req *osdpb.QueryObjectVersionRecordsRequest) (*osdpb.QueryObjectVersionRecordsResponse, error) {
	resp := &osdpb.QueryObjectVersionRecordsResponse{}
	bucketID := req.BucketId
	objectName := req.ObjectName
	if bucketID == "" || objectName == "" {
		return resp, nil
	}
	versionMap, err := s.mService.GetAllVersionRecords(bucketID, objectName)
	if err != nil {
		return resp, err
	}
	versionInfoMap := make(map[int64]*metadata.VersionInfo)
	for k, v := range versionMap {
		info := metadata.CreateVersionInfo(v)
		versionInfoMap[k] = info
	}
	rlt, err := json.Marshal(versionInfoMap)
	if err != nil {
		return resp, err
	}
	resp.ResultCode = osdpb.Result_SUCCESS
	resp.Desc = string(rlt)
	return resp, nil
}
