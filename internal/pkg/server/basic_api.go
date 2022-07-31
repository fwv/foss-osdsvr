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
	objectDir = flag.String("objectDir", "/home/fwv/oss/", "object directory")
)

// SayHello implements helloworld.GreeterServer
func (s *OsdServer) SayHello(ctx context.Context, in *osdpb.HelloRequest) (*osdpb.HelloReply, error) {
	zlog.Info("received message", zap.String("name", in.GetName()))
	return &osdpb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *OsdServer) UploadFile(stream osdpb.OsdService_UploadFileServer) error {
	filedata := make([]byte, 0)
	// receivedSize := 0
	writeoff := 0
	for {
		req, err := stream.Recv()
		// handle EOF
		if err == io.EOF {
			// todo
			if len(filedata) != 0 {
				content := string(filedata)
				zlog.Info("accept file content", zap.String("content", content))
			}
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
			// todo
			zlog.Info("file metadata detected.")
			fileName = req.MetaData.Name
		}
		var file *os.File
		if file == nil {
			fullFileName := strings.Join([]string{*objectDir, fileName}, "")
			file, err = os.OpenFile(fullFileName, os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				return err
			}
		}
		n, err := file.WriteAt(req.Chunk, int64(writeoff))
		if err != nil {
			return err
		}
		writeoff += n
		if err := file.Sync(); err != nil {
			return err
		}
		// handle chunk data
		filedata = append(filedata, req.GetChunk()...)
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
	readoff := 0
	data := make([]byte, 1024)
	for {
		n, err := file.ReadAt(data, int64(readoff))
		readoff += n
		if err != nil {
			if err == io.EOF {
				stream.Send(&osdpb.FileDownloadResponse{
					ResultCode: osdpb.Result_SUCCESS,
					Desc:       "",
				})
				// todo handle eof
				zlog.Info("read file finish")
			}
			return err
		}
		stream.Send(&osdpb.FileDownloadResponse{
			Chunk: data[:n],
		})
	}
}
