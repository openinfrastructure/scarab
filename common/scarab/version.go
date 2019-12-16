/*
Copyright © 2019 Open Infrastructure Services, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scarab

import (
	"fmt"

	"github.com/blang/semver"
)

var (
	// Build* variables are set by the Makefile via LDFLAGS

	// Build e.g. `"ac9d7bd"`
	Build string
	// BuildTime e.g. `"2019-12-16T21:32:36Z"`
	BuildTime string
	// BuildVersion e.g. `"0.1.0"`
	BuildVersion string

	// Version is computed by the initalizer from the Build* values
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
