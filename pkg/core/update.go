package core

import (
	"github.com/menathan/semverbot/pkg/versions"
)

// UpdateVersion fetches the latest tags from the git repository.
// Returns an error if updating the version went wrong.
func UpdateVersion() error {
	var versionAPI = versions.NewAPI("", "") // prefix and suffix are not used
	return versionAPI.UpdateVersion()
}
