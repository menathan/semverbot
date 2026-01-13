package core

import (
	"github.com/restechnica/semverbot/pkg/versions"
)

type GetVersionOptions struct {
	DefaultVersion  string
	Prefix          string
	Suffix          string
	SemVerParseMode string
	SemVerPerPrefix bool
}

// GetVersion gets the current version.
// Returns the current version.
func GetVersion(options *GetVersionOptions) string {
	var versionAPI = versions.NewAPI(options.Prefix, options.Suffix)
	return versionAPI.GetVersionOrDefault(options.DefaultVersion, options.SemVerParseMode, options.SemVerPerPrefix, options.Prefix)
}
