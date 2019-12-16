package scarab

import (
  "fmt"

  "github.com/blang/semver"
)

var (
  // Passed in from build system (Makefile)
  Build string
  BuildTime string
  BuildVersion string
  // Set by initializer
  Version semver.Version
)

func init() {
  v, err := semver.Make(BuildVersion)
  if err != nil {
    fmt.Println(err)
  }
  v.Build = []string{Build}
  Version = v
}
