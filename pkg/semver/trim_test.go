package semver

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	type Test struct {
		Name        string
		Prefix      string
		Major       string
		Minor       string
		Patch       string
		Suffix      string
		PreAndBuild string
	}

	var tests = []Test{
		{Name: "Default", Major: "0", Minor: "0", Patch: "0"},
		{Name: "Patch", Major: "0", Minor: "0", Patch: "1"},
		{Name: "Minor", Major: "0", Minor: "2", Patch: "0"},
		{Name: "Major", Major: "3", Minor: "0", Patch: "0"},
		{Name: "DiscardPrefix", Prefix: "v", Major: "1", Minor: "0", Patch: "0"},
		{Name: "DiscardPreAndBuildWithFull", Major: "1", Minor: "0", Patch: "0", PreAndBuild: "-any.valid+version"},
		{Name: "DiscardPreAndBuildWithFull2", Major: "2", Minor: "0", Patch: "0", PreAndBuild: "-pre+001"},
		{Name: "DiscardPreAndBuildWithPre", Major: "1", Minor: "0", Patch: "0", PreAndBuild: "-pre"},
		{Name: "DiscardPreAndBuildWithBuild", Major: "1", Minor: "0", Patch: "0", PreAndBuild: "+001"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var version = fmt.Sprintf(`%s%s.%s.%s%s%s`, test.Prefix, test.Major, test.Minor, test.Patch,
				test.Suffix, test.PreAndBuild)

			var want = strings.ReplaceAll(version, test.Prefix, "")
			want = strings.ReplaceAll(want, test.Suffix, "")
			want = strings.ReplaceAll(want, test.PreAndBuild, "")

			var got, err = Trim(version)

			assert.Equal(t, want, got, `want: "%s", got: "%s"`, want, got)

			if test.Prefix != "" {
				assert.False(t, strings.HasPrefix(got, test.Prefix))
			}

			if test.Suffix != "" {
				assert.False(t, strings.HasSuffix(got, test.Suffix))
			}

			if test.PreAndBuild != "" {
				assert.False(t, strings.HasSuffix(got, test.PreAndBuild))
			}

			assert.NoError(t, err)
		})
	}

	type ErrorTest struct {
		Name    string
		Version string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnErrorOnInvalidVersion", Version: "invalid"},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var _, got = Trim(test.Version)
			assert.Error(t, got)
		})
	}
}
