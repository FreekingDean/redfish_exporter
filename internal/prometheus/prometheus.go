package prometheus

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

func NewRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}

func RegisterHandler(mux *http.ServeMux, reg *prometheus.Registry, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
			return nil
		},
	})
}

func RegisterBasicCollectors(reg *prometheus.Registry, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := reg.Register(
				collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			)
			if err != nil {
				return err
			}

			return reg.Register(collectors.NewGoCollector())
		},
	})
}
