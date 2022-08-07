package main

import (
	"context"
	"flag"
	"osdsvr/internal/pkg/config"
	"osdsvr/internal/pkg/server"
	"osdsvr/internal/ssclient"
	"osdsvr/pkg/zlog"
	"runtime"

	"github.com/foss/linksvr/pkg/proto/linkpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	flag.Parse()
	done := make(chan bool)
	// grpc sever
	osdServer := server.NewOsdServer()

	// load init file
	if err := osdServer.LoadInitFile(*config.STORAGE_PATH, *config.OSD_INIT_FILE); err != nil {
		zlog.Error("faile to load init file", zap.Error(err))
	}
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	// register to linksvr
	var conn *grpc.ClientConn
	var err error
	conn, err = grpc.Dial(*config.LINKSVR_GRPC_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zlog.Error("", zap.Error(err))
		done <- true
	}
	defer conn.Close()
	c := linkpb.NewLinkServiceClient(conn)
	linkClient := ssclient.NewLinkService(c)

	go func() {
		err := linkClient.RegisterOsdServiceServer(context.TODO(), *config.OSD_NO, *config.GRPC_ADDR)
		if err != nil {
			zlog.Error("faile to register osd service", zap.Error(err))
		}
		zlog.Info("register to linksvr successfully", zap.Any("osd id", osdServer.OSDID))
	}()
	go func() {
		if err := osdServer.Serve(); err != nil {
			zlog.Error("faile to serve osdsvr", zap.Error(err))
			osdServer.Shutdown()
			done <- true
		}
		zlog.Info("osd grpc server start serving", zap.Any("osd id", osdServer.OSDID))
	}()
	<-done
}
