package core

import (
	"github.com/restechnica/semverbot/pkg/versions"
)

type GetVersionOptions struct {
	DefaultVersion  string
	GitTagPrefix    string
	GitTagSuffix    string
	SemVerParseMode string
	SemVerPerPrefix bool
}

// GetVersion gets the current version.
// Returns the current version.
func GetVersion(options *GetVersionOptions) string {
	var versionAPI = versions.NewAPI(options.GitTagPrefix, options.GitTagSuffix)
	return versionAPI.GetVersionOrDefault(options.DefaultVersion, options.SemVerParseMode, options.SemVerPerPrefix, options.GitTagPrefix)
}
