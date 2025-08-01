package main

import (
	"flag"
	"fmt"
	"github.com/go-git/go-git/v6"
	"github.com/pkg/errors"
	gittools "gver/git"
	"gver/semver"
	"log"
	"log/slog"
	"regexp"
)

const defaultMajorRegex = `^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\([\w\s-]*\))?(!:|:.*\n\n((.+\n)+\n)?BREAKING CHANGE:\s.+)`
const defaultMinorRegex = `^(feat)(\\([\\w\\s-]*\\))?:`
const defaultMainBranch = "main"
const defaultFormat = "full"

var verboseLogging bool
var help bool
var dir string
var build string
var mainBranchRegex string
var majorRegex string
var minorRegex string
var outputFormat string

var majorRegexp *regexp.Regexp
var minorRegexp *regexp.Regexp
var mainBranchRegexp *regexp.Regexp

func SetupFlags() {
	flag.BoolVar(&verboseLogging, "v", false, "Enables debug logging.")
	flag.BoolVar(&help, "h", false, "Shows the help.")
	flag.StringVar(&dir, "repo", "", "The git directory. (defaults to current directory)")
	flag.StringVar(&build, "build", "", "Dot-separated build identifier.")
	flag.StringVar(&mainBranchRegex, "mainBranch", defaultMainBranch, "The branch which is counted as release branch (regex supported). (defaults to main).")
	flag.StringVar(&majorRegex, "major", defaultMajorRegex, "The regex by which de major number is counted.")
	flag.StringVar(&minorRegex, "minor", defaultMinorRegex, "The regex by which de minor number is counted.")
	flag.StringVar(&outputFormat, "format", defaultFormat, "The requested output format. Can be either full or majorMinorPath. (defaults to full)")

	if mainBranchRegex == "" {
		mainBranchRegex = "main"
	}

	flag.Parse()
}

func CompileRegexes() {
	var err error
	majorRegexp, err = regexp.Compile(majorRegex)
	if err != nil {
		log.Fatal(err)
	}
	minorRegexp, err = regexp.Compile(minorRegex)
	if err != nil {
		log.Fatal(err)
	}
	mainBranchRegexp, err = regexp.Compile(mainBranchRegex)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to compile mainBranch regex"))
	}
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

	CompileRegexes()

	slog.Debug("Git director", slog.String("dir", dir))
	slog.Debug("Build version", slog.String("build", build))
	slog.Debug("Main branch name", slog.String("mainBranch", mainBranchRegex))

	repo, err := git.PlainOpen(dir)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to open git repository"))
	}

	version := semver.SemVer{Build: build}

	version.Major, version.Minor, version.Patch, err = gittools.MajorMinorPatch(repo, majorRegexp, minorRegexp)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to get majorMinorPatch"))
	}

	branchName, err := gittools.BranchName(repo)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to get branch name"))
	}

	if mainBranchRegexp.MatchString(branchName) {
		slog.Debug("Branch name matches main branch (not adding it to the version).",
			slog.String("mainBranch", mainBranchRegex),
			slog.String("currentBranch", branchName))
	} else {
		version.Prerelease = branchName
	}

	switch outputFormat {
	case "full":
		fmt.Println(version.Full())
	case "majorMinorPatch":
		fmt.Println(version.MajorMinorPatch())
	}
}
