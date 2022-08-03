package server

import (
	"context"
	"errors"
	"io"
	"os"
	"osdsvr/internal/pkg/config"
	"osdsvr/internal/pkg/core"
	"osdsvr/pkg/proto/osdpb"
	"osdsvr/pkg/zlog"
	"strings"

	"go.uber.org/zap"
)

var (
	ErrInvalidMetaData = errors.New("request metadata is invalid, please check if parameter is right")
)

// SayHello implements helloworld.GreeterServer
func (s *OsdServer) SayHello(ctx context.Context, in *osdpb.HelloRequest) (*osdpb.HelloReply, error) {
	zlog.Info("received message", zap.String("name", in.GetName()))
	return &osdpb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *OsdServer) UploadFile(stream osdpb.OsdService_UploadFileServer) error {
	// handle metadata
	req, _ := stream.Recv()
	fileName := ""
	bucketID := ""
	if req.MetaData != nil {
		fileName = req.MetaData.Name
		bucketID = req.MetaData.BucketId
	}
	tmpPath := s.oService.GetObjectTmpPath(fileName)
	done := make(chan bool)
	task := &core.Task{
		DoTask: func(v ...any) error {
			path := v[0].(string)
			bucketId := v[1].(string)
			st := v[2].(osdpb.OsdService_UploadFileServer)
			ch := v[3].(chan bool)
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
			// modify metadata
			ch <- true
			return nil
		},
		Param: []any{tmpPath, bucketID, stream, done},
	}
	s.scheduler.AddTask(task)
	<-done
	return nil
}

func (s *OsdServer) DownloadFile(in *osdpb.FileDownloadRequest, stream osdpb.OsdService_DownloadFileServer) error {
	if in.MetaData == nil {
		return errors.New("metadata is nil")
	}
	objectName := in.MetaData.Name
	fullObjectPath := strings.Join([]string{*config.STORAGE_PATH, objectName}, "")
	file, err := os.OpenFile(fullObjectPath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	readoff := 0
	data := make([]byte, *config.DOWNLOAD_CHUNK_SIZE)
	for {
		n, err := file.ReadAt(data, int64(readoff))
		readoff += n
		if n != 0 {
			zlog.Debug("send chunk data to stream", zap.Int("chunk size", n))
			if err := stream.Send(&osdpb.FileDownloadResponse{
				Chunk: data[:n],
			}); err != nil {
				zlog.Error("failed to send chunk data to stream", zap.Error(err))
				return err
			}
		}
		if err == io.EOF {
			zlog.Info("send file data to stream completed", zap.String("object name", objectName), zap.Int("total size", readoff))
			break
		} else if err != nil {
			zlog.Error("failed to read file", zap.Error(err))
			return err
		}
	}
	return nil
}
