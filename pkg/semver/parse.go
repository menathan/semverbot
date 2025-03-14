package semver

import (
	"errors"
	"strings"

	blangsemver "github.com/blang/semver/v4"
	"github.com/restechnica/semverbot/pkg/semver/semverparse"
)

// ParseStrict parses a version string into a semver version struct.
// It does not tolerate versions specifications that do not strictly adhere to semver specs.
// See the library documentation for more information.
// Returns the parsed blang/semver/v4 Version.
func ParseStrict(version string) (blangsemver.Version, error) {
	return blangsemver.Parse(version)
}

// ParseTolerant parses a version string into a semver version struct.
// It tolerates certain version specifications that do not strictly adhere to semver specs.
// See the library documentation for more information.
// Returns the parsed blang/semver/v4 Version.
func ParseTolerant(version string) (blangsemver.Version, error) {
	return blangsemver.ParseTolerant(version)
}

// ParseWithOptions parses a version string into a semver version struct.
// It tolerates certain version specifications that do not strictly adhere to semver specs.
// See the library documentation for more information.
// Returns the parsed blang/semver/v4 Version.
func ParseWithOptions(version string, parseMode string, semVerPerPrefix bool, prefix string) (blangsemver.Version, error) {
	switch parseMode {
	case semverparse.Strict:
		return blangsemver.Parse(version)
	case semverparse.Tolerant:
		if semVerPerPrefix {
			if !strings.HasPrefix(version, prefix) {
				// This should under normal circumstances never happen, as the prefix is added to the git tag cli command
				// For now, keeping it so we can write coherent tests
				return blangsemver.Version{}, errors.New("version '" + version + "' does not have the prefix '" + prefix + "'")
			}
			var versionWithoutPrefix = strings.Replace(version, prefix, "", 1)
			return blangsemver.ParseTolerant(versionWithoutPrefix)
		} else {
			return blangsemver.ParseTolerant(version)
		}
	default:
		return blangsemver.Version{}, errors.New("parseMode '" + parseMode + "' not recognized")
	}
}
