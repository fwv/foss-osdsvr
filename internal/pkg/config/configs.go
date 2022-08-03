package config

import "flag"

var (
	GRPC_ADDR = flag.String("GRPC_ADDR", ":5000", "grpc server listen address")

	SCHEDULE_CAPACITY    = flag.Int64("SCHEDULE_CAPACITY", 5000, "sheduler's task queue size")
	SCHEDULE_CONCURRENCY = flag.Int64("SCHEDULE_CONCURRENCY", 100, "sheduler's worker num")

	STORAGE_PATH        = flag.String("STORAGE_PATH", "/home/fwv/oss/", "object directory")
	DOWNLOAD_CHUNK_SIZE = flag.Int("DOWNLOAD_CHUNK_SIZE", 1024*64, "download chunk size")
)
