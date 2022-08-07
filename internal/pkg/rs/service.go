package rs

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"osdsvr/internal/pkg/core"
	"osdsvr/internal/pkg/object"
	"osdsvr/pkg/zlog"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/klauspost/reedsolomon"
	"go.uber.org/zap"
)

var dataShards = flag.Int("data", 4, "Number of shards to split the data into, must be below 257.")
var parShards = flag.Int("par", 2, "Number of parity shards")
var outDir = flag.String("out", "", "Alternative output directory")
var outFile = flag.String("outFile", "", "Alternative output path/file")

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateRedundantObject(bucketID, hash string, scheduler *core.Scheduler) error {
	zlog.Info("create redundant of object", zap.Any("bucketID", bucketID), zap.Any("object name", hash))
	if scheduler == nil {
		zlog.Error("scheduler is nil when create redundant object")
		return nil
	}
	params := make([]interface{}, 2)
	params[0] = bucketID
	params[1] = hash
	task := &core.Task{
		Param: params,
		DoTask: func(v ...interface{}) error {
			bucketID := v[0].(string)
			hash := v[1].(string)
			if err := s.Encode(bucketID, hash); err != nil {
				zlog.Error("", zap.Error(err))
			}
			return nil
		},
	}
	if err := scheduler.AddTask(task); err != nil {
		return err
	}
	return nil
}

func (s *Service) ReconstructObject(bucketID, hash string) error {
	return s.Decode(bucketID, hash)
}

func (s *Service) Encode(bucketID, hash string) error {
	// Parse command line parameters.
	if (*dataShards + *parShards) > 256 {
		zlog.Error("sum of data and parity shards cannot exceed 256")
		return errors.New("sum of data and parity shards cannot exceed 256")
	}
	oSvc := object.NewService()
	fname := oSvc.GetObjectPath(bucketID, hash)

	// Create encoding matrix.
	enc, err := reedsolomon.NewStream(*dataShards, *parShards)
	checkErr(err)

	// fmt.Println("Opening", fname)
	f, err := os.Open(fname)
	checkErr(err)

	instat, err := f.Stat()
	checkErr(err)

	shards := *dataShards + *parShards
	out := make([]*os.File, shards)

	// Create the resulting files.
	dir, file := filepath.Split(fname)
	if *outDir != "" {
		dir = *outDir
	}
	for i := range out {
		outfn := fmt.Sprintf("%s.%d", file, i)
		out[i], err = os.Create(filepath.Join(dir, outfn))
		checkErr(err)
	}
	zlog.Info("create redundant file data successfully", zap.Any("", strings.Join([]string{fname, ".", "0-", strconv.Itoa(len(out))}, "")))

	// Split into files.
	data := make([]io.Writer, *dataShards)
	for i := range data {
		data[i] = out[i]
	}
	// Do the split
	err = enc.Split(f, data, instat.Size())
	checkErr(err)

	// Close and re-open the files.
	input := make([]io.Reader, *dataShards)

	for i := range data {
		out[i].Close()
		f, err := os.Open(out[i].Name())
		checkErr(err)
		input[i] = f
		defer f.Close()
	}

	// Create parity output writers
	parity := make([]io.Writer, *parShards)
	for i := range parity {
		parity[i] = out[*dataShards+i]
		defer out[*dataShards+i].Close()
	}

	// Encode parity
	err = enc.Encode(input, parity)
	checkErr(err)
	// fmt.Printf("File split into %d data + %d parity shards.\n", *dataShards, *parShards)
	return nil
}

func (s *Service) Decode(bucketID, hash string) error {
	oSvc := object.NewService()
	fname := oSvc.GetObjectPath(bucketID, hash)

	// Create matrix
	enc, err := reedsolomon.NewStream(*dataShards, *parShards)
	checkErr(err)

	// Open the inputs
	shards, size, err := openInput(*dataShards, *parShards, fname)
	checkErr(err)

	// Verify the shards
	ok, err := enc.Verify(shards)
	if !ok {
		zlog.Info("shard data damage is detected! start reconstructing data")
		shards, size, err = openInput(*dataShards, *parShards, fname)
		checkErr(err)
		// Create out destination writers
		out := make([]io.Writer, len(shards))
		for i := range out {
			if shards[i] == nil {
				outfn := fmt.Sprintf("%s.%d", fname, i)
				zlog.Info("restore shard data", zap.Any("shard", outfn))
				out[i], err = os.Create(outfn)
				checkErr(err)
			}
		}
		err = enc.Reconstruct(shards, out)
		if err != nil {
			zlog.Error("reconstruct failed", zap.Error(err))
			return err
		}
		// Close output.
		for i := range out {
			if out[i] != nil {
				err := out[i].(*os.File).Close()
				checkErr(err)
			}
		}
		shards, size, err = openInput(*dataShards, *parShards, fname)
		ok, err = enc.Verify(shards)
		if !ok {
			zlog.Error("verification failed after reconstruction, data likely corrupted", zap.Error(err))
			return err
		}
		checkErr(err)
	}

	// Join the shards and write them
	outfn := *outFile
	if outfn == "" {
		outfn = fname
	}

	f, err := os.Create(outfn)
	checkErr(err)

	shards, size, err = openInput(*dataShards, *parShards, fname)
	checkErr(err)

	// We don't know the exact filesize.
	err = enc.Join(f, shards, int64(*dataShards)*size)
	checkErr(err)
	return nil
}

func openInput(dataShards, parShards int, fname string) (r []io.Reader, size int64, err error) {
	// Create shards and load the data.
	shards := make([]io.Reader, dataShards+parShards)
	for i := range shards {
		infn := fmt.Sprintf("%s.%d", fname, i)
		f, err := os.Open(infn)
		if err != nil {
			zlog.Info("shard data file is lost", zap.Any("shard", infn))
			shards[i] = nil
			continue
		} else {
			shards[i] = f
		}
		stat, err := f.Stat()
		checkErr(err)
		if stat.Size() > 0 {
			size = stat.Size()
		} else {
			shards[i] = nil
		}
	}
	return shards, size, nil
}

func checkErr(err error) {
	if err != nil {
		zlog.Error("")
	}
}
