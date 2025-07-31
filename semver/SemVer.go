package semver

import (
	"fmt"
	"strings"
)

type SemVer struct {
	Major   int
	Minor   int
	Patch   int
	Release string
	Build   string
}

func (v *SemVer) MajorMinorPatch() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *SemVer) Full() string {
	builder := strings.Builder{}

	builder.WriteString(v.MajorMinorPatch())

	if v.Release != "" {
		builder.WriteString("-")
		builder.WriteString(v.Release)
	}

	if v.Build != "" {
		builder.WriteString("+")
		builder.WriteString(v.Build)
	}

	return builder.String()
}
