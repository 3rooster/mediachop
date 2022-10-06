package config

import "fmt"

var (

	// GitCommit will be overwritten automatically by the build system
	GitCommit string
	// GoVersion ver
	GoVersion string
	// BuildTime time
	BuildTime string
)

// Version version string
func Version() string {
	return fmt.Sprintf("Git commit: %6s \nGo version: %6s \nBuild time: %6s \n",
		GitCommit, GoVersion, BuildTime)
}
