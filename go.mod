module osdsvr

go 1.18

require (
	github.com/foss/linksvr v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/stretchr/testify v1.8.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)

replace github.com/foss/linksvr => /home/fwv/code/foss-linksvr
