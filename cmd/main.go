package main

import (
	"flag"
	"osdsvr/internal/pkg/server"
	"osdsvr/pkg/zlog"

	"go.uber.org/zap"
)

func main() {
	flag.Parse()
	osdServer := server.NewOsdServer()
	if err := osdServer.Serve(); err != nil {
		zlog.Error("faile to serve osdsvr", zap.Error(err))
		osdServer.Shutdown()
	}
}
