package service

import (
	"fmt"
)

var (
	BuildVersion  string
	BuildCommitID string
)

func getBuildInfo() string {
	return fmt.Sprintf("Version: %s, Commit: %s", BuildVersion, BuildCommitID)
}
