package semver

import (
	"fmt"
	"strings"
	"testing"

	"github.com/restechnica/semverbot/pkg/semver/semverparse"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
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
		{Name: "DiscardEmptyPrefix", Prefix: "", Major: "1", Minor: "0", Patch: "0"},
		{Name: "DiscardCommonPrefix", Prefix: "v", Major: "1", Minor: "0", Patch: "0"},
		{Name: "KeepSuffix", Major: "1", Minor: "0", Patch: "0", PreAndBuild: "-any.valid+version"},
		{Name: "KeepSuffixAsPre", Major: "1", Minor: "0", Patch: "0", PreAndBuild: "-pre"},
		{Name: "KeepSuffixAsBuild", Major: "1", Minor: "0", Patch: "0", PreAndBuild: "+001"},
		{Name: "KeepPrebuild", Major: "2", Minor: "0", Patch: "0", PreAndBuild: "-pre+001"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var version = fmt.Sprintf(`%s%s.%s.%s%s%s`, test.Prefix, test.Major, test.Minor, test.Patch,
				test.Suffix, test.PreAndBuild)

			var got, err = ParseWithOptions(version, semverparse.Tolerant, false, test.Prefix)

			assert.Equal(t, test.Major, fmt.Sprint(got.Major), `want: "%s", got: "%d"`, test.Major, got.Major)
			assert.Equal(t, test.Minor, fmt.Sprint(got.Minor), `want: "%s", got: "%s"`, test.Minor, got.Minor)
			assert.Equal(t, test.Patch, fmt.Sprint(got.Patch), `want: "%s", got: "%s"`, test.Patch, got.Patch)

			if test.Prefix != "" {
				assert.False(t, strings.HasPrefix(got.String(), test.Prefix))
			}

			if test.Suffix != "" {
				assert.True(t, strings.HasSuffix(got.String(), test.Suffix))
			}

			if test.PreAndBuild != "" {
				assert.True(t, strings.HasSuffix(got.String(), test.PreAndBuild))
			}

			assert.NoError(t, err)
		})
	}
}
