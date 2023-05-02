/*
Package build is just a dummy store for build time information. These values should be overwritten during build time by
ldflags for metadata used throughout the program.
*/
package build

// nolint:gochecknoglobals
var (
	// Build represents the version of code in the binary.
	Build = "snapshot"
	// Commit SHA of the code used in the binary.
	Commit = "none"
	// Date the binary was built.
	Date = "unknown"
	// Version of go this binary was built from.
	Version = "unknown"
	// Program is the name of the binary/program.
	Program = "common-go"
	// OS is the OS the binary was built for.
	OS = "unknown"
	// Architecture represents the processor architecture the binary was built for.
	Architecture = "unknown"
	// ARM is only used for the ARM architecture, delineates the version of ARM the binary was built for.
	ARM = ""
)
