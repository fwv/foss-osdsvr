package core

import (
	"os"
	"osdsvr/pkg/zlog"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSaver(t *testing.T) {
	path := "/home/fwv/code/foss-osdsvr/testing/osd_file/"
	os.RemoveAll(path)
	// test store
	data := []byte("nihaonihao") // 10 byte
	saver := NewOSDSaver(path, 70)
	entry, err := saver.Store(data)
	zlog.Info("", zap.Any("entry", entry))
	assert.NoError(t, err)
	assert.Equal(t, int64(1), entry.fid)
	assert.Equal(t, int64(0), entry.loc)

	entry, err = saver.Store(data)
	zlog.Info("", zap.Any("entry", entry))
	assert.NoError(t, err)
	assert.Equal(t, int64(1), entry.fid)
	assert.Equal(t, int64(32), entry.loc)

	entry, err = saver.Store(data)
	zlog.Info("", zap.Any("entry", entry))
	assert.NoError(t, err)
	assert.Equal(t, int64(2), entry.fid)
	assert.Equal(t, int64(0), entry.loc)

	// test load
	readData, err := saver.Load(&Entry{fid: 1, loc: 0})
	assert.NoError(t, err)
	assert.Equal(t, "nihaonihao", string(readData))

	readData, err = saver.Load(&Entry{fid: 1, loc: 32})
	assert.NoError(t, err)
	assert.Equal(t, "nihaonihao", string(readData))

	readData, err = saver.Load(&Entry{fid: 2, loc: 0})
	assert.NoError(t, err)
	assert.Equal(t, "nihaonihao", string(readData))
}

func TestInitializer(t *testing.T) {
	path := "/home/fwv/code/foss-osdsvr/testing/osd_file/"
	os.RemoveAll(path)
	// test store
	data := []byte("nihaonihao") // 10 byte
	saver := NewOSDSaver(path, 40)
	entry, err := saver.Store(data)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), entry.fid)
	assert.Equal(t, int64(0), entry.loc)

	entry, err = saver.Store(data)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), entry.fid)
	assert.Equal(t, int64(0), entry.loc)

	saver1 := NewOSDSaver(path, 40)
	entry, err = saver1.Store(data)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), entry.fid)
	assert.Equal(t, int64(0), entry.loc)

}

// goos: linux
// goarch: amd64
// pkg: osdsvr/internal/pkg/core
// cpu: AMD Ryzen 5 4600H with Radeon Graphics
// BenchmarkStore-12
// 8673	    143049 ns/op
func BenchmarkStore(b *testing.B) {
	path := "/home/fwv/code/foss-osdsvr/testing/bench_store/"
	os.RemoveAll(path)
	data10B := []byte("nihaonihao") // 10 byte
	data64KB := make([]byte, 0)
	for i := 0; i < 6400; i++ {
		data64KB = append(data64KB, data10B...)
	}
	zlog.Info("", zap.Any("len of 64kB data", len(data64KB)))
	saver := NewOSDSaver(path, 1024*1024*4)
	for i := 0; i < b.N; i++ {
		saver.Store(data64KB)
	}
}
