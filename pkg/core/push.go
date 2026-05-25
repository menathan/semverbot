package core

import "github.com/menathan/semverbot/pkg/versions"

type PushVersionOptions struct {
	DefaultVersion  string
	Prefix          string
	Suffix          string
	SemVerPerPrefix bool
	SemVerParseMode string
}

// PushVersion pushes the current version.
// Returns an error if the push went wrong.
func PushVersion(options *PushVersionOptions) (err error) {
	var versionAPI = versions.NewAPI(options.Prefix, options.Suffix)
	var version = versionAPI.GetVersionOrDefault(options.DefaultVersion, options.SemVerParseMode, options.SemVerPerPrefix, options.Prefix)
	return versionAPI.PushVersion(version)
}
