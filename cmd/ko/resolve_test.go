// Copyright 2018 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/google/go-containerregistry/pkg/ko/build"
	"github.com/google/go-containerregistry/pkg/v1"
)

func TestDefaultGoBuildOptionsCnt(t *testing.T) {
	opts, err := gobuildOptions()
	if err != nil {
		t.Errorf("gobuildOptions() = %v", err)
	}

	// WithBaseImages, WithCreationTime
	if len(opts) != 2 {
		t.Fatalf("opt cnt = %d want %d", len(opts), 2)
	}
}

func TestDefaultGoBuildImageCreationTime(t *testing.T) {
	baseLayers := int64(3)
	base, err := random.Image(1024, baseLayers)
	if err != nil {
		t.Fatalf("random.Image() = %v", err)
	}
	importPath := "github.com/google/go-containerregistry"

	opts, err := gobuildOptions()
	if err != nil {
		t.Errorf("gobuildOptions() = %v", err)
	}

	// hack here
	opts[0] = build.WithBaseImages(func(string) (v1.Image, error) { return base, nil })
	ng, err := build.NewGo(opts...)

	img, err := ng.Build(filepath.Join(importPath, "cmd", "ko", "test"))
	if err != nil {
		t.Fatalf("Build() = %v", err)
	}

	cfg, err := img.ConfigFile()
	if err != nil {
		t.Errorf("ConfigFile() = %v", err)
	}

	actual := cfg.Created
	empty := v1.Time{}
	if actual.Time == empty.Time {
		t.Fatalf("created = %v, want current time", empty)
	}
}
