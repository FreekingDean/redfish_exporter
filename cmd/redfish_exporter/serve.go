package main

import (
	"github.com/FreekingDean/redfish_exporter/internal/collectors/chassiscollector"
	"github.com/FreekingDean/redfish_exporter/internal/config"
	"github.com/FreekingDean/redfish_exporter/internal/log"
	"github.com/FreekingDean/redfish_exporter/internal/prometheus"
	"github.com/FreekingDean/redfish_exporter/internal/redfish"
	"github.com/FreekingDean/redfish_exporter/internal/server"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func serve(configFile string) {
	configOptionProvider := func() []config.Option {
		return []config.Option{
			config.WithFilePath(configFile),
		}
	}
	app := fx.New(
		// Initialize FX
		fx.WithLogger(func(log *log.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Zap()}
		}),
		fx.RecoverFromPanics(),

		// Provide Dependencies
		fx.Provide(
			log.New,
			configOptionProvider,
			config.New,
			redfish.New,
			server.NewMux,
			server.New,
			prometheus.NewRegistry,
			chassiscollector.New,
		),

		// Invoke Service
		fx.Invoke(
			redfish.Start,
			prometheus.RegisterBasicCollectors,
			chassiscollector.Register,
			prometheus.RegisterHandler,
			server.Run,
		),
	)

	app.Run()
}

func newServeCmd() *cobra.Command {
	var cfg string
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the Redfish Exporter",
		Run:   func(cmd *cobra.Command, args []string) { serve(cfg) },
	}

	cmd.PersistentFlags().StringVar(&cfg, "config", "./config.yaml", "config file [./config.yaml]")
	return cmd
}
