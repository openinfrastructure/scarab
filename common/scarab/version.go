package scarab

import (
	"fmt"

	"github.com/blang/semver"
)

var (
	// Passed in from build system (Makefile)
	Build        string
	BuildTime    string
	BuildVersion string
	// Set by initializer
	Version semver.Version
)

// Take build strings from the build system and make a semver.Version
func version(ver string, build string) (semver.Version, error) {
	v, err := semver.Make(ver)
	if err != nil {
		return semver.Version{}, err
	}
	if len(build) > 0 {
		v.Build = []string{build}
	}
	return v, nil
}

func init() {
	v, err := version(BuildVersion, Build)
	if err != nil {
		fmt.Println(err)
		return
	}
	Version = v
}
