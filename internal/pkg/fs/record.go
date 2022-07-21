package fs

import (
	"encoding/binary"
)

const (
	EntryHeaderSize = 22
)

type R interface {
	Size() int64
	Encode() []byte
	GetValue() []byte
}

type (
	Record struct {
		header *Header
		data   *Data
	}

	Header struct {
		Id        int64
		ValueSize uint32
		Flag      uint16 //有效位，标识set和delete
		Timestamp int64
	}

	Data struct {
		Value []byte
	}
)

func NewRecord(id int64, value []byte, flag uint16, timestamp int64) *Record {
	record := &Record{
		header: &Header{
			Id:        id,
			ValueSize: uint32(len(value)),
			Flag:      flag,
			Timestamp: timestamp,
		},
		data: &Data{
			Value: value,
		},
	}
	return record
}

func (r *Record) Size() int64 {
	return int64(EntryHeaderSize + r.header.ValueSize)
}

//编码顺序为：|--Id(8)--|--ValueSize(4)--|--Flag(2)--|Timestamp(8)--|--Value--|
func (r *Record) Encode() []byte {
	valueSize := r.header.ValueSize

	//set DataItemHeader buf
	buf := make([]byte, r.Size())
	buf = r.setRecordHeaderBuf(buf)
	//set value
	copy(buf[(EntryHeaderSize):(EntryHeaderSize+valueSize)], r.data.Value)
	return buf
}

func (r *Record) setRecordHeaderBuf(buf []byte) []byte {
	binary.LittleEndian.PutUint64(buf[0:8], uint64(r.header.Id))
	binary.LittleEndian.PutUint32(buf[8:12], r.header.ValueSize)
	binary.LittleEndian.PutUint16(buf[12:14], r.header.Flag)
	binary.LittleEndian.PutUint64(buf[14:22], uint64(r.header.Timestamp))
	return buf
}

func (r *Record) GetValue() []byte {
	return r.data.Value
}
