package config

import "flag"

var (
	OSD_NO            = flag.Int64("OSD_NO", 1, "no of osd")
	GRPC_ADDR         = flag.String("GRPC_ADDR", ":5000", "grpc server listen address")
	LINKSVR_GRPC_ADDR = flag.String("LINKSVR_GRPC_ADDR", ":5001", "linksvr grpc server listen address")

	SCHEDULE_CAPACITY    = flag.Int64("SCHEDULE_CAPACITY", 5000, "sheduler's task queue size")
	SCHEDULE_CONCURRENCY = flag.Int64("SCHEDULE_CONCURRENCY", 50, "sheduler's worker num")

	STORAGE_PATH        = flag.String("STORAGE_PATH", "/root/oss/", "object directory")
	OSD_INIT_FILE       = flag.String("OSD_INIT_FILE", "osd.init", "osd sever init file")
	DOWNLOAD_CHUNK_SIZE = flag.Int("DOWNLOAD_CHUNK_SIZE", 1024*64, "download chunk size")

	REGISTER_RETRY_COUNT = flag.Int("REGISTER_RETRY_COUNT ", 5, "register retry count")

	RS_CODE_MODE = flag.Bool("RS_CODE_MODE", true, "reed solomon mode")
)
