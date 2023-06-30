package downloader

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gopkg.in/go-mixed/go-common.v1/logger.v1"
	"gopkg.in/go-mixed/go-common.v1/web.v1/controllers"
	"gorm.io/gorm"
	"path/filepath"
	"sd-downloader/internal/downloader/db"
	"sd-downloader/internal/settings"
)

type DownloaderController struct {
	controllers.Controller
	settings settings.DownloaderSettings
	log      *logger.Logger
	db       *gorm.DB
}

type ListRequest struct {
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Status string `json:"status"`
	Search string `json:"search"`
}

// List all download tasks, support pagination
func (c *DownloaderController) List(ctx *gin.Context) ([]*db.DownloadRecord, error) {
	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, errors.WithStack(err)
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Size <= 0 {
		req.Size = 20
	}

	var records []*db.DownloadRecord
	query := c.db.Model(&db.DownloadRecord{}).Order("id desc").Limit(req.Size).Offset((req.Page - 1) * req.Size)
	if req.Search != "" {
		query = query.Where(
			c.db.Where("url like ?", "%"+req.Search+"%").
				Or("path like ?", "%"+req.Search+"%").
				Or("file_name like ?", "%"+req.Search+"%"),
		)
	}
	if req.Status != DownloadUnknown {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Find(&records).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return records, nil
}

type DownloadRequest struct {
	Url      string `json:"url" validate:"required"`
	Path     string `json:"path" validate:"required"`
	FileName string `json:"file_name" validate:"required"`
}

// Download a file from DownloadRequest
func (c *DownloaderController) Download(ctx *gin.Context) (*db.DownloadRecord, error) {
	var req DownloadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, errors.WithStack(err)
	}

	if !c.settings.SD.ValidPath(req.Path) {
		return nil, errors.New("invalid path")
	}

	req.Path = filepath.Join(c.settings.SD.Path, req.Path)

	var record db.DownloadRecord
	if res := c.db.Where("url = ?", req.Url).First(&record); res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) { // DB error
			return nil, errors.WithStack(res.Error)
		}
	} else { // url already exists
		return &record, errors.Errorf("download task already exists. url=%s, id=%d", req.Url, record.ID)
	}

	// create a new download task
	record = db.DownloadRecord{
		Url:      req.Url,
		Path:     req.Path,
		FileName: req.FileName,
		Status:   DownloadPending,
	}

	if res := c.db.Create(&record); res.Error != nil {
		return nil, errors.WithStack(res.Error)
	}

	return &record, nil
}

// Get a download task by id
func (c *DownloaderController) Get(ctx *gin.Context) (*db.DownloadRecord, error) {
	id := ctx.Param("id")

	var record db.DownloadRecord
	if res := c.db.First(&record, id); res.Error != nil {
		return nil, errors.WithStack(res.Error)
	}

	return &record, nil
}

type PauseRequest struct {
	IDs []int64 `json:"ids" validate:"required"`
}

// Pause download tasks
func (c *DownloaderController) Pause(ctx *gin.Context) ([]*db.DownloadRecord, error) {
	var req PauseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, errors.WithStack(err)
	}

	records, err := c.getRecords(req.IDs)
	if err != nil {
		return nil, err
	}

	return records, nil
}

type ResumeRequest struct {
	IDs []int64 `json:"ids" validate:"required"`
}

// Resume paused download tasks
func (c *DownloaderController) Resume(ctx *gin.Context) ([]*db.DownloadRecord, error) {
	var req ResumeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, errors.WithStack(err)
	}

	records, err := c.getRecords(req.IDs)
	if err != nil {
		return nil, err
	}

	return records, nil
}

type DeleteRequest struct {
	IDs      []int64 `json:"ids" validate:"required"`
	WithFile bool    `json:"with_file"`
}

// Delete download tasks
func (c *DownloaderController) Delete(ctx *gin.Context) ([]*db.DownloadRecord, error) {
	var req DeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, errors.WithStack(err)
	}

	records, err := c.getRecords(req.IDs)
	if err != nil {
		return nil, err
	}

	c.db.Delete(records)

	return records, nil
}

type RetryRequest struct {
	IDs []int64 `json:"ids" validate:"required"`
}

// Retry failed download tasks
func (c *DownloaderController) Retry(ctx *gin.Context) ([]*db.DownloadRecord, error) {
	var req DeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, errors.WithStack(err)
	}

	records, err := c.getRecords(req.IDs)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (c *DownloaderController) getRecords(ids []int64) ([]*db.DownloadRecord, error) {
	var records []*db.DownloadRecord
	if res := c.db.Find(&records, ids); res.Error != nil {
		return nil, errors.WithStack(res.Error)
	}

	return records, nil
}
