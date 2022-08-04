package server

import (
	"net"
	"osdsvr/internal/pkg/config"
	"osdsvr/internal/pkg/core"
	"osdsvr/internal/pkg/metadata"
	"osdsvr/internal/pkg/object"
	"osdsvr/pkg/proto/osdpb"
	"osdsvr/pkg/zlog"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type OsdServer struct {
	osdpb.UnimplementedOsdServiceServer
	scheduler *core.Scheduler
	oService  *object.Service
	mService  *metadata.Service
}

func NewOsdServer() *OsdServer {
	return &OsdServer{
		scheduler: core.NewScheduler(*config.SCHEDULE_CAPACITY),
		oService:  object.NewService(),
	}
}

func (s *OsdServer) Serve() error {
	// start scheduler
	if s.scheduler != nil {
		if err := s.scheduler.Start(*config.SCHEDULE_CONCURRENCY); err != nil {
			zlog.Fatal("failed to start scheduler", zap.Error(err))
			return err
		}
	}
	// listen grcp
	lis, err := net.Listen("tcp", *config.GRPC_ADDR)
	if err != nil {
		zlog.Fatal("failed to listen.", zap.Error(err))
	}
	server := grpc.NewServer()
	osdpb.RegisterOsdServiceServer(server, s)
	zlog.Info("server start listening", zap.String("addr", *config.GRPC_ADDR))
	if err := server.Serve(lis); err != nil {
		zlog.Fatal("failed to serve", zap.Error(err))
		return err
	}
	return nil
}
func (s *OsdServer) Shutdown() {
	if s.scheduler != nil {
		s.scheduler.Shutdown()
	}
}
