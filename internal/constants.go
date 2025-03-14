package internal

import (
	"github.com/restechnica/semverbot/pkg/modes"
	"github.com/restechnica/semverbot/pkg/semver/semverparse"
)

const (
	// DefaultConfigFilePath the default relative filepath to the config file.
	DefaultConfigFilePath = ".semverbot.toml"

	// DefaultGitBranchDelimiters the default delimiters used by the git-branch mode.
	DefaultGitBranchDelimiters = "/"

	// DefaultGitCommitDelimiters the default delimiters used by the git-commit mode.
	DefaultGitCommitDelimiters = "[]/"

	// DefaultGitTagsPrefix the default prefix prepended to git tags.
	DefaultGitTagsPrefix = "v"

	// DefaultGitTagsSuffix the default suffix appended to git tags.
	DefaultGitTagsSuffix = ""

	// DefaultMode the default mode for incrementing versions.
	DefaultMode = modes.Auto

	// DefaultVersion the default version when no other version can be found.
	DefaultVersion = "0.0.0"

	// DefaultSemVerParseMode the mode used to parse semver versions.
	DefaultSemVerParseMode = semverparse.Tolerant

	// DefaultSemVerPerPrefix whether to strictly match the prefix when getting and predicting versions.
	DefaultSemVerPerPrefix = false
)
