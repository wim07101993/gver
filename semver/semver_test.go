package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMajorMinorPatch_String(t *testing.T) {
	cases := []struct {
		input  MajorMinorPatch
		output string
	}{
		{
			input:  MajorMinorPatch{},
			output: "0.0.0",
		},
		{
			input:  MajorMinorPatch{Major: 1},
			output: "1.0.0",
		},
		{
			input:  MajorMinorPatch{Minor: 8},
			output: "0.8.0",
		},
		{
			input:  MajorMinorPatch{Patch: 5},
			output: "0.0.5",
		},
		{
			input:  MajorMinorPatch{Major: 1, Minor: 2, Patch: 3},
			output: "1.2.3",
		},
		{
			input:  MajorMinorPatch{Major: 9, Minor: 8, Patch: 7},
			output: "9.8.7",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.input.String(), c.output)
	}
}

func TestSemVer_String(t *testing.T) {
	cases := []struct {
		input  SemVer
		output string
	}{
		{
			input:  SemVer{},
			output: "0.0.0",
		},
		{
			input:  SemVer{MajorMinorPatch: MajorMinorPatch{Major: 1, Minor: 2, Patch: 3}},
			output: "1.2.3",
		},
		{
			input: SemVer{
				MajorMinorPatch: MajorMinorPatch{Major: 9, Minor: 8, Patch: 7},
				Build:           "norelease",
			},
			output: "9.8.7+norelease",
		},
		{
			input: SemVer{
				MajorMinorPatch: MajorMinorPatch{Major: 9, Minor: 8, Patch: 7},
				Prerelease:      "nobuild",
			},
			output: "9.8.7-nobuild",
		},
		{
			input: SemVer{
				MajorMinorPatch: MajorMinorPatch{Major: 9, Minor: 8, Patch: 7},
				Prerelease:      "release",
				Build:           "test",
			},
			output: "9.8.7-release+test",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.input.String(), c.output)
	}
}

func TestParseSemver(t *testing.T) {
	cases := []struct {
		input  string
		output SemVer
	}{
		{
			input:  "0.0.4",
			output: SemVer{MajorMinorPatch: MajorMinorPatch{Patch: 4}},
		},
		{
			input:  "1.2.3",
			output: SemVer{MajorMinorPatch: MajorMinorPatch{Major: 1, Minor: 2, Patch: 3}},
		},
		{
			input:  "10.20.30",
			output: SemVer{MajorMinorPatch: MajorMinorPatch{Major: 10, Minor: 20, Patch: 30}},
		},
		{
			input: "1.1.2-prerelease+meta",
			output: SemVer{
				MajorMinorPatch: MajorMinorPatch{Major: 1, Minor: 1, Patch: 2},
				Prerelease:      "prerelease",
				Build:           "meta",
			},
		},
		{
			input: "1.1.2+meta",
			output: SemVer{
				MajorMinorPatch: MajorMinorPatch{Major: 1, Minor: 1, Patch: 2},
				Build:           "meta",
			},
		},
		{
			input: "1.1.2+meta-valid",
			output: SemVer{
				MajorMinorPatch: MajorMinorPatch{Major: 1, Minor: 1, Patch: 2},
				Build:           "meta-valid",
			},
		},
		{
			input: "1.0.0-alpha",
			output: SemVer{
				MajorMinorPatch: MajorMinorPatch{Major: 1, Minor: 0, Patch: 0},
				Prerelease:      "alpha",
			},
		},
		{
			input: "1.0.0-alpha.beta.1",
			output: SemVer{
				MajorMinorPatch: MajorMinorPatch{Major: 1, Minor: 0, Patch: 0},
				Prerelease:      "alpha.beta.1",
			},
		},
		{
			input: "1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
			output: SemVer{
				MajorMinorPatch: MajorMinorPatch{Major: 1, Minor: 0, Patch: 0},
				Prerelease:      "alpha-a.b-c-somethinglong",
				Build:           "build.1-aef.1-its-okay",
			},
		},
		{
			input: "2.0.0-rc.1+build.123",
			output: SemVer{
				MajorMinorPatch: MajorMinorPatch{Major: 2, Minor: 0, Patch: 0},
				Prerelease:      "rc.1",
				Build:           "build.123"},
		},
	}

	for _, c := range cases {
		semver, err := Parse(c.input)
		assert.Nil(t, err)
		assert.Equal(t, *semver, c.output)
	}
}

func TestSanitizeBuild(t *testing.T) {
	cases := []struct {
		input  string
		output string
	}{
		{"alpha.beta", "alpha.beta"},
		{"alpha-beta", "alpha-beta"},
		{"alpha+beta", "alpha-beta"},
		{"alpha,beta", "alpha-beta"},
		{"..alpha.beta", "alpha.beta"},
		{"alpha.beta..", "alpha.beta"},
	}

	for _, c := range cases {
		output := SanitizeBuild(c.input)
		assert.Equal(t, output, c.output)
	}
}

func TestSanitizePrerelease(t *testing.T) {
	cases := []struct {
		input  string
		output string
	}{
		{"alpha.beta", "alpha.beta"},
		{"alpha-beta", "alpha-beta"},
		{"alpha+beta", "alpha-beta"},
		{"alpha,beta", "alpha-beta"},
		{"..alpha.beta", "alpha.beta"},
		{"alpha.beta..", "alpha.beta"},
	}

	for _, c := range cases {
		output := SanitizePrerelease(c.input)
		assert.Equal(t, output, c.output)
	}
}
