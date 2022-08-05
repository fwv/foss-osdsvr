package ssclient

import (
	"context"

	"github.com/foss/linksvr/pkg/proto/linkpb"
)

type LinkSerivce struct {
	linkClient linkpb.LinkServiceClient
}

func NewLinkService(linkClient linkpb.LinkServiceClient) *LinkSerivce {
	return &LinkSerivce{
		linkClient: linkClient,
	}
}

func (s *LinkSerivce) RegisterOsdServiceServer(ctx context.Context, osdID int64, addr string) error {
	req := &linkpb.RegisterOSDRequest{
		Addr:  addr,
		OsdId: osdID,
	}
	_, err := s.linkClient.RegisterOSD(context.TODO(), req)
	if err != nil {
		return err
	}
	return nil
}
