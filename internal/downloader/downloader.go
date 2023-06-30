package downloader

import (
	"context"
	"github.com/spf13/cobra"
	"gopkg.in/go-mixed/go-common.v1/conf.v1"
	"gopkg.in/go-mixed/go-common.v1/logger.v1"
	"gopkg.in/go-mixed/go-common.v1/web.v1"
	"gopkg.in/go-mixed/go-common.v1/web.v1/controllers"
	"sd-downloader/internal/downloader/db"
	"sd-downloader/internal/settings"
)

func Cmd() *cobra.Command {
	downloadCmd := &cobra.Command{
		Use:   "downloader",
		Short: "A model downloader manager",
		Run: func(cmd *cobra.Command, args []string) {
			config, _ := cmd.Root().PersistentFlags().GetStringSlice("config")
			var settings settings.DownloaderSettings
			if err := conf.LoadSettings(&settings, config...); err != nil {
				panic(err.Error())
			}

			log, err := logger.NewLogger(settings.Log)
			if err != nil {
				panic(err.Error())
			}

			launchDownloader(settings, log)
		},
	}
	return downloadCmd
}

func launchDownloader(settings settings.DownloaderSettings, log *logger.Logger) {
	server := web.NewHttpServer(web.DefaultServerOptions(settings.Host)).SetLogger(log.ToILogger())
	engine := web.NewGinEngine(web.ApiGinOptions(settings.Debug, false), log.ZapLogger())

	c := &DownloaderController{
		Controller: controllers.Controller{},
		log:        log,
		settings:   settings,
		db:         db.DB(),
	}

	registerRoutes(engine, c)

	_ = server.SetDefaultServeHandler(engine, nil)
	log.Info("download manager of stable diffusion is running")

	if err := server.Run(context.Background(), nil); err != nil {
		log.Error(err.Error())
	}
}
