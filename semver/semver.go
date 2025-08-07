package semver

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var semverRegex = regexp.MustCompile(`^(?P<Major>0|[1-9]\d*)\.(?P<Minor>0|[1-9]\d*)\.(?P<Patch>0|[1-9]\d*)(?:-(?P<Prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<Build>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
var allowedBuildCharsRegex = regexp.MustCompile(`[^A-Za-z0-9-.]`)
var allowedPrereleaseCharsRegex = regexp.MustCompile(`[^A-Za-z0-9-.]`)

type SemVer struct {
	MajorMinorPatch MajorMinorPatch
	Prerelease      string
	Build           string
}

func (v *SemVer) AddMajorMinorPatch(mmp MajorMinorPatch) {
	v.MajorMinorPatch.Major += mmp.Major
	v.MajorMinorPatch.Minor += mmp.Minor
	v.MajorMinorPatch.Patch += mmp.Patch
}

func (v *SemVer) String() string {
	builder := strings.Builder{}

	builder.WriteString(v.MajorMinorPatch.String())

	if v.Prerelease != "" {
		builder.WriteString("-")
		builder.WriteString(v.Prerelease)
	}

	if v.Build != "" {
		builder.WriteString("+")
		builder.WriteString(v.Build)
	}

	return builder.String()
}

func Parse(s string) (*SemVer, error) {
	match := semverRegex.FindStringSubmatch(s)
	if len(match) == 0 {
		return nil, errors.New("failed to parse SemVer according to regex")
	}

	version := SemVer{}
	for i, name := range semverRegex.SubexpNames() {
		if i != 0 && name != "" {
			switch name {
			case "Major":
				version.MajorMinorPatch.Major, _ = strconv.ParseUint(match[i], 10, 64)
			case "Minor":
				version.MajorMinorPatch.Minor, _ = strconv.ParseUint(match[i], 10, 64)
			case "Patch":
				version.MajorMinorPatch.Patch, _ = strconv.ParseUint(match[i], 10, 64)
			case "Prerelease":
				version.Prerelease = match[i]
			case "Build":
				version.Build = match[i]
			}
		}
	}

	return &version, nil
}

func SanitizeBuild(s string) string {
	s = strings.Trim(s, ".")
	return allowedBuildCharsRegex.ReplaceAllString(s, "-")
}

func SanitizePrerelease(s string) string {
	s = strings.Trim(s, ".")
	return allowedPrereleaseCharsRegex.ReplaceAllString(s, "-")
}
