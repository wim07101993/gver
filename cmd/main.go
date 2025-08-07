package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gver"
	"log"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/pkg/errors"
)

const (
	SemverFormat          = "semver"
	MajorMinorPatchFormat = "majorminorpatch"
	AllFormats            = "all"
)
const defaultFormat = SemverFormat

var verboseLogging bool
var help bool
var dir string
var build string
var mainBranchRegex string
var majorRegex string
var patchRegex string
var outputFormat string

var majorRegexp *regexp.Regexp
var patchRegexp *regexp.Regexp
var mainBranchRegexp *regexp.Regexp

func SetupFlags() {
	flag.BoolVar(&verboseLogging, "v", false, "Enables debug logging.")
	flag.BoolVar(&help, "h", false, "Shows the help.")
	flag.StringVar(&dir, "repo", "", "The git directory. (defaults to current directory)")
	flag.StringVar(&build, "build", "", "Dot-separated build identifier.")
	flag.StringVar(&mainBranchRegex, "mainBranch", gver.DefaultMainBranch.String(), "The branch which is counted as release branch (regex supported). (defaults to main).")
	flag.StringVar(&majorRegex, "major", gver.DefaultMajorRegex.String(), "The regex by which de major number is counted.")
	flag.StringVar(&patchRegex, "minor", gver.DefaultPatchRegex.String(), "The regex by which de patch number is counted (only if message does not match major).")
	flag.StringVar(&outputFormat, "format", defaultFormat, `The requested output format (defaults to semver). Can be either
		- major (major)
		- minor (minor)
		- patch (patch)
		- majorMinorPath (major.minor.patch)
		- semver (major.minor.patch-prerelease+build)
		- all (a json document containing all other formats)`)

	if mainBranchRegex == "" {
		mainBranchRegex = "main"
	}

	flag.Parse()
}

func main() {
	SetupFlags()

	if verboseLogging {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if help {
		flag.PrintDefaults()
		return
	}

	majorRegexp = regexp.MustCompile(majorRegex)
	patchRegexp = regexp.MustCompile(patchRegex)
	mainBranchRegexp = regexp.MustCompile(mainBranchRegex)

	slog.Debug("Git director", slog.String("dir", dir))
	slog.Debug("Build version", slog.String("build", build))
	slog.Debug("Main branch name", slog.String("mainBranch", mainBranchRegex))

	repo, err := git.PlainOpen(dir)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to open git repository"))
	}

	semver, err := gver.Build(repo, gver.Options{
		MajorBumpRegex:  majorRegexp,
		PatchBumpRegex:  patchRegexp,
		MainBranchRegex: mainBranchRegexp,
		IncludeBranch:   true,
		Build:           build,
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to get latest version tag from repo"))
	}

	switch strings.ToLower(outputFormat) {
	case SemverFormat:
		fmt.Println(semver.String())
	case MajorMinorPatchFormat:
		fmt.Println(semver.MajorMinorPatch.String())
	case AllFormats:
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(struct {
			Major           uint64
			Minor           uint64
			Patch           uint64
			Prerelease      string
			BuildMetadata   string
			MajorMinorPatch string
			Semver          string
		}{
			Major:           semver.MajorMinorPatch.Major,
			Minor:           semver.MajorMinorPatch.Minor,
			Patch:           semver.MajorMinorPatch.Patch,
			Prerelease:      semver.Prerelease,
			BuildMetadata:   semver.Build,
			MajorMinorPatch: semver.MajorMinorPatch.String(),
			Semver:          semver.String(),
		})
		if err != nil {
			panic(err)
		}
	}
}
