package semver

import "fmt"

type MajorMinorPatch struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func (mmp MajorMinorPatch) String() string {
	return fmt.Sprintf("%d.%d.%d", mmp.Major, mmp.Minor, mmp.Patch)
}
