package object

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/url"
	"os"
	"osdsvr/internal/pkg/config"
	"osdsvr/internal/pkg/fs"
	"osdsvr/pkg/proto/osdpb"
	"osdsvr/pkg/zlog"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) WriteObject(objPath string, stream osdpb.OsdService_UploadFileServer) (string, error) {
	writeoff := 0
	sha := sha256.New()
	for {
		req, err := stream.Recv()
		// handle EOF
		if err == io.EOF {
			hash := base64.StdEncoding.EncodeToString(sha.Sum(nil))
			hash = url.PathEscape(hash)
			// zlog.Debug("accept file completed", zap.String("object save path", objPath), zap.Int("content size", writeoff))
			return hash, stream.SendAndClose(&osdpb.FileUploadResponse{
				ResultCode: osdpb.Result_SUCCESS,
				Desc:       "upload file sucessfully",
			})
		}
		// handle error
		if err != nil {
			zlog.Error("failed to receive chunk", zap.Error(err))
			return "", stream.SendAndClose(&osdpb.FileUploadResponse{
				ResultCode: osdpb.Result_FAILED,
				Desc:       "receive chunk failed",
			})
		}
		var file *os.File
		if file == nil {
			file, err = os.OpenFile(objPath, os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				return "", err
			}
			defer file.Close()
		}
		n, err := file.WriteAt(req.Chunk, int64(writeoff))
		if err != nil {
			return "", err
		}
		// zlog.Debug("write chunk data to file completed", zap.Int("chunk size", n))

		if _, err := sha.Write(req.Chunk); err != nil {
			zlog.Error("failed to calculate hash sha256")
			return "", err
		}
		writeoff += n
		// if err := file.Sync(); err != nil {
		// 	return err
		// }
	}
}

func (s *Service) GetObject(hash string, bucketID string, stream osdpb.OsdService_DownloadFileServer) error {
	// read data
	objectPath := s.GetObjectPath(bucketID, hash)
	file, err := os.OpenFile(objectPath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	readoff := 0
	data := make([]byte, *config.DOWNLOAD_CHUNK_SIZE)
	for {
		n, err := file.ReadAt(data, int64(readoff))
		readoff += n
		if n != 0 {
			zlog.Debug("send chunk data to stream", zap.Int("chunk size", n))
			if err := stream.Send(&osdpb.FileDownloadResponse{
				Chunk: data[:n],
			}); err != nil {
				zlog.Error("failed to send chunk data to stream", zap.Error(err))
				return err
			}
		}
		if err == io.EOF {
			zlog.Info("send file data to stream completed", zap.String("object name", hash), zap.Int("total size", readoff))
			break
		} else if err != nil {
			zlog.Error("failed to read file", zap.Error(err))
			return err
		}
	}
	return nil
}

func (s *Service) RenameObject(oldObjPath, newObjPath string) error {
	if err := os.Rename(oldObjPath, newObjPath); err != nil {
		zlog.Error("failed to rename object", zap.Error(err))
		return err
	}
	// zlog.Debug("rename object sucessfully", zap.Any("tmp object", oldObjPath), zap.Any("new object", newObjPath))
	return nil
}

func (s *Service) CheckObject() error {
	return nil
}

func (s *Service) GetObjectTmpPath(objectName string) string {
	rid := uuid.New().String()
	tmpDir := strings.Join([]string{*config.STORAGE_PATH, "tmp/"}, "")
	fs.CreatePathIfNotExists(tmpDir)
	object := strings.Join([]string{objectName, "tmp", rid}, "-")
	return strings.Join([]string{tmpDir, object}, "")
}

func (s *Service) GetObjectPath(bucketID string, hash string) string {
	dir := strings.Join([]string{*config.STORAGE_PATH, bucketID, "/"}, "")
	fs.CreatePathIfNotExists(dir)
	path := strings.Join([]string{dir, hash}, "")
	return path
}
