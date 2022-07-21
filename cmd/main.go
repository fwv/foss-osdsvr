package main

import (
	"flag"
	"net"
	"osdsvr/internal/pkg/server"
	"osdsvr/pkg/proto/osdpb"
	"osdsvr/pkg/zlog"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", ":50051", "grpc server listen address")
)

func main() {
	flag.Parse()

	// 监听grpc端口
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		zlog.Fatal("failed to listen.", zap.Error(err))
	}
	s := grpc.NewServer()

	osdServer := server.NewOsdServer()

	// 注册osd服务
	osdpb.RegisterOsdServiceServer(s, osdServer)

	zlog.Info("server start listening", zap.String("addr", *addr))
	if err := s.Serve(lis); err != nil {
		zlog.Fatal("failed to serve", zap.Error(err))
	}
}
