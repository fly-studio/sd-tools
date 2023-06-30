package db

import (
	"github.com/glebarez/sqlite"
	"gopkg.in/go-mixed/go-common.v1/utils/io"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

func DB() *gorm.DB {
	dir := filepath.Join(ioUtils.GetCurrentDir(), "data")
	_ = os.MkdirAll(dir, 0)

	db, err := gorm.Open(sqlite.Open(filepath.Join(dir, "downloader.db")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err = db.AutoMigrate(&DownloadRecord{}, &DownloadChunk{}); err != nil {
		panic(err)
	}

	return db
}
