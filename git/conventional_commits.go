package git

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

func MajorMinorPatch(repo *git.Repository, majorRegexp *regexp.Regexp, minorRegexp *regexp.Regexp) (major int, minor int, patch int, err error) {
	version, hash, err := LatestVersionTag(repo)
	if err != nil {
		return 0, 0, 0, errors.Wrap(err, "Failed to get latest version tag")
	}
	if version != nil {
		major = version.Major
		minor = version.Minor
		patch = version.Patch
	}

	commits, err := repo.Log(&git.LogOptions{All: true})
	if err != nil {
		return 0, 0, 0, errors.Wrap(err, "Failed to get git log")
	}

	var firstCommit *object.Commit
	var commitCount int

	err = commits.ForEach(func(commit *object.Commit) error {
		firstCommit = commit
		commitCount++

		if commit.Hash == hash {
			return storer.ErrStop
		}

		var shortMessage string
		if len(commit.Message) > 100 {
			shortMessage = commit.Message[:100]
		} else {
			shortMessage = commit.Message
		}
		if majorRegexp.MatchString(commit.Message) {
			major++
			slog.Debug("MajorMinorPatch: Major version increase.",
				slog.String("commitMessage", shortMessage),
				slog.Int("major", major),
				slog.Int("minor", minor),
				slog.Int("patch", patch))
			return nil
		}
		if major > 0 {
			// one the latest major version bump is found, we do not need to look for any minor or patches anymore
			return nil
		}

		if minorRegexp.MatchString(commit.Message) {
			minor++
			slog.Debug("MajorMinorPatch: Minor version increase.",
				slog.String("commitMessage", shortMessage),
				slog.Int("major", major),
				slog.Int("minor", minor),
				slog.Int("patch", patch))
			return nil
		}
		if minor > 0 {
			// one the latest minor version bump is found, we do not need to look for any patches anymore
			return nil
		}

		patch++
		slog.Debug("MajorMinorPatch: Patch version increase.",
			slog.String("commitMessage", shortMessage),
			slog.Int("major", major),
			slog.Int("minor", minor),
			slog.Int("patch", patch))
		return nil
	})
	if err != nil {
		return 0, 0, 0, errors.Wrap(err, "failed to determine majorMinorPatch from log")
	}

	slog.Info("Parsed all commits.", slog.Int("count", commitCount))

	if firstCommit == nil {
		slog.Warn("No git history found (no commits)")
	} else if firstCommit.NumParents() != 0 {
		slog.Warn("Did not find the complete git history. Version might be incorrect.")
	}

	return major, minor, patch, nil
}

func LatestVersionTag(repo *git.Repository) (version *semver.SemVer, hash plumbing.Hash, err error) {
	tags, err := repo.Tags()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to get git tags"))
	}

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
		slog.Info("Found version tag", slog.Any("version", *version))
		return storer.ErrStop
	})
	if err != nil {
		err = errors.Wrap(err, "failed to get latest tag")
		return nil, plumbing.Hash{}, err
	}

	return version, hash, nil
}
