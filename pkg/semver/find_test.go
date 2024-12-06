package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	type Test struct {
		Name      string
		Prefix    string
		Versions  []string
		WantIndex int
	}

	var defaultPrefix = "v"
	var tests = []Test{
		// Using defaultPrefix - Match
		{Name: "MatchVersionUsingDefaultPrefixNoPrefix", Prefix: defaultPrefix, Versions: []string{"1.0.1", "0.2.1", "0.2.0"}, WantIndex: 0},
		{Name: "MatchVersionUsingDefaultPrefixNoPrefixOneSuffix", Prefix: defaultPrefix, Versions: []string{"1.0.1-pre", "0.2.1-pre", "0.2.0-pre"}, WantIndex: 0},
		{Name: "MatchVersionUsingDefaultPrefixNoPrefixMultipleSuffix", Prefix: defaultPrefix, Versions: []string{"1.0.1-pre", "0.2.1-pre2", "0.2.0-pre"}, WantIndex: 0},
		{Name: "MatchVersionUsingDefaultPrefixNoPrefixDifferentOrder", Prefix: defaultPrefix, Versions: []string{"0.2.0", "1.3.1", "1.3.0"}, WantIndex: 1},
		{Name: "MatchVersionUsingDefaultPrefixAllDefaultPrefix", Prefix: defaultPrefix, Versions: []string{"v1.0.1", "v0.2.1", "v0.2.0"}, WantIndex: 0},
		{Name: "MatchVersionUsingDefaultPrefixAllDefaultPrefixOneSuffix", Prefix: defaultPrefix, Versions: []string{"v1.0.1-pre", "v0.2.1-pre", "v0.2.0-pre"}, WantIndex: 0},
		{Name: "MatchVersionUsingDefaultPrefixAllDefaultPrefixMultipleSuffix", Prefix: defaultPrefix, Versions: []string{"v1.3.1-pre1", "v0.2.0", "v0.2.0-pre3"}, WantIndex: 0},
		{Name: "MatchVersionUsingDefaultPrefixAllDefaultPrefixDifferentOrder", Prefix: defaultPrefix, Versions: []string{"v0.2.0", "v1.3.1", "v1.3.0"}, WantIndex: 1},
		{Name: "MatchVersionUsingDefaultPrefixNoPrefixDifferentOrder", Prefix: defaultPrefix, Versions: []string{"0.2.0", "1.3.1", "1.3.0"}, WantIndex: 1},
		{Name: "MatchVersionUsingDefaultPrefixMultiplePrefixNoSuffix", Prefix: defaultPrefix, Versions: []string{"1.3.1", "v0.2.0", "test/1.3.0"}, WantIndex: 0},
		{Name: "MatchVersionUsingDefaultPrefixMultiplePrefixOneSuffix", Prefix: defaultPrefix, Versions: []string{"1.3.1-pre", "v0.2.0-pre", "test/1.3.0-pre"}, WantIndex: 0},
		{Name: "MatchVersionUsingDefaultPrefixMultiplePrefixMultipleSuffix", Prefix: defaultPrefix, Versions: []string{"1.3.1-pre", "v0.2.0-pre2", "test/1.3.0-pre"}, WantIndex: 0},

		// Using defaultPrefix - No Match
		{Name: "NoMatchVersionUsingDefaultPrefixOnePrefixNoSuffix", Prefix: defaultPrefix, Versions: []string{"test/1.0.1", "test/0.2.1", "test/0.2.0", "0.1.0"}, WantIndex: 3},
		{Name: "NoMatchVersionUsingDefaultPrefixOnePrefixOneSuffix", Prefix: defaultPrefix, Versions: []string{"test/1.0.1-pre", "test/0.2.1-pre", "test/0.2.0-pre", "0.1.0"}, WantIndex: 3},
		{Name: "NoMatchVersionUsingDefaultPrefixOnePrefixMultipleSuffix", Prefix: defaultPrefix, Versions: []string{"test/1.0.1-pre", "test/0.2.1-pre2", "test/0.2.0-pre", "0.1.0"}, WantIndex: 3},

		// Using defaultPrefix - Invalid
		{Name: "InvalidVersionUsingDefaultPrefixOneInvalid", Prefix: defaultPrefix, Versions: []string{"invalid1", "0.2.1", "0.2.0"}, WantIndex: 1},
		{Name: "InvalidVersionUsingDefaultPrefixMultipleInvalid", Prefix: defaultPrefix, Versions: []string{"invalid1", "0.2.1", "0.2.0invalid2"}, WantIndex: 1},
		{Name: "InvalidVersionUsingDefaultPrefixAllInvalid", Prefix: defaultPrefix, Versions: []string{"invalid1", "invalid2", "v1.0.0invalid", "0.1.0"}, WantIndex: 3},

		// Using empty "" Prefix - Match
		{Name: "MatchVersionUsingNoPrefixNoPrefix", Prefix: "", Versions: []string{"1.0.1", "0.2.1", "0.2.0"}, WantIndex: 0},
		{Name: "MatchVersionUsingNoPrefixNoPrefixOneSuffix", Prefix: "", Versions: []string{"1.0.1-pre", "0.2.1-pre", "0.2.0-pre"}, WantIndex: 0},
		{Name: "MatchVersionUsingNoPrefixNoPrefixMultipleSuffix", Prefix: "", Versions: []string{"1.0.1-pre", "0.2.1-pre2", "0.2.0-pre"}, WantIndex: 0},
		{Name: "MatchVersionUsingNoPrefixNoPrefixDifferentOrder", Prefix: "", Versions: []string{"0.2.0", "1.3.1", "1.3.0"}, WantIndex: 1},
		{Name: "MatchVersionUsingNoPrefixMultiplePrefixNoSuffix", Prefix: "", Versions: []string{"1.3.1", "test2/0.2.0", "test/1.3.0"}, WantIndex: 0},
		{Name: "MatchVersionUsingNoPrefixMultiplePrefixOneSuffix", Prefix: "", Versions: []string{"1.3.1-pre", "test2/0.2.0-pre", "test/1.3.0-pre"}, WantIndex: 0},
		{Name: "MatchVersionUsingNoPrefixMultiplePrefixMultipleSuffix", Prefix: "", Versions: []string{"1.3.0-pre", "test2/0.2.0-pre2", "test/1.3.1-pre"}, WantIndex: 0},

		// Using empty "" Prefix - No Match
		{Name: "NoMatchVersionUsingNoPrefixAllDefaultPrefix", Prefix: "", Versions: []string{"v1.0.1", "v0.2.1", "v0.2.0", "0.1.0"}, WantIndex: 3},
		{Name: "NoMatchVersionUsingNoPrefixAllDefaultPrefixOneSuffix", Prefix: "", Versions: []string{"v1.0.1-pre", "v0.2.1-pre", "v0.2.0-pre", "0.1.0"}, WantIndex: 3},
		{Name: "NoMatchVersionUsingNoPrefixAllDefaultPrefixMultipleSuffix", Prefix: "", Versions: []string{"v1.0.1-pre", "v0.2.1-pre2", "v0.2.0-pre", "0.1.0"}, WantIndex: 3},
		{Name: "NoMatchVersionUsingNoPrefixDifferentPrefixNoSuffix", Prefix: "", Versions: []string{"test/1.0.1", "test/0.2.1", "test/0.2.0", "0.1.0"}, WantIndex: 3},
		{Name: "NoMatchVersionUsingNoPrefixDifferentPrefixOneSuffix", Prefix: "", Versions: []string{"test/1.0.1-pre", "test/0.2.1-pre", "test/0.2.0-pre", "0.1.0"}, WantIndex: 3},
		{Name: "NoMatchVersionUsingNoPrefixDifferentPrefixMultipleSuffix", Prefix: "", Versions: []string{"test/1.0.1-pre", "test/0.2.1-pre2", "test/0.2.0-pre", "0.1.0"}, WantIndex: 3},
		{Name: "NoMatchVersionUsingNoPrefixDifferentPrefixDifferentOrder", Prefix: "", Versions: []string{"test/0.2.0", "test/1.3.1", "test/1.3.0", "0.1.0"}, WantIndex: 3},

		// Using empty "" Prefix - Invalid
		{Name: "InvalidVersionUsingNoPrefixOneInvalid", Prefix: "", Versions: []string{"0.2.0", "0.2.1invalid", "0.3.0"}, WantIndex: 2},
		{Name: "InvalidVersionUsingNoPrefixMultipleInvalid", Prefix: "", Versions: []string{"invalid1", "0.2.1", "0.2.0invalid2"}, WantIndex: 1},
		{Name: "InvalidVersionUsingNoPrefixAllInvalid", Prefix: "", Versions: []string{"invalid1", "invalid2", "test/1.0.0invalid", "0.1.0"}, WantIndex: 3},

		// Using custom "/test" Prefix - Match
		{Name: "MatchVersionUsingCustomPrefixSamePrefixNoSuffix", Prefix: "test/", Versions: []string{"test/1.0.1", "test/0.2.1", "test/0.2.0"}, WantIndex: 0},
		{Name: "MatchVersionUsingCustomPrefixSamePrefixOneSuffix", Prefix: "test/", Versions: []string{"test/1.0.1-pre", "test/0.2.1-pre", "test/0.2.0-pre"}, WantIndex: 0},
		{Name: "MatchVersionUsingCustomPrefixSamePrefixMultipleSuffix", Prefix: "test/", Versions: []string{"test/1.0.1-pre", "test/0.2.1-pre2", "test/0.2.0-pre"}, WantIndex: 0},
		{Name: "MatchVersionUsingCustomPrefixSamePrefixDifferentOrder", Prefix: "test/", Versions: []string{"test/0.2.0", "test/1.3.1", "test/1.3.0"}, WantIndex: 1},
		{Name: "MatchVersionUsingCustomPrefixMultiplePrefixNoSuffix", Prefix: "test/", Versions: []string{"1.3.1", "test2/0.2.0", "test/1.3.0"}, WantIndex: 2},
		{Name: "MatchVersionUsingCustomPrefixMultiplePrefixOneSuffix", Prefix: "test/", Versions: []string{"1.3.1-pre", "test2/0.2.0-pre", "test/1.3.0-pre"}, WantIndex: 2},
		{Name: "MatchVersionUsingCustomPrefixMultiplePrefixMultipleSuffix", Prefix: "test/", Versions: []string{"1.3.0-pre", "test2/0.2.0-pre2", "test/1.3.1-pre"}, WantIndex: 2},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var versions = test.Versions
			var want = versions[test.WantIndex]
			var got, err = Find(test.Prefix, versions)

			assert.Equal(t, want, got, `want: "%s", got: "%s"`, want, got)
			assert.NoError(t, err)
		})
	}

	type ErrorTest struct {
		Name     string
		Prefix   string
		Versions []string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnErrorUsingCustomPrefixNoPrefix", Prefix: "test/", Versions: []string{"1.0.1", "0.2.1", "0.2.0"}},
		{Name: "ReturnErrorUsingCustomPrefixNoPrefixOneSuffix", Prefix: "test/", Versions: []string{"1.0.1-pre", "0.2.1-pre", "0.2.0-pre"}},
		{Name: "ReturnErrorUsingCustomPrefixNoPrefixMultipleSuffix", Prefix: "test/", Versions: []string{"1.0.1-pre", "0.2.1-pre2", "0.2.0-pre"}},
		{Name: "ReturnErrorUsingCustomPrefixNoPrefixDifferentOrder", Prefix: "test/", Versions: []string{"0.2.0", "1.3.1", "1.3.0"}},
		{Name: "ReturnErrorUsingCustomPrefixAllDefaultPrefix", Prefix: "test/", Versions: []string{"v1.0.1", "v0.2.1", "v0.2.0"}},
		{Name: "ReturnErrorUsingCustomPrefixAllDefaultPrefixOneSuffix", Prefix: "test/", Versions: []string{"v1.0.1-pre", "v0.2.1-pre", "v0.2.0-pre"}},
		{Name: "ReturnErrorUsingCustomPrefixAllDefaultPrefixMultipleSuffix", Prefix: "test/", Versions: []string{"v1.0.1-pre", "v0.2.1-pre2", "v0.2.0-pre"}},
		{Name: "ReturnErrorUsingCustomPrefixMultiplePrefix", Prefix: "test/", Versions: []string{"v1.3.2", "v1.3.1", "c1.3.1"}},
		{Name: "ReturnErrorUsingCustomPrefixMultiplePrefixOneSuffix", Prefix: "test/", Versions: []string{"v1.3.2-pre", "v1.3.1-pre", "c1.3.1-pre"}},
		{Name: "ReturnErrorUsingCustomPrefixMultiplePrefixMultipleSuffix", Prefix: "test/", Versions: []string{"v1.3.2-pre", "v1.3.1-pre2", "c1.3.1-pre"}},
		{Name: "ReturnErrorUsingCustomPrefixOneInvalid", Prefix: "test/", Versions: []string{"0.2.0", "0.2.1invalid", "0.3.0"}},
		{Name: "ReturnErrorUsingCustomPrefixMultipleInvalid", Prefix: "test/", Versions: []string{"invalid1", "0.2.1", "0.2.0invalid2"}},
		{Name: "ReturnErrorUsingCustomPrefixAllInvalid", Prefix: "test/", Versions: []string{"invalid1", "invalid2", "test/1.0.0invalid"}},
		{Name: "ReturnErrorOnInvalidVersions", Versions: []string{"invalid", "semver", "versions"}},
		{Name: "ReturnErrorOnNoVersions", Versions: []string{}},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var versions = test.Versions
			var _, got = Find(test.Prefix, versions)
			assert.Error(t, got)
		})
	}
}
