package config

import "fmt"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func AppVersionInfo() (s string) {
	s = fmt.Sprintf("version %v, commit %v, built at %v", version, commit, date)
	return
}
