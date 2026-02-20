// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	sizesFlag       = flag.String("sizes", "amd64/linux:128,256,512;wasm/wasip1:128", "semicolon-separated list of arch:size1,size2,etc")
	dirFlag         = flag.String("dir", ".", "directory to process")
	archsimdPfxFlag = flag.String("prefix", "simd", "prefix for the archsimd package")
	midwayPackage   = flag.String("midway", "github.com/dr2chase/midway/midway", "package name for midway helpers")
)

type ArchSizes struct {
	arch, os string
	sizes    []int
}

var _knownArches = []string{"amd64", "arm64", "wasm"}
var knownArches = setFrom(_knownArches)

func setFrom(ss []string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range ss {
		m[s] = true
	}
	return m
}

func main() {
	flag.Parse()

	var allRewrites []ArchSizes

	bySemis := strings.Split(*sizesFlag, ";")
	for _, archSizes := range bySemis {
		as := strings.Split(archSizes, ":")

		if len(as) != 2 {
			log.Fatalf("expected arch:sizes, got %s instead", archSizes)
		}
		os := "linux"
		arch := as[0]
		if slash := strings.Index(arch, "/"); slash != -1 {
			os = arch[slash+1:]
			arch = arch[:slash]
		} else if arch == "wasm" {
			os = "wasip1"
		}
		if !knownArches[arch] {
			log.Fatalf("Expected an architecture in %v, saw %s instead", _knownArches, arch)
		}

		sizesStr := strings.Split(as[1], ",")
		var sizes []int
		for _, s := range sizesStr {
			var k int
			if _, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &k); err != nil {
				log.Fatalf("invalid size %q: %v", s, err)
			}
			sizes = append(sizes, k)
		}
		if len(sizes) == 0 {
			log.Fatalf("no vector sizes specified in %s", archSizes)
		}
		allRewrites = append(allRewrites, ArchSizes{arch: arch, os: os, sizes: sizes})
	}

	for _, rw := range allRewrites {

		fmt.Printf("Rewriting for arch %s and sizes: %v in directory: %s\n", rw.arch, rw.sizes, *dirFlag)

		if err := run(*dirFlag, rw); err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}

func run(dir string, archSizes ArchSizes) error {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedImports,
		Dir:        dir,
		Env:        append(os.Environ(), "GOOS="+archSizes.os, "GOARCH="+archSizes.arch, "GOEXPERIMENT=simd"),
		BuildFlags: []string{"-tags=midway"},
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return fmt.Errorf("loading packages: %v", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		return fmt.Errorf("package load errors")
	}

	// We assume we are processing the package in the directory (usually only one).
	// If there are multiple (e.g. test packages), we might pick the main one or iterate.
	// For now, process all loaded packages that match the directory pattern?
	// Usually `packages.Load` with `.` loads the package in current dir.

	for _, pkg := range pkgs {
		fmt.Printf("Analyzing package: %s\n", pkg.ID)
		analyzer := NewAnalyzer(pkg)
		if err := analyzer.Analyze(); err != nil {
			return fmt.Errorf("analysis failed: %v", err)
		}

		rewriter := NewRewriter(pkg, analyzer, archSizes.arch, archSizes.sizes)
		if err := rewriter.Rewrite(); err != nil {
			return fmt.Errorf("rewrite failed: %v", err)
		}
	}

	return nil
}
