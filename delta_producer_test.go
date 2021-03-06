//
// Copyright © 2017-2020 Solus Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package libeopkg

import (
	"os"
	"testing"
)

// Our test files, known to produce a valid delta package
const (
	deltaOldPkg = "testdata/delta/nano-4.6-117-1-x86_64.eopkg"
	deltaNewPkg = "testdata/delta/nano-4.7-118-1-x86_64.eopkg"
	notAFile    = "testdata/bob"
	notAPkg     = "testdata/not.xml"
)

func TestBasicDelta(t *testing.T) {
	producer, err := NewDeltaProducer("TESTING", deltaOldPkg, deltaNewPkg)
	if err != nil {
		t.Fatalf("Failed to create delta producer for existing pkgs: %v", err)
	}
	defer producer.Close()
	path, err := producer.Create()
	if err != nil {
		t.Fatalf("Failed to produce delta packages: %v", err)
	}
	defer os.Remove(path)

	pkg, err := Open(path)
	if err != nil {
		t.Fatalf("Failed to open our delta package: %v", err)
	}
	defer pkg.Close()
	if err = pkg.ReadAll(); err != nil {
		t.Fatalf("Failed to read metadata on delta: %v", err)
	}
	if pkg.Meta.Package.Name != "nano" {
		t.Fatalf("Invalid delta name: %s", pkg.Meta.Package.Name)
	}
	if pkg.Meta.Package.GetRelease() != 118 {
		t.Fatalf("Invalid release number in delta: %d", pkg.Meta.Package.GetRelease())
	}
}

func TestBasicDeltaOldMissing(t *testing.T) {
	producer, err := NewDeltaProducer("TESTING", notAFile, deltaNewPkg)
	if err == nil {
		t.Fatalf("Should have failed to create delta producer for non-existent pkg: %s", notAFile)
	}
	producer.Close()
}

func TestBasicDeltaOldInvalid(t *testing.T) {
	producer, err := NewDeltaProducer("TESTING", notAPkg, deltaNewPkg)
	if err == nil {
		t.Fatalf("Should have failed to create delta producer for invalid pkg: %s", notAPkg)
	}
	producer.Close()
}

func TestBasicDeltaNewMissing(t *testing.T) {
	producer, err := NewDeltaProducer("TESTING", deltaOldPkg, notAFile)
	if err == nil {
		t.Fatalf("Should have failed to create delta producer for non-existent pkg: %s", notAFile)
	}
	producer.Close()
}

func TestBasicDeltaNewInvalid(t *testing.T) {
	producer, err := NewDeltaProducer("TESTING", deltaOldPkg, notAPkg)
	if err == nil {
		t.Fatalf("Should have failed to create delta producer for invalid pkg: %s", notAPkg)
	}
	producer.Close()
}

func TestBasicDeltaImpossibleEqual(t *testing.T) {
	producer, err := NewDeltaProducer("TESTING", deltaOldPkg, deltaOldPkg)
	if err == nil {
		t.Fatalf("Should have failed to create delta producer for identical pkg: %s", deltaOldPkg)
	}
	producer.Close()
}

func TestBasicDeltaImpossibleGreater(t *testing.T) {
	producer, err := NewDeltaProducer("TESTING", deltaNewPkg, deltaOldPkg)
	if err == nil {
		t.Fatalf("Should have failed to create delta producer for newer old pkg: %s", deltaNewPkg)
	}
	producer.Close()
}
