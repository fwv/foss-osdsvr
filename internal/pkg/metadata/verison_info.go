package metadata

import "encoding/json"

type VersionInfo struct {
	Version     int64  `json:"verison"`
	Hash        string `json:"objectName"`
	Size        int64  `json:"size"`
	Upload_time int64  `json:"uploadTime"`
}

func CreateVersionInfo(record *VersionRecord) *VersionInfo {
	return &VersionInfo{
		Version:     record.Version,
		Hash:        record.Hash,
		Size:        record.Size,
		Upload_time: record.UploadTime,
	}
}

func (v *VersionInfo) ToJsonString() string {
	data, _ := json.Marshal(v)
	return string(data)
}
