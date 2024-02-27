package buildinfo

import (
	"runtime"

	"github.com/fatih/color"
)

// -X with build
var (
	MainVersion string
	GoVersion   = runtime.Version()
	GoOSArch    = runtime.GOOS + "/" + runtime.GOARCH
	GitSha      string
	BuildTime   string
)

var (
	bold     = color.New(color.Bold)
	bluebold = color.New(color.FgBlue, color.Bold)
)

func Version() string {
	s1 := bluebold.Sprintf("Version: ")
	s2 := bold.Sprintf("%s %s (commit-id=%s)", AppName, MainVersion, GitSha)
	s3 := bluebold.Sprintf("Runtime: ")
	s4 := bold.Sprintf("%s %s RELEASE.%s", GoVersion, GoOSArch, BuildTime)
	return s1 + s2 + "\r\n" + s3 + s4
}
