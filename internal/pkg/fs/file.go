package fs

import "encoding/binary"

var (
	FILE_SUFFIX_STR string = ".dat"
)

type File interface {
	ReadAt(off int) (r *Record, err error)
	WriteAt(r *Record, off int64) (start int, n int, err error)
	Append(r *Record) (start int, err error)
	Sync() (err error)
	Close() (err error)
	Size() int64
}

type OSDFile struct {
	path       string
	OSDFileID  int64
	writeOff   int64
	io         *FileIO
	actualSize int64
}

func NewOSDFile(path string) (fd *OSDFile, err error) {
	fileIO, err := NewFileIO(path, 0)
	if err != nil {
		return nil, err
	}
	n, _ := fileIO.GetFileSize()
	return &OSDFile{
		path:       path,
		writeOff:   n,
		actualSize: n,
		io:         fileIO,
	}, nil
}

func (f *OSDFile) ReadAt(off int) (r *Record, err error) {
	buf := make([]byte, EntryHeaderSize)
	if _, err := f.io.ReadAt(buf, int64(off)); err != nil {
		return nil, err
	}
	header := readHeader(buf)
	r = &Record{
		header: header,
		data:   &Data{},
	}
	// read value
	off += EntryHeaderSize
	valBuf := make([]byte, header.ValueSize)
	_, err = f.io.ReadAt(valBuf, int64(off))
	if err != nil {
		return nil, err
	}
	r.data.Value = valBuf
	return
}

func (f *OSDFile) WriteAt(r *Record, off int64) (start int, n int, err error) {
	n, err = f.io.WriteAt(r.Encode(), off)
	if err != nil {
		return -1, 0, err
	}
	start = int(f.writeOff)
	f.writeOff += int64(n)
	f.actualSize += int64(n)
	return start, n, nil
}

func (f *OSDFile) Append(r *Record) (start int, err error) {
	s, _, err := f.WriteAt(r, f.writeOff)
	return s, err
}

func (f *OSDFile) Sync() (err error) {
	return f.io.Sync()

}
func (f *OSDFile) Size() int64 {
	return f.actualSize
}

func (f *OSDFile) Close() (err error) {
	return f.io.Close()
}

func readHeader(buf []byte) *Header {
	return &Header{
		Id:        int64(binary.LittleEndian.Uint64(buf[0:8])),
		ValueSize: binary.LittleEndian.Uint32(buf[8:12]),
		Flag:      binary.LittleEndian.Uint16(buf[12:14]),
		Timestamp: int64(binary.LittleEndian.Uint64(buf[14:22])),
	}
}
