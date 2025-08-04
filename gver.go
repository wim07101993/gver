package gver

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/storer"
	"github.com/pkg/errors"
	"gver/semver"
	"log"
	"log/slog"
	"regexp"
)

var DefaultMajorRegex = regexp.MustCompile(`^(?P<type>[\w\s-]+)!:\s(?P<message>(?:.*\n)*.*)$`)
var DefaultPatchRegex = regexp.MustCompile(`^fix:\s(?P<message>(?:.*\n)*.*)$`)
var DefaultMainBranch = regexp.MustCompile("main")

func Build(repo *git.Repository, options Options) (*semver.SemVer, error) {
	tags, err := repo.Tags()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to get git tags"))
	}

	version, hash, err := LatestVersionTag(tags)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to get latest version tag from repo"))
	}
	if version == nil {
		version = &semver.SemVer{}
	}

	commits, err := repo.Log(&git.LogOptions{All: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get git log")
	}

	mmp, err := BuildMajorMinorPatchFromCommits(commits, hash, options.MajorBumpRegex, options.PatchBumpRegex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to determine major-minor-patch from commit history")
	}

	version.AddMajorMinorPatch(mmp)

	if options.IncludeBranch {
		head, err := repo.Head()
		if err != nil {
			return nil, errors.Wrap(err, "failed getting the head of the repo")
		}
		branch := head.Name().Short()
		slog.Debug("BranchName: Found branch name", slog.String("branch", branch))
		if options.MainBranchRegex.MatchString(branch) {
			slog.Debug("Branch name matches main branch (not adding it to the version).",
				slog.String("mainBranch", options.MainBranchRegex.String()),
				slog.String("currentBranch", branch))
		} else {
			version.Prerelease = semver.SanitizePrerelease(branch)
		}
	}

	if options.Build != "" {
		version.Build = semver.SanitizeBuild(options.Build)
	}
	return version, nil
}

func BuildMajorMinorPatchFromCommits(
	commits object.CommitIter,
	start plumbing.Hash,
	majorRegexp *regexp.Regexp,
	minorRegexp *regexp.Regexp) (semver.MajorMinorPatch, error) {

	var firstCommit *object.Commit
	var commitCount int

	version := semver.MajorMinorPatch{}
	err := commits.ForEach(func(commit *object.Commit) error {
		firstCommit = commit
		commitCount++

		if commit.Hash == start {
			slog.Debug("Found start commit, stop increasing version",
				slog.String("hash", commit.Hash.String()))
			return storer.ErrStop
		}

		var shortMessage string
		if len(commit.Message) > 100 {
			shortMessage = commit.Message[:100]
		} else {
			shortMessage = commit.Message
		}

		if majorRegexp.MatchString(commit.Message) {
			version.Major++
			slog.Debug("MajorMinorPatch: Major version increase.",
				slog.String("commitMessage", shortMessage),
				slog.Any("version", version))
			return nil
		}
		if version.Major > 0 {
			// one the latest major version bump is found, we do not need to look for any minor or patches anymore
			return nil
		}

		if minorRegexp.MatchString(commit.Message) {
			version.Minor++
			slog.Debug("MajorMinorPatch: Minor version increase.",
				slog.String("commitMessage", shortMessage),
				slog.Any("version", version))
			return nil
		}
		if version.Minor > 0 {
			// one the latest minor version bump is found, we do not need to look for any patches anymore
			return nil
		}

		version.Patch++
		slog.Debug("MajorMinorPatch: Patch version increase.",
			slog.String("commitMessage", shortMessage),
			slog.Any("version", version))
		return nil
	})
	if err != nil {
		return semver.MajorMinorPatch{}, errors.Wrap(err, "failed to determine majorMinorPatch from log")
	}

	slog.Debug("Parsed all commits.", slog.Int("count", commitCount))

	if firstCommit == nil {
		slog.Warn("No git history found (no commits)")
	} else if firstCommit.NumParents() != 0 {
		slog.Warn("Did not find the complete git history. Version might be incorrect.")
	}

	return version, nil
}

func LatestVersionTag(tags storer.ReferenceIter) (version *semver.SemVer, hash plumbing.Hash, err error) {
	err = tags.ForEach(func(reference *plumbing.Reference) error {
		name := reference.Name().Short()
		if name[0] == 'v' {
			name = name[1:]
		}

		v, err := semver.Parse(name)
		if err != nil {
			slog.Debug("Found tag but could not parse version", slog.String("tag", reference.Name().Short()))
			return nil
		}

		version = v
		hash = reference.Hash()
		slog.Debug("Found version tag", slog.Any("version", *version))
		return storer.ErrStop
	})
	if err != nil {
		err = errors.Wrap(err, "failed to get latest tag")
		return nil, plumbing.Hash{}, err
	}

	return version, hash, nil
}
