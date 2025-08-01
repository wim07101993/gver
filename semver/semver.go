package semver

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const Regex = `^(?P<Major>0|[1-9]\d*)\.(?P<Minor>0|[1-9]\d*)\.(?P<Patch>0|[1-9]\d*)(?:-(?P<Prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<Build>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

var verregexp *regexp.Regexp

func init() {
	var err error
	verregexp, err = regexp.Compile(Regex)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to compile semver regex"))
	}
}

type SemVer struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Build      string
}

func (v *SemVer) MajorMinorPatch() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *SemVer) Full() string {
	builder := strings.Builder{}

	builder.WriteString(v.MajorMinorPatch())

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
	match := verregexp.FindStringSubmatch(s)
	if len(match) == 0 {
		return nil, errors.New("failed to parse SemVer according to regex")
	}

	version := SemVer{}
	for i, name := range verregexp.SubexpNames() {
		if i != 0 && name != "" {
			switch name {
			case "Major":
				version.Major, _ = strconv.Atoi(match[i])
			case "Minor":
				version.Minor, _ = strconv.Atoi(match[i])
			case "Patch":
				version.Patch, _ = strconv.Atoi(match[i])
			case "Prerelease":
				version.Prerelease = match[i]
			case "Build":
				version.Build = match[i]
			}
		}
	}

	return &version, nil
}
