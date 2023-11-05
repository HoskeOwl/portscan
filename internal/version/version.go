package version

import "fmt"

var (
	Version        = "dev"
	CommitHash     = "null"
	BuildTimestamp = "null"
)

func BuildVersion() string {
	return fmt.Sprintf("%s-%s (%s)", Version, CommitHash, BuildTimestamp)
}
