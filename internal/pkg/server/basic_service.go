package server

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"osdsvr/internal/pkg/config"
	"osdsvr/internal/pkg/core"
	"osdsvr/internal/pkg/metadata"
	"osdsvr/internal/pkg/object"
	"osdsvr/internal/pkg/rs"
	"osdsvr/internal/ssclient"
	"osdsvr/pkg/proto/osdpb"
	"osdsvr/pkg/zlog"

	"github.com/foss/linksvr/pkg/proto/linkpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type OsdServer struct {
	osdpb.UnimplementedOsdServiceServer
	core.OSDInitializer
	OSDID      int64
	scheduler  *core.Scheduler
	oService   *object.Service
	mService   *metadata.Service
	rsService  *rs.Service
	LinkClient *ssclient.LinkSerivce
}

func NewOsdServer() *OsdServer {
	s := &OsdServer{
		OSDID:     *config.OSD_NO,
		scheduler: core.NewScheduler(*config.SCHEDULE_CAPACITY),
		oService:  object.NewService(),
		mService:  metadata.NewService(),
		rsService: rs.NewService(),
	}
	return s
}

func (s *OsdServer) LoadInitFile(path string, initFileName string) error {
	s.OSDInitializer.LoadInitFile(path, initFileName)
	// read init from init file
	if s.InitFile == nil {
		zlog.Error("init file is nil", zap.String("path", path), zap.String("file name", initFileName))
		return errors.New("init file is nill")
	}
	// skip read for the first time
	n, _ := s.InitFile.GetFileSize()
	if n == 0 {
		s.SyncInitFile()
		return nil
	}

	data, err := s.InitFile.ReadAllBytes()
	if err != nil {
		return err
	}

	init := &core.OSDServerInit{}
	proto.Unmarshal(data, init)
	s.OSDID = init.ID
	zlog.Info("load init file for OSDSaver successfully", zap.String("path", path), zap.Int64("OSD ID", init.ID))
	return nil
}

func (s *OsdServer) SyncInitFile() error {
	if s.InitFile != nil {
		init := &core.OSDServerInit{
			ID: s.OSDID,
		}
		initData, err := proto.Marshal(init)
		if err != nil {
			return err
		}
		os.Truncate(s.GetInitFilePath(), 0)
		_, err = s.InitFile.WriteAt(initData, 0)
		if err != nil {
			return err
		}
		zlog.Info("sync init file sucessfully", zap.Any("init file", s.GetInitFilePath()), zap.Any("OSD ID", s.OSDID))
	}
	return nil
}

func (s *OsdServer) Serve() error {
	// init
	if err := s.LoadInitFile(*config.STORAGE_PATH, *config.OSD_INIT_FILE); err != nil {
		zlog.Fatal("failed to init", zap.Error(err))
		return err
	}
	// register osd
	conn, err := grpc.Dial(*config.LINKSVR_GRPC_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return err
	}
	// defer conn.Close()
	c := linkpb.NewLinkServiceClient(conn)
	s.LinkClient = ssclient.NewLinkService(c)
	if err := s.LinkClient.RegisterOsdServiceServer(context.TODO(), s.OSDID, *config.GRPC_ADDR); err != nil {
		zlog.Fatal("failed to register osdsvr", zap.Error(err))
		return err
	}

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
