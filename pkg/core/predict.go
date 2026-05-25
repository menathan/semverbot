package core

import (
	"github.com/menathan/semverbot/pkg/modes"
	"github.com/menathan/semverbot/pkg/semver"
	"github.com/menathan/semverbot/pkg/versions"
)

type PredictVersionOptions struct {
	DefaultVersion      string
	GitBranchDelimiters string
	GitCommitDelimiters string
	Prefix              string
	Suffix              string
	Mode                string
	SemverMap           semver.Map
	SemVerPerPrefix     bool
	SemVerParseMode     string
}

// PredictVersion predicts a version based on a modes.Mode and a modes.Map.
// The modes.Map values will be matched against git information to detect which semver level to increment.
// Returns the next version or an error if the prediction failed.
func PredictVersion(options *PredictVersionOptions) (prediction string, err error) {
	var gitBranchMode = modes.NewGitBranchMode(options.GitBranchDelimiters, options.SemverMap)
	var gitCommitMode = modes.NewGitCommitMode(options.GitCommitDelimiters, options.SemverMap)

	var versionAPI = versions.NewAPI(options.Prefix, options.Suffix)
	var version = versionAPI.GetVersionOrDefault(options.DefaultVersion, options.SemVerParseMode, options.SemVerPerPrefix, options.Prefix)

	var modeAPI = modes.NewAPI(gitBranchMode, gitCommitMode)
	var mode = modeAPI.SelectMode(options.Mode)

	return versionAPI.PredictVersion(version, mode)
}
