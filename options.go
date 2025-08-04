package gver

import "regexp"

type Options struct {
	MajorBumpRegex  *regexp.Regexp
	PatchBumpRegex  *regexp.Regexp
	MainBranchRegex *regexp.Regexp
	IncludeBranch   bool
	Build           string
}
