package collectors

import (
	"github.com/FreekingDean/redfish_exporter/internal/redfish"
)

func HealthToFloat(health redfish.Health) (float64, bool) {
	switch health {
	case redfish.HealthOK:
		return float64(1), true
	case redfish.HealthWarning:
		return float64(2), true
	case redfish.HealthCritical:
		return float64(3), true
	}

	return float64(0), false
}

func StateToFloat(state redfish.State) (float64, bool) {
	switch state {
	case redfish.StateEnabled:
		return float64(1), true
	case redfish.StateDisabled:
		return float64(2), true
	case redfish.StateStandbyOffline:
		return float64(3), true
	case redfish.StateStandbySpare:
		return float64(4), true
	case redfish.StateInTest:
		return float64(5), true
	case redfish.StateStarting:
		return float64(6), true
	case redfish.StateAbsent:
		return float64(7), true
	case redfish.StateUnavailableOffline:
		return float64(8), true
	case redfish.StateDeferring:
		return float64(9), true
	case redfish.StateQuiesced:
		return float64(10), true
	case redfish.StateUpdating:
		return float64(11), true
	}
	return float64(0), false
}

//func parseCommonPowerState(status redfish.PowerState) (float64, bool) {
//	if bytes.Equal([]byte(status), []byte("On")) {
//		return float64(1), true
//	} else if bytes.Equal([]byte(status), []byte("Off")) {
//		return float64(2), true
//	} else if bytes.Equal([]byte(status), []byte("PoweringOn")) {
//		return float64(3), true
//	} else if bytes.Equal([]byte(status), []byte("PoweringOff")) {
//		return float64(4), true
//	}
//	return float64(0), false
//}
//
//func parseLinkStatus(status redfish.LinkStatus) (float64, bool) {
//	if bytes.Equal([]byte(status), []byte("LinkUp")) {
//		return float64(1), true
//	} else if bytes.Equal([]byte(status), []byte("NoLink")) {
//		return float64(2), true
//	} else if bytes.Equal([]byte(status), []byte("LinkDown")) {
//		return float64(3), true
//	}
//	return float64(0), false
//}
//
//func parsePortLinkStatus(status redfish.NetworkPortLinkStatus) (float64, bool) {
//	if bytes.Equal([]byte(status), []byte("Up")) {
//		return float64(1), true
//	}
//	return float64(0), false
//}
//func boolToFloat64(data bool) float64 {
//
//	if data {
//		return float64(1)
//	}
//	return float64(0)
//
//}
//
//func parsePhySecIntrusionSensor(method redfish.IntrusionSensor) (float64, bool) {
//	if bytes.Equal([]byte(method), []byte("Normal")) {
//		return float64(1), true
//	}
//	if bytes.Equal([]byte(method), []byte("TamperingDetected")) {
//		return float64(2), true
//	}
//	if bytes.Equal([]byte(method), []byte("HardwareIntrusion")) {
//		return float64(3), true
//	}
//
//	return float64(0), false
//}
