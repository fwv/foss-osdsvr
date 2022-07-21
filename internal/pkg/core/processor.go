package core

import (
	"osdsvr/pkg/proto/osdpb"

	"google.golang.org/protobuf/proto"
)

type Processor interface {
	StoreFile(file *osdpb.File) (oid int64, err error)
	LoadFile(oid int64) (file *osdpb.File, err error)
}

type OSDObj struct {
	oid  int64
	file *osdpb.File
}

type OSDComponents struct {
	saver Saver
	index Index
}

type OSDProcessor struct {
	generator     IDGenerator
	fileComponent OSDComponents
	metaComponent OSDComponents
}

func (p *OSDProcessor) StoreFile(file *osdpb.File) (oid int64, err error) {
	// generate oid for file
	oid, err = p.generator.GenerateOid()
	if err != nil {
		return -1, err
	}
	o := &OSDObj{
		oid:  0,
		file: file,
	}
	// store object
	err = p.storeOSDObj(o)
	if err != nil {
		return -1, err
	}
	return oid, nil
}

func (p *OSDProcessor) LoadFile(oid int64) (file *osdpb.File, err error) {
	return nil, nil
}

func (p *OSDProcessor) storeOSDObj(o *OSDObj) error {
	// store file content
	if err := p.fileComponent.saver.AsyncStore(o.file.Content, p.fileComponent.index.Put); err != nil {
		return err
	}

	// store metadata content
	metaData, err := proto.Marshal(o.file.MetaData)
	if err != nil {
		return err
	}
	if err := p.metaComponent.saver.AsyncStore(metaData, p.metaComponent.index.Put); err != nil {
		return err
	}
	return nil
}
