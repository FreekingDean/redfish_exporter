package version

import "fmt"

var (
	major int
	minor int
	patch int
)

// Version is the current version of the application.
func Version() string {
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
