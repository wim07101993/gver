package semver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSemVer_MajorMinorPatch(t *testing.T) {
	cases := []struct {
		input  SemVer
		output string
	}{
		{
			input:  SemVer{},
			output: "0.0.0",
		},
		{
			input:  SemVer{Major: 1},
			output: "1.0.0",
		},
		{
			input:  SemVer{Minor: 8},
			output: "0.8.0",
		},
		{
			input:  SemVer{Patch: 5},
			output: "0.0.5",
		},
		{
			input:  SemVer{Major: 1, Minor: 2, Patch: 3},
			output: "1.2.3",
		},
		{
			input:  SemVer{Major: 9, Minor: 8, Patch: 7, Prerelease: "release", Build: "test"},
			output: "9.8.7",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.input.MajorMinorPatch(), c.output)
	}
}

func TestSemVer_Full(t *testing.T) {
	cases := []struct {
		input  SemVer
		output string
	}{
		{
			input:  SemVer{},
			output: "0.0.0",
		},
		{
			input:  SemVer{Major: 1, Minor: 2, Patch: 3},
			output: "1.2.3",
		},
		{
			input:  SemVer{Major: 9, Minor: 8, Patch: 7, Build: "norelease"},
			output: "9.8.7+norelease",
		},
		{
			input:  SemVer{Major: 9, Minor: 8, Patch: 7, Prerelease: "nobuild"},
			output: "9.8.7-nobuild",
		},
		{
			input:  SemVer{Major: 9, Minor: 8, Patch: 7, Prerelease: "release", Build: "test"},
			output: "9.8.7-release+test",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.input.Full(), c.output)
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		input  string
		output SemVer
	}{
		{
			input:  "0.0.4",
			output: SemVer{Patch: 4},
		},
		{
			input:  "1.2.3",
			output: SemVer{Major: 1, Minor: 2, Patch: 3},
		},
		{
			input:  "10.20.30",
			output: SemVer{Major: 10, Minor: 20, Patch: 30},
		},
		{
			input:  "1.1.2-prerelease+meta",
			output: SemVer{Major: 1, Minor: 1, Patch: 2, Prerelease: "prerelease", Build: "meta"},
		},
		{
			input:  "1.1.2+meta",
			output: SemVer{Major: 1, Minor: 1, Patch: 2, Build: "meta"},
		},
		{
			input:  "1.1.2+meta-valid",
			output: SemVer{Major: 1, Minor: 1, Patch: 2, Build: "meta-valid"},
		},
		{
			input:  "1.0.0-alpha",
			output: SemVer{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha"},
		},
		{
			input:  "1.0.0-alpha.beta.1",
			output: SemVer{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.beta.1"},
		},
		{
			input:  "1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
			output: SemVer{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha-a.b-c-somethinglong", Build: "build.1-aef.1-its-okay"},
		},
		{
			input:  "2.0.0-rc.1+build.123",
			output: SemVer{Major: 2, Minor: 0, Patch: 0, Prerelease: "rc.1", Build: "build.123"},
		},
	}

	for _, c := range cases {
		semver, err := Parse(c.input)
		assert.Nil(t, err)
		assert.Equal(t, *semver, c.output)
	}
}
