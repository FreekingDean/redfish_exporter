package chassiscollector

import (
	"fmt"
	"strconv"

	"github.com/FreekingDean/redfish_exporter/internal/collectors"
	"github.com/FreekingDean/redfish_exporter/internal/redfish"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	networkAdapterStateMetric   = "network_adapter_state"
	networkAdapterHealthMetric  = "network_adapter_health"
	networkAdapterTXBytesMetric = "network_adapter_tx_bytes"
	networkAdapterRXBytesMetric = "network_adapter_rx_bytes"

	networkPortStateMetric      = "network_port_state"
	networkPortHealthMetric     = "network_port_health"
	networkPortLinkStatusMetric = "network_port_link_status"
)

var (
	networkAdapterLabels = []string{"network_adapter", "network_adapter_id"}
	networkPortLabels    = []string{"network_port", "network_port_id", "network_port_speed", "network_port_connection_type", "network_port_physical_number"}
)

func networkMetrics() map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		networkAdapterStateMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, networkAdapterStateMetric),
			collectors.StateHelp("chassis.network_adapter"),
			append(labels, networkAdapterLabels...),
			nil,
		),
		networkAdapterHealthMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, networkAdapterHealthMetric),
			collectors.HealthHelp("chassis.network_adapter"),
			append(labels, networkAdapterLabels...),
			nil,
		),
		networkAdapterTXBytesMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, networkAdapterTXBytesMetric),
			"Transmitted bytes of the network adapter",
			append(labels, networkAdapterLabels...),
			nil,
		),
		networkAdapterRXBytesMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, networkAdapterRXBytesMetric),
			"Received bytes of the network adapter",
			append(labels, networkAdapterLabels...),
			nil,
		),
		networkPortStateMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, networkPortStateMetric),
			collectors.StateHelp("chassis.network_port"),
			append(append(labels, networkAdapterLabels...), networkPortLabels...),
			nil,
		),
		networkPortHealthMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, networkPortHealthMetric),
			collectors.HealthHelp("chassis.network_port"),
			append(append(labels, networkAdapterLabels...), networkPortLabels...),
			nil,
		),
		networkPortLinkStatusMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, networkPortLinkStatusMetric),
			"Link status of the network port",
			append(append(labels, networkAdapterLabels...), networkPortLabels...),
			nil,
		),
	}
}

func (c *Collector) collectNetworkMetrics(ch chan<- prometheus.Metric, chassis *redfish.Chassis) {
	c.logger.Debug("Collecting network metrics")
	adapters, err := chassis.NetworkAdapters()
	if err != nil {
		c.logger.Error(fmt.Sprintf("Failed to get network adapter information for chassis %s", chassis.ID), zap.Error(err))
		return
	} else if adapters == nil {
		c.logger.Warn(fmt.Sprintf("No network adapter information for chassis %s", chassis.ID))
		return
	}

	for _, adapter := range adapters {
		labels := []string{"network_adapter", chassis.ID, adapter.Name, adapter.ID}
		if health, ok := collectors.HealthToFloat(adapter.Status.Health); ok {
			ch <- prometheus.MustNewConstMetric(c.metrics[networkAdapterHealthMetric], prometheus.GaugeValue, health, labels...)
		}
		if state, ok := collectors.StateToFloat(adapter.Status.State); ok {
			ch <- prometheus.MustNewConstMetric(c.metrics[networkAdapterStateMetric], prometheus.GaugeValue, state, labels...)
		}
		ch <- prometheus.MustNewConstMetric(c.metrics[networkAdapterTXBytesMetric], prometheus.CounterValue, float64(adapter.Metrics.TXBytes), labels...)
		ch <- prometheus.MustNewConstMetric(c.metrics[networkAdapterRXBytesMetric], prometheus.CounterValue, float64(adapter.Metrics.RXBytes), labels...)

		ports, err := adapter.NetworkPorts()
		if err != nil {
			c.logger.Error(fmt.Sprintf("Failed to get network port information for network adapter %s", adapter.ID), zap.Error(err))
			continue
		} else if ports == nil {
			c.logger.Warn(fmt.Sprintf("No network port information for network adapter %s", adapter.ID))
			continue
		}

		for _, port := range ports {
			portLabels := []string{port.Name, port.ID, strconv.Itoa(port.CurrentLinkSpeedMbps), string(port.ActiveLinkTechnology), port.PhysicalPortNumber}
			portLabels = append(labels, portLabels...)
			if health, ok := collectors.HealthToFloat(port.Status.Health); ok {
				ch <- prometheus.MustNewConstMetric(c.metrics[networkPortHealthMetric], prometheus.GaugeValue, health, portLabels...)
			}
			if state, ok := collectors.StateToFloat(port.Status.State); ok {
				ch <- prometheus.MustNewConstMetric(c.metrics[networkPortStateMetric], prometheus.GaugeValue, state, portLabels...)
			}
			linkStatus := 0.0
			if port.LinkStatus == redfish.NetworkPortLinkStatusUp {
				linkStatus = 1.0
			}
			ch <- prometheus.MustNewConstMetric(c.metrics[networkPortLinkStatusMetric], prometheus.GaugeValue, linkStatus, portLabels...)
		}
	}
}
