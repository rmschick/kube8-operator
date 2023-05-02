package types

import (
	"github.com/FishtechCSOC/common-go/pkg/build"
)

const (
	cyderesSuffix    = "cyderes.io"
	programSeparator = "@"
	versionSeparator = "/"
)

func CreateUserAgent() string {
	return build.Program + programSeparator + cyderesSuffix + versionSeparator + build.Build
}
