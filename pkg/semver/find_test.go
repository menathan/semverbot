package semver

import (
	"testing"

	"github.com/restechnica/semverbot/pkg/semver/semverparse"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	type Test struct {
		Name            string
		Versions        []string
		WantIndex       int
		SemVerParseMode string
		SemVerPerPrefix bool
		Prefix          string
	}

	var tests = []Test{
		// Using empty "" Prefix - Match
		{Name: "MatchVersionUsingNoPrefixNoPrefix", Versions: []string{"1.0.1", "0.2.1", "0.2.0"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "MatchVersionUsingNoPrefixNoPrefixOneSuffix", Versions: []string{"1.0.1-pre", "0.2.1-pre", "0.2.0-pre"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "MatchVersionUsingNoPrefixNoPrefixMultipleSuffix", Versions: []string{"1.0.1-pre", "0.2.1-pre2", "0.2.0-pre"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "MatchVersionUsingNoPrefixNoPrefixDifferentOrder", Versions: []string{"0.2.0", "1.3.1", "1.3.0"}, WantIndex: 1, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "MatchVersionUsingNoPrefixWithVPrefixDifferentOrder", Versions: []string{"v0.2.0", "v1.3.1", "v1.3.0"}, WantIndex: 1, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "MatchVersionUsingNoPrefixMultiplePrefix1", Versions: []string{"1.3.1", "bar/0.2.0", "foo/1.5.0"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "MatchVersionUsingNoPrefixMultiplePrefix2", Versions: []string{"1.3.1-pre", "bar/0.2.0-pre", "foo/1.3.0-pre"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},

		// Using empty "" Prefix - No Match
		{Name: "NoMatchVersionUsingNoPrefixDifferentPrefixNoSuffix", Versions: []string{"foo/1.0.1", "foo/0.2.1", "foo/0.2.0", "0.1.0"}, WantIndex: 3, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "NoMatchVersionUsingNoPrefixDifferentPrefixOneSuffix", Versions: []string{"foo/1.0.1-pre", "foo/0.2.1-pre", "foo/0.2.0-pre", "0.1.0"}, WantIndex: 3, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "NoMatchVersionUsingNoPrefixDifferentPrefixMultipleSuffix", Versions: []string{"foo/1.0.1-pre", "foo/0.2.1-pre2", "foo/0.2.0-pre", "0.1.0"}, WantIndex: 3, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "NoMatchVersionUsingNoPrefixDifferentPrefixDifferentOrder", Versions: []string{"foo/0.2.0", "foo/1.3.1", "foo/1.3.0", "0.1.0"}, WantIndex: 3, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},

		// Using empty "" Prefix - Invalid
		{Name: "InvalidVersionUsingNoPrefixOneInvalid", Versions: []string{"0.2.0", "0.2.1invalid", "0.3.0"}, WantIndex: 2, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "InvalidVersionUsingNoPrefixMultipleInvalid", Versions: []string{"invalid1", "0.2.1", "0.2.0invalid2"}, WantIndex: 1, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "InvalidVersionUsingNoPrefixAllInvalid", Versions: []string{"invalid1", "invalid2", "foo/1.0.0invalid", "0.1.0"}, WantIndex: 3, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},

		// Using custom "foo/" Prefix - Match
		{Name: "MatchVersionUsingCustomPrefixSamePrefixNoSuffix", Versions: []string{"foo/1.0.1", "foo/0.2.1", "foo/0.2.0"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "MatchVersionUsingCustomPrefixSamePrefixOneSuffix", Versions: []string{"foo/1.0.1-pre", "foo/0.2.1-pre", "foo/0.2.0-pre"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "MatchVersionUsingCustomPrefixSamePrefixMultipleSuffix", Versions: []string{"foo/1.0.1-pre", "foo/0.2.1-pre2", "foo/0.2.0-pre"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "MatchVersionUsingCustomPrefixSamePrefixDifferentOrder", Versions: []string{"foo/0.2.0", "foo/1.3.1", "foo/1.3.0"}, WantIndex: 1, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "MatchVersionUsingCustomPrefixMultiplePrefixNoSuffix", Versions: []string{"1.3.1", "bar/0.2.0", "foo/1.3.0"}, WantIndex: 2, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "MatchVersionUsingCustomPrefixMultiplePrefixOneSuffix", Versions: []string{"1.3.1-pre", "bar/0.2.0-pre", "foo/1.3.0-pre"}, WantIndex: 2, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "MatchVersionUsingCustomPrefixMultiplePrefixMultipleSuffix", Versions: []string{"1.3.0-pre", "bar/0.2.0-pre2", "foo/1.3.1-pre"}, WantIndex: 2, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var versions = test.Versions
			var want, _ = ParseWithOptions(versions[test.WantIndex], test.SemVerParseMode, test.SemVerPerPrefix, test.Prefix)
			var got, err = Find(versions, test.SemVerParseMode, test.SemVerPerPrefix, test.Prefix)

			assert.Equal(t, want.String(), got, `want: "%s", got: "%s"`, want, got)
			assert.NoError(t, err)
		})
	}

	var testsStrictMatch = []Test{
		// Using empty "" Prefix - Match
		{Name: "StrictMatchVersionUsingNoPrefixNoPrefix", Versions: []string{"1.0.1", "0.2.1", "0.2.0"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "StrictMatchVersionUsingNoPrefixNoPrefixOneSuffix", Versions: []string{"1.0.1-pre", "0.2.1-pre", "0.2.0-pre"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "StrictMatchVersionUsingNoPrefixNoPrefixMultipleSuffix", Versions: []string{"1.0.1-pre", "0.2.1-pre2", "0.2.0-pre"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "StrictMatchVersionUsingNoPrefixNoPrefixDifferentOrder", Versions: []string{"0.2.0", "1.3.1", "1.3.0"}, WantIndex: 1, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "StrictMatchVersionUsingNoPrefixMultiplePrefixNoSuffix", Versions: []string{"1.3.1", "bar/0.2.0", "foo/1.3.0"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "StrictMatchVersionUsingNoPrefixMultiplePrefixOneSuffix", Versions: []string{"1.3.1-pre", "bar/0.2.0-pre", "foo/1.3.0-pre"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "StrictMatchVersionUsingNoPrefixMultiplePrefixMultipleSuffix", Versions: []string{"1.3.0-pre", "bar/0.2.0-pre2", "foo/1.3.1-pre"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},

		// Using empty "" Prefix - No Match
		{Name: "NoStrictMatchVersionUsingNoPrefixDifferentPrefixNoSuffix", Versions: []string{"foo/1.0.1", "foo/0.2.1", "foo/0.2.0", "0.1.0"}, WantIndex: 3, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "NoStrictMatchVersionUsingNoPrefixDifferentPrefixOneSuffix", Versions: []string{"foo/1.0.1-pre", "foo/0.2.1-pre", "foo/0.2.0-pre", "0.1.0"}, WantIndex: 3, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "NoStrictMatchVersionUsingNoPrefixDifferentPrefixMultipleSuffix", Versions: []string{"foo/1.0.1-pre", "foo/0.2.1-pre2", "foo/0.2.0-pre", "0.1.0"}, WantIndex: 3, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "NoStrictMatchVersionUsingNoPrefixDifferentPrefixDifferentOrder", Versions: []string{"foo/0.2.0", "foo/1.3.1", "foo/1.3.0", "0.1.0"}, WantIndex: 3, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},

		// Using custom "foo/" Prefix - Match
		{Name: "StrictMatchVersionUsingCustomPrefixSamePrefixNoSuffix", Versions: []string{"foo/1.0.1", "foo/0.2.1", "foo/0.2.0"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "StrictMatchVersionUsingCustomPrefixSamePrefixOneSuffix", Versions: []string{"foo/1.0.1-pre", "foo/0.2.1-pre", "foo/0.2.0-pre"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "StrictMatchVersionUsingCustomPrefixSamePrefixMultipleSuffix", Versions: []string{"foo/1.0.1-pre", "foo/0.2.1-pre2", "foo/0.2.0-pre"}, WantIndex: 0, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "StrictMatchVersionUsingCustomPrefixSamePrefixDifferentOrder", Versions: []string{"foo/0.2.0", "foo/1.3.1", "foo/1.3.0"}, WantIndex: 1, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "StrictMatchVersionUsingCustomPrefixMultiplePrefixNoSuffix", Versions: []string{"1.3.1", "bar/0.2.0", "foo/1.3.0"}, WantIndex: 2, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "StrictMatchVersionUsingCustomPrefixMultiplePrefixOneSuffix", Versions: []string{"1.3.1-pre", "bar/0.2.0-pre", "foo/1.3.0-pre"}, WantIndex: 2, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "StrictMatchVersionUsingCustomPrefixMultiplePrefixMultipleSuffix", Versions: []string{"1.3.0-pre", "bar/0.2.0-pre2", "foo/1.3.1-pre"}, WantIndex: 2, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
	}

	for _, test := range testsStrictMatch {
		t.Run(test.Name, func(t *testing.T) {
			var versions = test.Versions
			var want, _ = ParseWithOptions(versions[test.WantIndex], test.SemVerParseMode, test.SemVerPerPrefix, test.Prefix)
			var got, err = Find(versions, test.SemVerParseMode, test.SemVerPerPrefix, test.Prefix)

			assert.Equal(t, want.String(), got, `want: "%s", got: "%s"`, want, got)
			assert.NoError(t, err)
		})
	}

	type ErrorTest struct {
		Name            string
		Versions        []string
		SemVerParseMode string
		SemVerPerPrefix bool
		Prefix          string
	}

	var errorTestsStrictMatch = []ErrorTest{
		{Name: "ReturnErrorUsingCustomPrefixNoPrefix", Versions: []string{"1.0.1", "0.2.1", "0.2.0"}, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "ReturnErrorUsingCustomPrefixNoPrefixOneSuffix", Versions: []string{"1.0.1-pre", "0.2.1-pre", "0.2.0-pre"}, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "ReturnErrorUsingCustomPrefixNoPrefixMultipleSuffix", Versions: []string{"1.0.1-pre", "0.2.1-pre2", "0.2.0-pre"}, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "ReturnErrorUsingCustomPrefixNoPrefixDifferentOrder", Versions: []string{"0.2.0", "1.3.1", "1.3.0"}, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "ReturnErrorUsingCustomPrefixMultiplePrefix", Versions: []string{"v1.3.2", "v1.3.1", "c1.3.1"}, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "ReturnErrorUsingCustomPrefixMultiplePrefixOneSuffix", Versions: []string{"v1.3.2-pre", "v1.3.1-pre", "c1.3.1-pre"}, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "ReturnErrorUsingCustomPrefixMultiplePrefixMultipleSuffix", Versions: []string{"v1.3.2-pre", "v1.3.1-pre2", "c1.3.1-pre"}, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "ReturnErrorUsingCustomPrefixOneInvalid", Versions: []string{"0.2.0", "0.2.1invalid", "0.3.0"}, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "ReturnErrorUsingCustomPrefixMultipleInvalid", Versions: []string{"invalid1", "0.2.1", "0.2.0invalid2"}, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "ReturnErrorUsingCustomPrefixAllInvalid", Versions: []string{"invalid1", "invalid2", "foo/1.0.0invalid"}, SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "ReturnErrorOnInvalidVersions", Versions: []string{"invalid", "semver", "versions"}},
		{Name: "ReturnErrorOnNoVersions", Versions: []string{}},
	}

	for _, test := range errorTestsStrictMatch {
		t.Run(test.Name, func(t *testing.T) {
			var versions = test.Versions
			var _, got = Find(versions, test.SemVerParseMode, test.SemVerPerPrefix, test.Prefix)
			assert.Error(t, got)
		})
	}
}
