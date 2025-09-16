package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the current version of the application
	// This will be set by build flags: -ldflags "-X github.com/wepala/vine-pod/pkg/version.Version=v1.0.0"
	Version = "dev"

	// Commit is the git commit hash
	// This will be set by build flags: -ldflags "-X github.com/wepala/vine-pod/pkg/version.Commit=$(git rev-parse HEAD)"
	Commit = "unknown"

	// BuildTime is when the binary was built
	// This will be set by build flags: -ldflags "-X github.com/wepala/vine-pod/pkg/version.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
	BuildTime = "unknown"
)

// Info represents version information
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

// Get returns version information
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a string representation of version info
func (i Info) String() string {
	return fmt.Sprintf("Version: %s, Commit: %s, BuildTime: %s, GoVersion: %s, Platform: %s",
		i.Version, i.Commit, i.BuildTime, i.GoVersion, i.Platform)
}
