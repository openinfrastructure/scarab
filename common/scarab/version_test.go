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
	"testing"
)

func TestVersionSimple(t *testing.T) {
	got, err := version("0.1.0", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if got.Major != 0 {
		t.Errorf("Expected major version 0, got %v", got.Major)
	}
	if got.Minor != 1 {
		t.Errorf("Expected minor version 1, got %v", got.Minor)
	}
	if got.Patch != 0 {
		t.Errorf("Expected patch version 0, got %v", got.Patch)
	}
	if fmt.Sprintf("%v", got) != "0.1.0" {
		t.Errorf("Expected \"0.1.0\", got \"%v\"", got)
	}
}

func TestVersionWithBuild(t *testing.T) {
	got, err := version("0.1.0", "7a145d2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if fmt.Sprintf("%v", got) != "0.1.0+7a145d2" {
		t.Errorf("Expected \"0.1.0+7a145d2\", got \"%v\"", got)
	}
}
