package core

import (
	"errors"
	"osdsvr/pkg/proto/osdpb"
	"osdsvr/pkg/zlog"
)

type OSD struct {
	p Processor
}

func (o *OSD) BeforeServe() {}

func (o *OSD) Serve() {}

func (o *OSD) AfterServe() {}

func (o *OSD) UploadFile(file *osdpb.File) (oid int64, err error) {
	if file == nil {
		zlog.Error("file is empty")
		return -1, errors.New("file is empty")
	}

	if o.p == nil {
		zlog.Error("processor is empty")
		return -1, errors.New("processor is empty")
	}
	return o.p.StoreFile(file)
}

func (o *OSD) FetchFile(oid int64) (file *osdpb.File, err error) {
	if o.p == nil {
		zlog.Error("processor is empty")
		return nil, errors.New("processor is empty")
	}
	return o.p.LoadFile(file.MetaData.Oid)
}
