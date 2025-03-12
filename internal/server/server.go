package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FreekingDean/redfish_exporter/internal/config"
	"github.com/FreekingDean/redfish_exporter/internal/log"
	"github.com/prometheus/exporter-toolkit/web"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New(mux *http.ServeMux) *http.Server {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	return &http.Server{
		Handler: mux,
	}
}

func Run(server *http.Server, cfg config.Config, logger *log.Logger, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addresses := []string{cfg.Web.ListenAddress()}

			flagConfig := &web.FlagConfig{
				WebListenAddresses: &addresses,
				WebConfigFile:      &cfg.Web.ConfigFile,
			}

			go func() {
				if err := web.ListenAndServe(server, flagConfig, logger.Slog()); err != nil {
					if err != http.ErrServerClosed {
						logger.Error("Failed to start server", zap.Error(err))
					}
				}
			}()
			return nil
		},

		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping application...")
			return server.Shutdown(ctx)
		},
	})
}
