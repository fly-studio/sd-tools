package downloader

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/go-mixed/go-common.v1/web.v1/controllers"
	"sd-downloader/internal/downloader/db"
)

func registerRoutes(r *gin.Engine, c *DownloaderController) {
	g := r.Group("/api/v1/downloader")
	{
		// list all download tasks
		g.GET("/", controllers.Handle[[]*db.DownloadRecord](c, c.List))

		// add a new download task
		g.POST("/download", controllers.Handle[*db.DownloadRecord](c, c.Download))

		// get a download task
		g.GET("/:id", controllers.Handle[*db.DownloadRecord](c, c.Get))

		// delete download tasks
		g.DELETE("/", controllers.Handle[[]*db.DownloadRecord](c, c.Delete))

		// pause download tasks
		g.PUT("/pause", controllers.Handle[[]*db.DownloadRecord](c, c.Pause))

		// resume download tasks
		g.PUT("/resume", controllers.Handle[[]*db.DownloadRecord](c, c.Resume))

		g.PUT("/retry", controllers.Handle[[]*db.DownloadRecord](c, c.Retry))
	}
}
