package server

import (
	"net/http"

	"github.com/FreekingDean/redfish_exporter/internal/config"
)

func Run(cfg config.Config) error {
	_ = &http.Server{Addr: cfg.Web.ListenAddress()}
	return nil
	//return web.ListenAndServe(server, cfg.PrometheusConfig(), kitlogger)
}
