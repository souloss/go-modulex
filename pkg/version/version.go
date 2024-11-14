package version

import "runtime"

var (
	ProjectName  string = "ModuleX"
	Version      string = "dev"
	BuildTime    string
	GitBranch    string
	GitRevision  string
	PlatformName string = runtime.GOOS + "/" + runtime.GOARCH
	GOVersion    string = runtime.Version()
)
