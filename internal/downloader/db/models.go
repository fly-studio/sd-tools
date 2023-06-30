package db

import "gorm.io/gorm"

type DownloadRecord struct {
	gorm.Model
	Url      string `json:"url" gorm:"column:url,unique,not null,size:1024"`
	Path     string `json:"path" gorm:"column:path,not null,size:100"`
	FileName string `json:"file_name" gorm:"column:file_name,not null,size:1024"`
	Status   string `json:"status" gorm:"column:status,not null,size:30"`
	Read     int64  `json:"read" gorm:"column:read"`
	Size     int64  `json:"size" gorm:"column:size"`
	Error    string `json:"error" gorm:"column:error,text"`

	Chunks []*DownloadChunk `json:"-" gorm:"foreignKey:DownloadRecordID;references:ID"`
}

func (m *DownloadRecord) TableName() string {
	return "download_records"
}

type DownloadChunk struct {
	gorm.Model
	DownloadRecordID uint `json:"download_record_id" gorm:"column:download_record_id"`

	Offset int64 `json:"offset" gorm:"column:offset,not null"`
	Size   int64 `json:"size" gorm:"column:size,not null"`
	Read   int64 `json:"read" gorm:"column:read,not null"`
}

func (m *DownloadChunk) TableName() string {
	return "download_chunks"
}
