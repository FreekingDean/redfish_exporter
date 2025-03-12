package chassiscollector

import (
	"fmt"

	"github.com/FreekingDean/redfish_exporter/internal/collectors"
	"github.com/FreekingDean/redfish_exporter/internal/redfish"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	powerVoltageStateMetric                    = "power_voltage_state"
	powerVoltageHealthMetric                   = "power_voltage_health"
	powerVoltageVoltsMetric                    = "power_voltage_volts"
	powerAverageConsumedWattsMetric            = "power_average_consumed_watts"
	powerPowerSupplyStateMetric                = "power_power_supply_state"
	powerPowerSupplyHealthMetric               = "power_power_supply_health"
	powerPowerSupplyInputWattsMetric           = "power_power_supply_input_watts"
	powerPowerSupplyOutputWattsMetric          = "power_power_supply_output_watts"
	powerPowerSupplyEfficiencyPercentageMetric = "power_power_supply_efficiency_percentage"
	powerPowerSupplyPowerCapacityWattsMetric   = "power_power_supply_power_capacity_watts"
	powerPowerSupplyLastPowerOutputWattsMetric = "power_power_supply_last_power_output_watts"
)

var (
	powerLabels = []string{"name", "member_id"}
)

func powerMetrics() map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		powerVoltageStateMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, powerVoltageStateMetric),
			collectors.StateHelp("chassis.power_voltage"),
			append(labels, powerLabels...),
			nil,
		),
		powerVoltageHealthMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, powerVoltageHealthMetric),
			collectors.HealthHelp("chassis.power_voltage"),
			append(labels, powerLabels...),
			nil,
		),
		powerVoltageVoltsMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, powerVoltageVoltsMetric),
			"Voltage of the power supply",
			append(labels, powerLabels...),
			nil,
		),
		powerAverageConsumedWattsMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, powerAverageConsumedWattsMetric),
			"Average power consumed in watts",
			append(labels, powerLabels...),
			nil,
		),
		powerPowerSupplyStateMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, powerPowerSupplyStateMetric),
			collectors.StateHelp("chassis.power_supply"),
			append(labels, powerLabels...),
			nil,
		),
		powerPowerSupplyHealthMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, powerPowerSupplyHealthMetric),
			collectors.HealthHelp("chassis.power_supply"),
			append(labels, powerLabels...),
			nil,
		),
		powerPowerSupplyInputWattsMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, powerPowerSupplyInputWattsMetric),
			"Power supply input watts",
			append(labels, powerLabels...),
			nil,
		),
		powerPowerSupplyOutputWattsMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, powerPowerSupplyOutputWattsMetric),
			"Power supply output watts",
			append(labels, powerLabels...),
			nil,
		),
		powerPowerSupplyEfficiencyPercentageMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, powerPowerSupplyEfficiencyPercentageMetric),
			"Power supply efficiency percentage",
			append(labels, powerLabels...),
			nil,
		),
		powerPowerSupplyPowerCapacityWattsMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, powerPowerSupplyPowerCapacityWattsMetric),
			"Power supply power capacity watts",
			append(labels, powerLabels...),
			nil,
		),
		powerPowerSupplyLastPowerOutputWattsMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, powerPowerSupplyLastPowerOutputWattsMetric),
			"Power supply last power output watts",
			append(labels, powerLabels...),
			nil,
		),
	}
}

func (c *Collector) collectPowerMetrics(ch chan<- prometheus.Metric, chassis *redfish.Chassis) {
	c.logger.Debug("Collecting power metrics")
	power, err := chassis.Power()
	if err != nil {
		c.logger.Error(fmt.Sprintf("Failed to get power information for chassis %s", chassis.ID), zap.Error(err))
		return
	} else if power == nil {
		c.logger.Warn(fmt.Sprintf("No power information for chassis %s", chassis.ID))
		return
	}

	for _, voltage := range power.Voltages {
		labelValues := []string{"power_voltage", chassis.ID, voltage.Name, voltage.MemberID}
		if state, ok := collectors.StateToFloat(voltage.Status.State); ok {
			ch <- prometheus.MustNewConstMetric(c.metrics[powerVoltageStateMetric], prometheus.GaugeValue, state, labelValues...)
		}
		if health, ok := collectors.HealthToFloat(voltage.Status.Health); ok {
			ch <- prometheus.MustNewConstMetric(c.metrics[powerVoltageHealthMetric], prometheus.GaugeValue, health, labelValues...)
		}

		ch <- prometheus.MustNewConstMetric(c.metrics[powerVoltageVoltsMetric], prometheus.GaugeValue, float64(voltage.ReadingVolts), labelValues...)
	}

	for _, powerControl := range power.PowerControl {
		labelValues := []string{"power_control", chassis.ID, powerControl.Name, powerControl.MemberID}
		ch <- prometheus.MustNewConstMetric(c.metrics[powerAverageConsumedWattsMetric], prometheus.GaugeValue, float64(powerControl.PowerMetrics.AverageConsumedWatts), labelValues...)
	}

	for _, powerSupply := range power.PowerSupplies {
		labelValues := []string{"power_supply", chassis.ID, powerSupply.Name, powerSupply.MemberID}
		if state, ok := collectors.StateToFloat(powerSupply.Status.State); ok {
			ch <- prometheus.MustNewConstMetric(c.metrics[powerPowerSupplyStateMetric], prometheus.GaugeValue, state, labelValues...)
		}
		if health, ok := collectors.HealthToFloat(powerSupply.Status.Health); ok {
			ch <- prometheus.MustNewConstMetric(c.metrics[powerPowerSupplyHealthMetric], prometheus.GaugeValue, health, labelValues...)
		}
		ch <- prometheus.MustNewConstMetric(c.metrics[powerPowerSupplyInputWattsMetric], prometheus.GaugeValue, float64(powerSupply.PowerInputWatts), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.metrics[powerPowerSupplyOutputWattsMetric], prometheus.GaugeValue, float64(powerSupply.PowerOutputWatts), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.metrics[powerPowerSupplyEfficiencyPercentageMetric], prometheus.GaugeValue, float64(powerSupply.EfficiencyPercent), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.metrics[powerPowerSupplyPowerCapacityWattsMetric], prometheus.GaugeValue, float64(powerSupply.PowerCapacityWatts), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.metrics[powerPowerSupplyLastPowerOutputWattsMetric], prometheus.GaugeValue, float64(powerSupply.LastPowerOutputWatts), labelValues...)
	}
}
