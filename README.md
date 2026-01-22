# Midway: SIMD Source-to-Source Rewriter

`midway` is a CLI tool that automatically generates vector-size-specific
SIMD code from Go written to use a "scalable SIMD" API, provided in
`midway/simd/mocks.go`. This is a step towards a more general solution
for writing SIMD code in Go across architectures and hardware with a
variety of vector widths.  The intent is that the generated code for
different vector widths (for example, AVX and AVX-512, or ARM SVE's
scalable vector lengths) will all be present in a single binary that
chooses the right implementation at runtime.  The tool is not yet
general purpose; it is specialized for the API provided in
`midway/simd/mocks.go` and some missing operations need to be filled
in.  

## Overview

The tool operates on Go files marked with `//go:build midway`. It performs two main transformations:
1.  **Dispatcher Generation**: Rewrites the original file in a new `..._simd.go`, replacinmg functions dependent on SIMD types with a switch statement that calls the appropriate architecture-specific implementation based on `midway.MaxVectorSize()`. The build tag is updated to `!midway` so it compiles as standard Go.
2.  **Specialization**: Generates specialized implementation files (e.g., `..._simd128.go`, `..._simd256.go`) where abstract `simd` types are replaced with concrete `archsimd` types (e.g., `simd.Int32s` becomes `archsimd.Int32x4` for 128-bit, `archsimd.Int32x8` for 256-bit).

## Installation

```bash
go install github.com/dr2chase/midway/cmd/midway
```

## Usage

Run the tool on a directory containing `midway`-tagged files:

```bash
midway -dir <directory> -sizes <sizes> [options]
go mod tidy # fill in entry for github.com/dr2chase/midway/midway
```

### Flags

-   `-dir string`: Directory to process (default: current directory `"."`).
-   `-sizes string`: Comma-separated list of vector sizes to generate (e.g., `"128,256,512"`). Default is `"128"`.
-   `-prefix string`: Prefix for the `archsimd` package path (default: `"simd"`).
-   `-midway string`: Package path for midway helpers (default: `"github.com/dr2chase/midway/midway"`).

### Example

Input (`example.go`):
```go
//go:build midway

package example

import "github.com/dr2chase/midway/simd"

func Add(a, b simd.Int32s) simd.Int32s {
	return a // Implementation
}
```

Command:
```bash
midway -dir . -sizes 128,256
```

Output:

1.  **Refactored `example_simd.go`** (Dispatcher):
    ```go
    //go:build !midway

    package example

    import "github.com/dr2chase/midway/midway"

    func Add(a, b simd.Int32s) simd.Int32s {
        switch midway.MaxVectorSize() {
        case 256:
            return Add_simd256(a, b)
        case 128:
            return Add_simd128(a, b)
        default:
            panic("unsupported vector size")
        }
    }
    ```

2.  **Generated `example_simd128.go`**:
    ```go
    //go:build !midway

    package example

    import "simd/archsimd" // concrete types

    func Add_simd128(a, b archsimd.Int32x4) archsimd.Int32x4 {
        return a
    }
    ```

3.  **Generated `example_simd256.go`**:
    Similar to 128-bit, but using `archsimd.Int32x8`.

## Development

The project includes test data in `testdata/simple`
(hand-validated, won't run) and `testdata/ip` (runs, for GOARCH=amd64
and GOEXPERIMENT=simd) demonstrating various usage patterns,
including:
-   Dependent types and aliases.
-   Dependent struct fields.
-   Global variables.
-   Generic functions (instantiated with SIMD types).

## Known issues

Top-level initialized SIMD variables probably don't work yet.

Some of the operations present in simd/mocks.go lack implementations for
particular vector-length/element-type combinations.

