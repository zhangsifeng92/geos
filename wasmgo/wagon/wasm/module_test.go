// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/zhangsifeng92/geos/wasmgo/wagon/wasm"
)

var testPaths = []string{
	"testdata",
	"../exec/testdata",
	"../exec/testdata/spec",
}

func TestReadModule(t *testing.T) {
	for _, dir := range testPaths {
		fnames, err := filepath.Glob(filepath.Join(dir, "*.wasm"))
		if err != nil {
			t.Fatal(err)
		}
		for _, fname := range fnames {
			name := fname
			t.Run(filepath.Base(name), func(t *testing.T) {
				raw, err := ioutil.ReadFile(name)
				if err != nil {
					t.Fatal(err)
				}

				r := bytes.NewReader(raw)
				m, err := wasm.ReadModule(r, nil)
				if err != nil {
					t.Fatalf("error reading module %v", err)
				}
				if m == nil {
					t.Fatalf("error reading module: (nil *Module)")
				}
			})
		}
	}
}
