package collectors

import "fmt"

const (
	CommonHealthHelp   = "1(OK),2(Warning),3(Critical)"
	CommonSeverityHelp = CommonHealthHelp
	CommonStateHelp    = "1(Enabled),2(Disabled),3(StandbyOffinline),4(StandbySpare),5(InTest),6(Starting),7(Absent),8(UnavailableOffline),9(Deferring),10(Quiesced),11(Updating)"
)

func HealthHelp(component string) string {
	return fmt.Sprintf("health of %s,%s", component, CommonHealthHelp)
}

func StateHelp(component string) string {
	return fmt.Sprintf("state of %s,%s", component, CommonStateHelp)
}
