package chassiscollector

import (
	"context"
	"sync"

	"github.com/FreekingDean/redfish_exporter/internal/collectors"
	"github.com/FreekingDean/redfish_exporter/internal/log"
	"github.com/FreekingDean/redfish_exporter/internal/redfish"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

const (
	subsystem = "chassis"
)

var (
	labels          = []string{"resource", "chassis_id"}
	modelLabels     = []string{"manufacturer", "model", "part_number", "sku"}
	healthMetric    = "health"
	stateMetric     = "state"
	modelInfoMetric = "model_info"
)

type collectorFunc func(chan<- prometheus.Metric, *redfish.Chassis)

type Collector struct {
	logger         *log.Logger
	redfish        *redfish.Client
	metrics        map[string]*prometheus.Desc
	scrapeStatus   *prometheus.GaugeVec
	collectorFuncs []collectorFunc
}

func New(logger *log.Logger, client *redfish.Client) *Collector {
	return &Collector{
		logger:  logger,
		redfish: client,
		metrics: make(map[string]*prometheus.Desc),
		scrapeStatus: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: collectors.Namespace,
				Name:      "collector_scrape_status",
				Help:      "collector_scrape_status",
			},
			[]string{"collector"},
		),
	}
}

func Register(collector *Collector, registry *prometheus.Registry, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			metricGroups := []map[string]*prometheus.Desc{
				basicChassisMetrics(),
				thermalChassisMetrics(),
				fanMetrics(),
				powerMetrics(),
				networkMetrics(),
			}
			for _, metrics := range metricGroups {
				for metricName, metric := range metrics {
					collector.metrics[metricName] = metric
				}
			}
			collector.collectorFuncs = []collectorFunc{
				collector.collectBasicMetrics,
				collector.collectThermalMetrics,
				collector.collectPowerMetrics,
				collector.collectNetworkMetrics,
			}

			return registry.Register(collector)
		},
	})
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric
	}
	c.scrapeStatus.Describe(ch)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.logger.Debug("Collecting chassis metrics")

	chassiss, err := c.redfish.GetService().Chassis()
	if err != nil {
		c.logger.Error("Failed to get chassis", log.Error(err))
		c.scrapeStatus.WithLabelValues("chassis").Set(float64(0))
		return
	}

	wg := sync.WaitGroup{}
	for _, chassis := range chassiss {
		for _, collectorFunc := range c.collectorFuncs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				collectorFunc(ch, chassis)
			}()
		}
	}
	wg.Wait()

	c.logger.Debug("Finished collecting chassis metrics")
	c.scrapeStatus.WithLabelValues("chassis").Set(float64(1))
}
