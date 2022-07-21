package server

import (
	"context"
	"io"
	"osdsvr/pkg/proto/osdpb"
	"osdsvr/pkg/zlog"

	"go.uber.org/zap"
)

// SayHello implements helloworld.GreeterServer
func (s *OsdServer) SayHello(ctx context.Context, in *osdpb.HelloRequest) (*osdpb.HelloReply, error) {
	zlog.Info("received message", zap.String("name", in.GetName()))
	return &osdpb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *OsdServer) UploadFile(stream osdpb.OsdService_UploadFileServer) error {
	filedata := make([]byte, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// upload completed
			// todo
			if len(filedata) != 0 {
				content := string(filedata)
				zlog.Info("accept file content", zap.String("content", content))
			}
			return stream.SendAndClose(&osdpb.FileUploadResponse{
				ResultCode: osdpb.UploadResult_SUCCESS,
				Desc:       "upload file sucessfully",
			})
		}
		if err != nil {
			zlog.Error("faileb to receive chunk", zap.Error(err))
			return stream.SendAndClose(&osdpb.FileUploadResponse{
				ResultCode: osdpb.UploadResult_FAILED,
				Desc:       "receive chunk failed",
			})
		}
		if req.MetaData != nil {
			// todo
			zlog.Info("file metadata detected.")
		}
		filedata = append(filedata, req.GetChunk()...)
	}
}
