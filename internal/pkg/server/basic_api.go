package server

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"
	"osdsvr/pkg/proto/osdpb"
	"osdsvr/pkg/zlog"
	"strings"

	"go.uber.org/zap"
)

var (
	objectDir         = flag.String("objectDir", "/home/fwv/oss/", "object directory")
	downloadChunkSize = flag.Int("downloadChunkSize", 1024*64, "download chunk size")
)

// SayHello implements helloworld.GreeterServer
func (s *OsdServer) SayHello(ctx context.Context, in *osdpb.HelloRequest) (*osdpb.HelloReply, error) {
	zlog.Info("received message", zap.String("name", in.GetName()))
	return &osdpb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *OsdServer) UploadFile(stream osdpb.OsdService_UploadFileServer) error {
	writeoff := 0
	fullFileName := ""
	for {
		req, err := stream.Recv()
		// handle EOF
		if err == io.EOF {
			// todo
			zlog.Info("accept file completed", zap.String("object save path", fullFileName), zap.Int("content size", writeoff))
			return stream.SendAndClose(&osdpb.FileUploadResponse{
				ResultCode: osdpb.Result_SUCCESS,
				Desc:       "upload file sucessfully",
			})
		}

		// handle error
		if err != nil {
			zlog.Error("failed to receive chunk", zap.Error(err))
			return stream.SendAndClose(&osdpb.FileUploadResponse{
				ResultCode: osdpb.Result_FAILED,
				Desc:       "receive chunk failed",
			})
		}
		fileName := ""
		// handle metadata
		if req.MetaData != nil {
			fileName = req.MetaData.Name
		}
		var file *os.File
		if file == nil {
			fullFileName = strings.Join([]string{*objectDir, fileName}, "")
			file, err = os.OpenFile(fullFileName, os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				return err
			}
			defer file.Close()
		}
		n, err := file.WriteAt(req.Chunk, int64(writeoff))
		if err != nil {
			return err
		}
		zlog.Debug("write chunk data to file completed", zap.Int("chunk size", n))
		writeoff += n
		if err := file.Sync(); err != nil {
			return err
		}
	}
}

func (s *OsdServer) DownloadFile(in *osdpb.FileDownloadRequest, stream osdpb.OsdService_DownloadFileServer) error {
	if in.MetaData == nil {
		return errors.New("metadata is nil")
	}
	objectName := in.MetaData.Name
	fullObjectPath := strings.Join([]string{*objectDir, objectName}, "")
	file, err := os.OpenFile(fullObjectPath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	readoff := 0
	data := make([]byte, *downloadChunkSize)
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
