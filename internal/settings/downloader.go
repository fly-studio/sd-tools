package settings

import (
	"gopkg.in/go-mixed/go-common.v1/logger.v1"
	httpUtils "gopkg.in/go-mixed/go-common.v1/utils/http"
)

type DownloaderSettings struct {
	Debug bool                 `yaml:"debug"`
	Host  string               `yaml:"host" validate:"required,hostname_port"`
	Log   logger.LoggerOptions `yaml:"log"`

	SD SDOptions `yaml:"sd"`

	Proxy ProxyOptions `yaml:"proxy"`

	MaxTasks       int `yaml:"max_tasks" validate:"required,min=1"`
	ThreadsPerTask int `yaml:"threads_per_task" validate:"required,min=1"`
}

type ProxyOptions struct {
	Address string            `yaml:"address" validate:"url"`
	Domains httpUtils.Domains `yaml:"domain" validate:""`
}
