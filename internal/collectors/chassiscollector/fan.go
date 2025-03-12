package chassiscollector

import (
	"strings"

	"github.com/FreekingDean/redfish_exporter/internal/collectors"
	"github.com/FreekingDean/redfish_exporter/internal/redfish"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	fanStateMetric                        = "fan_state"
	fanHealthMetric                       = "fan_health"
	fanRPMMetric                          = "fan_rpm"
	fanRPMPercentageMetric                = "fan_rpm_percentage"
	fanRPMMinMetric                       = "fan_rpm_min"
	fanRPMMaxMetric                       = "fan_rpm_max"
	fanRPMLowerThresholdNonCriticalMetric = "fan_rpm_lower_threshold_non_critical"
	fanRPMLowerThresholdCriticalMetric    = "fan_rpm_lower_threshold_critical"
	fanRPMLowerThresholdFatalMetric       = "fan_rpm_lower_threshold_fatal"
	fanRPMUpperThresholdNonCriticalMetric = "fan_rpm_upper_threshold_non_critical"
	fanRPMUpperThresholdCriticalMetric    = "fan_rpm_upper_threshold_critical"
	fanRPMUpperThresholdFatalMetric       = "fan_rpm_upper_threshold_fatal"
)

var (
	fanLabels = []string{"fan", "fan_id", "fan_unit"}
)

func fanMetrics() map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		fanStateMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, fanStateMetric),
			collectors.StateHelp("chassis.fan"),
			append(labels, fanLabels...),
			nil,
		),
		fanHealthMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, fanHealthMetric),
			collectors.HealthHelp("chassis.fan"),
			append(labels, fanLabels...),
			nil,
		),
		fanRPMMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, fanRPMMetric),
			"RPM of the fan",
			append(labels, fanLabels...),
			nil,
		),
		fanRPMPercentageMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, fanRPMPercentageMetric),
			"Percentage of the fan's RPM compared to the miniumum-maximum RPM",
			append(labels, fanLabels...),
			nil,
		),
		fanRPMMinMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, fanRPMMinMetric),
			"Minimum possible RPM of the fan",
			append(labels, fanLabels...),
			nil,
		),
		fanRPMMaxMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, fanRPMMaxMetric),
			"Maximum possible RPM of the fan",
			append(labels, fanLabels...),
			nil,
		),
		fanRPMLowerThresholdNonCriticalMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, fanRPMLowerThresholdNonCriticalMetric),
			"threshold below the normal range that is not considered critical",
			append(labels, fanLabels...),
			nil,
		),
		fanRPMLowerThresholdCriticalMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, fanRPMLowerThresholdCriticalMetric),
			"threshold below the normal range that is not considered fatal",
			append(labels, fanLabels...),
			nil,
		),
		fanRPMLowerThresholdFatalMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, fanRPMLowerThresholdFatalMetric),
			"threshold below the normal range that is considered fatal",
			append(labels, fanLabels...),
			nil,
		),
		fanRPMUpperThresholdNonCriticalMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, fanRPMUpperThresholdNonCriticalMetric),
			"threshold above the normal range that is not considered critical",
			append(labels, fanLabels...),
			nil,
		),
		fanRPMUpperThresholdCriticalMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, fanRPMUpperThresholdCriticalMetric),
			"threshold above the normal range that is not considered fatal",
			append(labels, fanLabels...),
			nil,
		),
		fanRPMUpperThresholdFatalMetric: prometheus.NewDesc(
			prometheus.BuildFQName(collectors.Namespace, subsystem, fanRPMUpperThresholdFatalMetric),
			"threshold above the normal range that is considered fatal",
			append(labels, fanLabels...),
			nil,
		),
	}
}
func (c *Collector) collectFanMetrics(ch chan<- prometheus.Metric, chassis *redfish.Chassis, thermal *redfish.Thermal) {
	c.logger.Debug("Collecting fan metrics")
	for _, fan := range thermal.Fans {
		labelValues := []string{"fan", chassis.ID, fan.Name, fan.MemberID, strings.ToLower(string(fan.ReadingUnits))}

		if health, ok := collectors.HealthToFloat(fan.Status.Health); ok {
			ch <- prometheus.MustNewConstMetric(c.metrics[fanHealthMetric], prometheus.GaugeValue, health, labelValues...)
		}
		if state, ok := collectors.StateToFloat(fan.Status.State); ok {
			ch <- prometheus.MustNewConstMetric(c.metrics[fanStateMetric], prometheus.GaugeValue, state, labelValues...)
		}
		rpm := float64(fan.Reading)
		percentage := float64(fan.Reading)
		if fan.ReadingUnits == redfish.PercentReadingUnits {
			rpm = rpm * float64(fan.MaxReadingRange) / 100
		} else {
			percentage = (percentage / float64(fan.MaxReadingRange)) * 100
		}

		ch <- prometheus.MustNewConstMetric(c.metrics[fanRPMMetric], prometheus.GaugeValue, float64(rpm), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.metrics[fanRPMPercentageMetric], prometheus.GaugeValue, float64(percentage), labelValues...)

		ch <- prometheus.MustNewConstMetric(c.metrics[fanRPMMinMetric], prometheus.GaugeValue, float64(fan.MinReadingRange), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.metrics[fanRPMMaxMetric], prometheus.GaugeValue, float64(fan.MaxReadingRange), labelValues...)

		ch <- prometheus.MustNewConstMetric(c.metrics[fanRPMLowerThresholdNonCriticalMetric], prometheus.GaugeValue, float64(fan.LowerThresholdNonCritical), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.metrics[fanRPMLowerThresholdCriticalMetric], prometheus.GaugeValue, float64(fan.LowerThresholdCritical), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.metrics[fanRPMLowerThresholdFatalMetric], prometheus.GaugeValue, float64(fan.LowerThresholdFatal), labelValues...)

		ch <- prometheus.MustNewConstMetric(c.metrics[fanRPMUpperThresholdNonCriticalMetric], prometheus.GaugeValue, float64(fan.UpperThresholdNonCritical), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.metrics[fanRPMUpperThresholdCriticalMetric], prometheus.GaugeValue, float64(fan.UpperThresholdCritical), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.metrics[fanRPMUpperThresholdFatalMetric], prometheus.GaugeValue, float64(fan.UpperThresholdFatal), labelValues...)
	}
}
