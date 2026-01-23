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
	sizesFlag       = flag.String("sizes", "128,256,512", "comma-separated list of vector sizes (e.g., 128,256)")
	dirFlag         = flag.String("dir", ".", "directory to process")
	archsimdPfxFlag = flag.String("prefix", "simd", "prefix for the archsimd package")
	midwayPackage   = flag.String("midway", "github.com/dr2chase/midway/midway", "package name for midway helpers")
)

func main() {
	flag.Parse()

	sizesStr := strings.Split(*sizesFlag, ",")
	var sizes []int
	for _, s := range sizesStr {
		var k int
		if _, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &k); err != nil {
			log.Fatalf("invalid size %q: %v", s, err)
		}
		sizes = append(sizes, k)
	}

	if len(sizes) == 0 {
		log.Fatal("no vector sizes specified")
	}

	fmt.Printf("Rewriting for sizes: %v in directory: %s\n", sizes, *dirFlag)

	if err := run(*dirFlag, sizes); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run(dir string, sizes []int) error {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedImports,
		Dir:        dir,
		Env:        append(os.Environ(), "GOOS=linux", "GOARCH=amd64", "GOEXPERIMENT=simd"),
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

		rewriter := NewRewriter(pkg, analyzer, sizes)
		if err := rewriter.Rewrite(); err != nil {
			return fmt.Errorf("rewrite failed: %v", err)
		}
	}

	return nil
}
