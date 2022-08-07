module osdsvr

go 1.16

replace github.com/foss/linksvr => /root/code/foss-linksvr

replace github.com/foss/osdsvr => /root/code/foss-osdsvr

require (
	github.com/foss/linksvr v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.3.0
	github.com/klauspost/reedsolomon v1.10.0
	go.uber.org/zap v1.21.0
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.1
)
