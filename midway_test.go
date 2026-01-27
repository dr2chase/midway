// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestSimdOpsMatchesMocks verifies that the output of cmd/simdops matches the contents of simd/mocks.go
func TestSimdOpsMatchesMocks(t *testing.T) {
	// Run go run cmd/simdops/analyze_simd_ops.go
	cmd := exec.Command("go", "run", "cmd/simdops/analyze_simd_ops.go")
	if testing.Verbose() {
		t.Logf("go run cmd/simdops/analyze_simd_ops.go > simd/mocks.go-generated # real command in a tmpdir")
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run cmd/simdops: %v\nOutput:\n%s", err, output)
	}

	// Read simd/mocks.go
	mocksPath := filepath.Join("simd", "mocks.go")
	expected, err := os.ReadFile(mocksPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", mocksPath, err)
	}

	// Compare (trimming spaces/newlines slightly might be needed if go run adds output, but usually it shouldn't for stdout)
	// normalize newlines just in case
	if !bytes.Equal(output, expected) {
		t.Errorf("simd/mocks.go does not match output of cmd/simdops")
		// Save output for debugging if mismatch
		_ = os.WriteFile("simd/mocks.go-generated", output, 0644)
		t.Logf("Wrote generated output to simd/mocks.go-generated for comparison")
	}
}

// TestMidwaySimpleGeneration verifies that running cmd/midway in testdata/simple will generate files "test0_simd*.go" that match their counterparts in testdata/simd/golden
func TestMidwaySimpleGeneration(t *testing.T) {
	// Create a temp dir
	tmpDir := t.TempDir()

	// Copy testdata/simple to tmpDir
	simpleDir := filepath.Join("testdata", "simple")
	if err := copyDir(simpleDir, tmpDir); err != nil {
		t.Fatalf("failed to copy %s to %s: %v", simpleDir, tmpDir, err)
	}

	// Initialize module for packages.Load to work
	initModule(t, tmpDir, "simple")

	// Run cmd/midway
	// We run it using "go run ./cmd/midway"
	// The -dir flag points to the directory to process
	cmd := exec.Command("go", "run", "./cmd/midway", "-dir", tmpDir)
	if testing.Verbose() {
		t.Logf("go run ./cmd/midway -dir testdata/simple # real command in a tmpdir")
	}
	// cmd/midway logs to stderr/stdout
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to run cmd/midway: %v\nOutput:\n%s", err, output)
	}

	// Verify generated files match golden files
	// Golden files are in testdata/simple/golden (based on our research)
	// The user mentioned "testdata/simd/golden" in the prompt, but we found "testdata/simple/golden".
	// We will try finding "testdata/simd/golden" if it exists, otherwise "testdata/simple/golden".
	goldenDir := filepath.Join("testdata", "simd", "golden")
	if _, err := os.Stat(goldenDir); os.IsNotExist(err) {
		goldenDir = filepath.Join("testdata", "simple", "golden")
	}

	filesToCheck := []string{"test0_simd.go", "test0_simd128.go", "test0_simd256.go", "test0_simd512.go"}

	for _, fname := range filesToCheck {
		generatedPath := filepath.Join(tmpDir, fname)
		goldenPath := filepath.Join(goldenDir, fname)

		genContent, err := os.ReadFile(generatedPath)
		if err != nil {
			t.Errorf("failed to read generated file %s: %v", fname, err)
			continue
		}

		goldenContent, err := os.ReadFile(goldenPath)
		if err != nil {
			t.Errorf("failed to read golden file %s: %v", goldenPath, err)
			continue
		}

		if !bytes.Equal(genContent, goldenContent) {
			t.Errorf("content mismatch for %s", fname)
			_ = os.WriteFile("testdata/simple/"+fname+"-generated", genContent, 0644)
			t.Logf("Wrote generated output to testdata/simple/%s-generated for comparison", fname)
		}
	}
}

func TestMidwayIpCompilation(t *testing.T) {
	testMidwayCompilation(t, "ip")
}

func TestMidwaySplitPkgCompilation(t *testing.T) {
	testMidwayCompilation(t, "splitpkg")
}

// testMidwayCompilation verifies that running cmd/midway in testdata/subdir
// generates files that will compile (with GOARCH=amd64 and GOEXPERIMENT=simd)
func testMidwayCompilation(t *testing.T, subdir string) {
	tmpDir := t.TempDir()
	dir := filepath.Join("testdata", subdir)

	// Copy subdir dir contents to tmpDir
	if err := copyDir(dir, tmpDir); err != nil {
		t.Fatalf("failed to copy %s to %s: %v", dir, tmpDir, err)
	}

	// Initialize module for packages.Load to work
	initModule(t, tmpDir, subdir)

	// Run cmd/midway
	cmd := exec.Command("go", "run", "./cmd/midway", "-dir", tmpDir)
	if testing.Verbose() {
		t.Logf("go run ./cmd/midway -dir testdata/%s # real command in a tmpdir", subdir)
	}
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to run cmd/midway: %v\nOutput:\n%s", err, output)
	}

	// Compile with GOARCH=amd64 GOEXPERIMENT=simd go build .
	buildCmd := exec.Command("go", "build", ".")
	buildCmd.Dir = tmpDir
	buildCmd.Env = append(os.Environ(), "GOARCH=amd64", "GOEXPERIMENT=simd", "CGO_ENABLED=0")
	if testing.Verbose() {
		t.Logf("( cd testdata/%s; GOARCH=amd64 GOEXPERIMENT=simd go build . ) # real command in a tmpdir", subdir)
	}
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to compile generated code in %s: %v\nOutput:\n%s", tmpDir, err, output)
	}
}

// Helper to copy a directory recursively (simple version for testdata)
func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// initModule initializes a module in the given temporary directory
// and adds a replace directive to point to the local copy of midway.
func initModule(t *testing.T, dir, modName string) {
	cmd := exec.Command("go", "mod", "init", modName)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to init module in %s: %v\nOutput:\n%s", dir, err, output)
	}

	// Add replace directive to point to local midway
	// Assuming the test runs from the root of midway repo, we can use ".." relative path or absolute path.
	// We'll use absolute path to be safe.
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	// If we are running in src/midway, cwd is the root.
	// replace github.com/dr2chase/midway => cwd
	cmdReplace := exec.Command("go", "mod", "edit", "-replace", "github.com/dr2chase/midway="+cwd)
	cmdReplace.Dir = dir
	if output, err := cmdReplace.CombinedOutput(); err != nil {
		t.Fatalf("failed to add replace directive in %s: %v\nOutput:\n%s", dir, err, output)
	}

	// Run go mod tidy to resolve dependencies using the replacement
	cmdTidy := exec.Command("go", "mod", "tidy")
	cmdTidy.Dir = dir
	if output, err := cmdTidy.CombinedOutput(); err != nil {
		t.Fatalf("failed to run go mod tidy in %s: %v\nOutput:\n%s", dir, err, output)
	}
}
