package chassiscollector

import (
	"fmt"

	"github.com/FreekingDean/redfish_exporter/internal/collectors"
	"github.com/FreekingDean/redfish_exporter/internal/redfish"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	tempSensorStateMetric  = "temperature_sensor_state"
	tempSensorHealthMetric = "temperature_sensor_health"
	tempSensorTempMetric   = "temperature_celsius"
)

var (
	tempSensorLabels = []string{"sensor", "sensor_id"}
)

func thermalChassisMetrics() map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		tempSensorStateMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, tempSensorStateMetric),
			collectors.StateHelp("chassis.temprature_sensor"),
			append(labels, tempSensorLabels...),
			nil,
		),
		tempSensorHealthMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, tempSensorHealthMetric),
			collectors.HealthHelp("chassis.temprature_sensor"),
			append(labels, tempSensorLabels...),
			nil,
		),
		tempSensorTempMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, tempSensorTempMetric),
			"celcius temperature of the chassis component",
			append(labels, tempSensorLabels...),
			nil,
		),
	}
}

func (c *Collector) collectThermalMetrics(ch chan<- prometheus.Metric, chassis *redfish.Chassis) {
	thermal, err := chassis.Thermal()
	if err != nil {
		c.logger.Error(fmt.Sprintf("Failed to get thermal information for chassis %s", chassis.ID), zap.Error(err))
		return
	} else if thermal == nil {
		c.logger.Warn(fmt.Sprintf("No thermal information for chassis %s", chassis.ID))
		return
	}

	for _, tempSensor := range thermal.Temperatures {
		labelValues := []string{"temperature", chassis.ID, thermal.Name, tempSensor.MemberID}
		c.logger.Debug(fmt.Sprintf("Collecting thermal sensor metrics for %s", tempSensor.MemberID))
		if health, ok := collectors.HealthToFloat(tempSensor.Status.Health); ok {
			ch <- prometheus.MustNewConstMetric(c.metrics[tempSensorHealthMetric], prometheus.GaugeValue, health, labelValues...)
		}
		if state, ok := collectors.StateToFloat(tempSensor.Status.State); ok {
			ch <- prometheus.MustNewConstMetric(c.metrics[tempSensorStateMetric], prometheus.GaugeValue, state, labelValues...)
		}
		ch <- prometheus.MustNewConstMetric(c.metrics[tempSensorTempMetric], prometheus.GaugeValue, float64(tempSensor.ReadingCelsius), labelValues...)
	}
	c.collectFanMetrics(ch, chassis, thermal)
}
