package semver

import (
	"fmt"

	blangsemver "github.com/blang/semver/v4"
	"github.com/rs/zerolog/log"
)

// Find finds the biggest valid semver version in a slice of strings.
// The initial order of the versions does not matter.
// Returns the biggest valid semver version if found, otherwise an error stating no valid semver version has been found.
func Find(versions []string, semVerParseMode string, semVerPerPrefix bool, prefix string) (found string, err error) {
	var parsedVersions blangsemver.Versions
	var parsedVersion blangsemver.Version

	for _, version := range versions {
		if parsedVersion, err = ParseWithOptions(version, semVerParseMode, semVerPerPrefix, prefix); err != nil {
			log.Warn().Msg("could not parse tag '" + version + "'")
		}
		if !(parsedVersion.Major == 0 && parsedVersion.Minor == 0 && parsedVersion.Patch == 0) {
			log.Debug().Msg("parsed tag '" + version + "' as version " + parsedVersion.String())
			parsedVersions = append(parsedVersions, parsedVersion)
		}
	}

	if len(parsedVersions) == 0 {
		return found, fmt.Errorf("could not find a valid semver version")
	}

	blangsemver.Sort(parsedVersions)

	var targetVersion = parsedVersions[len(parsedVersions)-1]

	return targetVersion.String(), nil
}
