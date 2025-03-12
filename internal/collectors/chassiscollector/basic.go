package chassiscollector

import (
	"fmt"

	"github.com/FreekingDean/redfish_exporter/internal/collectors"
	"github.com/FreekingDean/redfish_exporter/internal/redfish"
	"github.com/prometheus/client_golang/prometheus"
)

func basicChassisMetrics() map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		healthMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, "health"),
			collectors.HealthHelp("chassis"),
			labels,
			nil, // Constant Labels
		),
		stateMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, "state"),
			fmt.Sprintf("state of chassis,%s", collectors.CommonStateHelp),
			labels,
			nil, // Constant Labels
		),
		modelInfoMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, "model_info"),
			"organization responsible for producing the chassis, the name by which the manufacturer generally refers to the chassis, and a part number and sku assigned by the organization that is responsible for producing or manufacturing the chassis",
			append(labels, modelLabels...),
			nil, // Constant Labels
		),
	}
}

func (c *Collector) collectBasicMetrics(ch chan<- prometheus.Metric, chassis *redfish.Chassis) {
	c.logger.Debug("Collecting basic chassis metrics")
	labels := []string{chassis.ID, chassis.Name}
	if health, ok := collectors.HealthToFloat(chassis.Status.Health); ok {
		ch <- prometheus.MustNewConstMetric(c.metrics[healthMetric], prometheus.GaugeValue, health, labels...)
	}
	if state, ok := collectors.StateToFloat(chassis.Status.State); ok {
		ch <- prometheus.MustNewConstMetric(c.metrics[stateMetric], prometheus.GaugeValue, state, labels...)
	}
}
