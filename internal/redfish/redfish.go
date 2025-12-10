package redfish

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/FreekingDean/redfish_exporter/internal/config"
	"github.com/FreekingDean/redfish_exporter/internal/log"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	HealthOK       = common.OKHealth
	HealthCritical = common.CriticalHealth
	HealthWarning  = common.WarningHealth

	StateEnabled            = common.EnabledState
	StateDisabled           = common.DisabledState
	StateStandbyOffline     = common.StandbyOfflineState
	StateStandbySpare       = common.StandbySpareState
	StateInTest             = common.InTestState
	StateStarting           = common.StartingState
	StateAbsent             = common.AbsentState
	StateUnavailableOffline = common.UnavailableOfflineState
	StateDeferring          = common.DeferringState
	StateQuiesced           = common.QuiescedState
	StateUpdating           = common.UpdatingState

	PercentReadingUnits = redfish.PercentReadingUnits

	NetworkPortLinkStatusUp   = redfish.UpPortLinkStatus
	NetworkPortLinkStatusDown = redfish.DownPortLinkStatus
)

type (
	Chassis = redfish.Chassis
	Thermal = redfish.Thermal
	Health  = common.Health
	State   = common.State
)

func NewClientConfig(cfg config.Config) *gofish.ClientConfig {
	defaultTransport := http.DefaultTransport.(*http.Transport)
	transport := &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   time.Duration(10) * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			// Manually Added additional CipherSuites to support TLS 1.0
			CipherSuites: []uint16{
				// TLS 1.0 - 1.2 cipher suites.
				tls.TLS_RSA_WITH_RC4_128_SHA,
				tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
				tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				// TLS 1.3 cipher suites.
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
			},
		},
	}

	config := gofish.ClientConfig{
		Endpoint:   fmt.Sprintf("https://%s", cfg.Host.Endpoint),
		Username:   cfg.Host.Username,
		Password:   cfg.Host.Password,
		BasicAuth:  cfg.Host.BasicAuth,
		Insecure:   true,
		HTTPClient: &http.Client{Transport: transport},
	}

	return &config
}

type Client struct {
	*gofish.APIClient
}

func NewClient(logger *log.Logger, clientConfig *gofish.ClientConfig) (*Client, error) {
	logger.Debug("Connecting to redfish service", zap.String("endpoint", clientConfig.Endpoint))

	client, err := gofish.Connect(*clientConfig)
	if err != nil {
		logger.Error("Failed to connect to redfish service", zap.String("endpoint", clientConfig.Endpoint), zap.Error(err))
		return nil, err
	}

	return &Client{client}, nil
}

func Start(client *Client, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			client.APIClient.Logout()
			return nil
		},
	})
}
