package server

import (
	"osdsvr/pkg/proto/osdpb"
)

type OsdServer struct {
	osdpb.UnimplementedOsdServiceServer
}

func NewOsdServer() *OsdServer {
	return &OsdServer{}
}

func (s *OsdServer) Shutdown() {
	// todo: gracefully shutdown
}
