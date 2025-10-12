package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
	GoVersion = runtime.Version()
)

type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
}

func Get() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
	}
}

func String() string {
	info := Get()
	return fmt.Sprintf("obsidian-index version %s\n"+
		"  Git commit: %s\n"+
		"  Build date: %s\n"+
		"  Go version: %s",
		info.Version,
		info.GitCommit,
		info.BuildDate,
		info.GoVersion)
}
