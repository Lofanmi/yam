package api

import (
	"github.com/uber/jaeger-client-go/config"
)

type TracerConfig struct {
	ServiceName string
	Jaeger      config.Configuration
}
