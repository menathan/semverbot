package versions

import (
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/restechnica/semverbot/internal/fakes"
	"github.com/restechnica/semverbot/internal/mocks"
	"github.com/restechnica/semverbot/pkg/cli"
	"github.com/restechnica/semverbot/pkg/git"
	"github.com/restechnica/semverbot/pkg/modes"
	"github.com/restechnica/semverbot/pkg/semver/semverparse"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func TestAPI_GetVersion(t *testing.T) {
	type Test struct {
		Name            string
		GitTag          string
		ExpectedVersion string
		SemVerParseMode string
		SemVerPerPrefix bool
		Prefix          string
	}

	var tests = []Test{
		{Name: "GetVersionStrict", GitTag: "0.1.0", ExpectedVersion: "0.1.0", SemVerParseMode: semverparse.Strict, SemVerPerPrefix: false, Prefix: ""},

		{Name: "GetVersionTolerant", GitTag: "1.0.0", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "GetVersionTolerant", GitTag: "1.0", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "GetVersionTolerant", GitTag: "1", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "GetVersionTolerantWithPrefixInTagAndConfig", GitTag: "v1.0.0", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: "v"},
		{Name: "GetVersionTolerantWithPrefixInTagButDifferentInConfig", GitTag: "v1.0.0", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "GetVersionTolerantWithPrefixInConfigButNotInTag", GitTag: "1.0.0", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: "v"},
		{Name: "GetVersionTolerantWithPrefixInConfigButNotInTag2", GitTag: "1.0.0", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: "abc"},

		{Name: "GetVersionTolerantWithSemVerPerPrefix", GitTag: "1.0.0", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: ""},
		{Name: "GetVersionTolerantWithSemVerPerPrefixWithPrefixInTagAndConfig1", GitTag: "v1.0.0", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "v"},
		{Name: "GetVersionTolerantWithSemVerPerPrefixWithPrefixInTagAndConfig2", GitTag: "abc1.0.0", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "abc"},
		{Name: "GetVersionTolerantWithSemVerPerPrefixWithPrefixInTagAndConfig3", GitTag: "foo/1.0.0", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "GetVersionTolerantWithSemVerPerPrefixWithPrefixInTagAndConfig4", GitTag: "bar@1.0.0", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "bar@"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return(test.GitTag, nil)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: "", GitAPI: gitAPI}

			var got, err = versionAPI.GetVersion(test.SemVerParseMode, test.SemVerPerPrefix, test.Prefix)

			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedVersion, got, `want: "%s, got: "%s"`, test.ExpectedVersion, got)
		})
	}

	type TestWithMultipleTags struct {
		Name            string
		Tags            string
		ExpectedVersion string
		SemVerParseMode string
		SemVerPerPrefix bool
		Prefix          string
	}

	var tests_with_multiple_tags = []TestWithMultipleTags{
		{Name: "GetVersionStrict", Tags: "1.0.0\n1.1.0\n2.0.0", ExpectedVersion: "2.0.0", SemVerParseMode: semverparse.Strict, SemVerPerPrefix: false, Prefix: ""},
		{Name: "GetVersionStrictUnordered", Tags: "1.0.0\n2.0.0\n1.1.0", ExpectedVersion: "2.0.0", SemVerParseMode: semverparse.Strict, SemVerPerPrefix: false, Prefix: ""},
		{Name: "GetVersionStrictIgnored", Tags: "1.0.0\n2.0\n1.1.0", ExpectedVersion: "1.1.0", SemVerParseMode: semverparse.Strict, SemVerPerPrefix: false, Prefix: ""},

		{Name: "GetVersionTolerant", Tags: "1.0.0\n1.1.0\n2.0.0", ExpectedVersion: "2.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "GetVersionTolerantCommon1", Tags: "1.0.0\n2.0\n1.1.0", ExpectedVersion: "2.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "GetVersionTolerantCommon2", Tags: "1.0.0\nv2.0\n2.1", ExpectedVersion: "2.1.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "GetVersionTolerantCommon3", Tags: "1.0.0\n1.1.0\nv2", ExpectedVersion: "2.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "GetVersionTolerantCommon4", Tags: "1.0.0\nv1.1.0\n2", ExpectedVersion: "2.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: ""},
		{Name: "GetVersionTolerantWithPrefixInConfigNotInTag1", Tags: "1.0.0\n1.1.0\n2.0.0", ExpectedVersion: "2.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: "v"},
		{Name: "GetVersionTolerantWithPrefixInConfigNotInTag2", Tags: "1.0.0\n1.1.0\n2.0.0", ExpectedVersion: "2.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: "abc"},
		{Name: "GetVersionTolerantWithPrefixInConfigNotInTag3", Tags: "1.0.0\n1.1.0\n2.0.0", ExpectedVersion: "2.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: "foo/"},
		{Name: "GetVersionTolerantWithPrefixInConfigInTag", Tags: "1.0.0\n1.1.0\nv2.0.0", ExpectedVersion: "2.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: "v"},

		// Configuring prefixes should not increase tolerance
		{Name: "GetVersionTolerantWithPrefixInConfigButNotMatched1", Tags: "1.0.0\n1.1.0\nabc2.0.0", ExpectedVersion: "1.1.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: "abc"},
		{Name: "GetVersionTolerantWithPrefixInConfigButNotMatched2", Tags: "1.0.0\n1.1.0\nfoo/2.0.0", ExpectedVersion: "1.1.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: false, Prefix: "foo/"},

		// Setting DefaultSemVerPerPrefix to true will increase tolerance
		{Name: "GetVersionTolerantWithSemVerPerPrefix1", Tags: "1.0.0\nfoo/1\nfoo/2.0\nbar/2.1", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: ""},
		{Name: "GetVersionTolerantWithSemVerPerPrefix3", Tags: "v1.0.0\nfoo/1\nfoo/2.0\nbar/2.1", ExpectedVersion: "1.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "v"},
		{Name: "GetVersionTolerantWithSemVerPerPrefix4", Tags: "foo/1.1.0\nbar/2.0.0\nbaz/3.0.0", ExpectedVersion: "1.1.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "foo/"},
		{Name: "GetVersionTolerantWithSemVerPerPrefix5", Tags: "foo/1.1.0\nbar/2.0.0\nbaz/3.0.0", ExpectedVersion: "2.0.0", SemVerParseMode: semverparse.Tolerant, SemVerPerPrefix: true, Prefix: "bar/"},
	}

	for _, test := range tests_with_multiple_tags {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return(test.Tags, nil)
			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: "", GitAPI: gitAPI}

			var got, err = versionAPI.GetVersion(test.SemVerParseMode, test.SemVerPerPrefix, test.Prefix)

			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedVersion, got, `want: "%s, got: "%s"`, test.ExpectedVersion, got)
		})
	}

	type TestErrorWithMultipleTags struct {
		Name   string
		Tags   string
		Prefix string
		Error  error
	}

	var tests_errors_with_multiple_tags = []TestErrorWithMultipleTags{
		{Name: "ReturnError1", Prefix: "", Tags: "foo/1.0.0\nbar/2.0.0\nbaz/3.0.0", Error: fmt.Errorf("could not find a valid semver version")},
		{Name: "ReturnError1", Prefix: "buzz/", Tags: "foo/1.0.0\nbar/2.0.0\nbaz/3.0.0", Error: fmt.Errorf("could not find a valid semver version")},
		{Name: "ReturnError1", Prefix: "buzz/", Tags: "1.0.0\n2.0.0\n3.0.0", Error: fmt.Errorf("could not find a valid semver version")},
	}

	for _, test := range tests_errors_with_multiple_tags {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return(test.Tags, nil)
			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: "", GitAPI: gitAPI}

			var _, got = versionAPI.GetVersion(semverparse.Tolerant, true, test.Prefix)

			assert.Error(t, got)
			assert.Equal(t, test.Error, got, `want: "%s, got: "%s"`, test.Error, got)
		})
	}

	type GitErrorTest struct {
		Error  error
		Name   string
		Prefix string
		Suffix string
	}

	var gitErrorTests = []GitErrorTest{
		{Name: "ReturnErrorOnGitError", Prefix: "", Suffix: "", Error: fmt.Errorf("some-error")},
	}

	for _, test := range gitErrorTests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return("", test.Error)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var _, got = versionAPI.GetVersion(semverparse.Tolerant, false, test.Prefix)

			assert.Error(t, got)
			assert.Equal(t, test.Error, got, `want: "%s, got: "%s"`, test.Error, got)
		})
	}

	type SemverErrorTest struct {
		Error    error
		Name     string
		Prefix   string
		Suffix   string
		Versions string
	}

	var semverErrorTests = []SemverErrorTest{
		{Name: "ReturnErrorOnNoVersions", Versions: "", Error: fmt.Errorf("could not find a valid semver version")},
		{Name: "ReturnErrorOnInvalidVersions", Versions: "invalid1 invalid2", Error: fmt.Errorf("could not find a valid semver version")},
	}

	for _, test := range semverErrorTests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return(test.Versions, nil)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var _, got = versionAPI.GetVersion(semverparse.Tolerant, false, test.Prefix)

			assert.Error(t, got)
			assert.Equal(t, test.Error, got, `want: "%s, got: "%s"`, test.Error, got)
		})
	}
}

func TestAPI_GetVersionOrDefault(t *testing.T) {
	type Test struct {
		Name    string
		Prefix  string
		Suffix  string
		Version string
	}

	var tests = []Test{
		{Name: "GetVersionWithoutError1", Prefix: "", Suffix: "", Version: "0.1.0"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return(test.Version, nil)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var got, err = versionAPI.GetVersion(semverparse.Tolerant, false, test.Prefix)

			assert.NoError(t, err)
			assert.Equal(t, test.Version, got, `want: "%s, got: "%s"`, test.Version, got)
		})
	}

	type ErrorTest struct {
		Error  error
		Prefix string
		Suffix string
		Name   string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnDefaultVersionOnGitApiError", Prefix: "", Suffix: "", Error: fmt.Errorf("some-error")},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return("", test.Error)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var got = versionAPI.GetVersionOrDefault(cli.DefaultVersion, semverparse.Tolerant, false, test.Prefix)

			assert.Equal(t, cli.DefaultVersion, got, `want: "%s, got: "%s"`, cli.DefaultVersion, got)
		})
	}
}

func TestAPI_PredictVersion(t *testing.T) {
	type Test struct {
		Mode    modes.Mode
		Name    string
		Prefix  string
		Suffix  string
		Version string
		Want    string
	}

	var tests = []Test{
		{Name: "ReturnPatchPrediction", Prefix: "", Suffix: "", Mode: modes.NewPatchMode(), Version: "0.0.0", Want: "0.0.1"},
		{Name: "ReturnMinorPrediction", Prefix: "", Suffix: "", Mode: modes.NewMinorMode(), Version: "0.0.0", Want: "0.1.0"},
		{Name: "ReturnMajorPrediction", Prefix: "", Suffix: "", Mode: modes.NewMajorMode(), Version: "0.0.0", Want: "1.0.0"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return(test.Version, nil)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var got, err = versionAPI.PredictVersion(test.Version, test.Mode)

			assert.NoError(t, err)
			assert.Equal(t, test.Want, got, `want: "%s, got: "%s"`, test.Want, got)
		})
	}

	type ErrorTest struct {
		Error   error
		Name    string
		Prefix  string
		Suffix  string
		Version string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnErrorOnModeIncrementError", Prefix: "", Error: fmt.Errorf("some-error"), Version: "invalid"},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix}

			var mode = mocks.NewMockMode()
			mode.On("Increment", mock.Anything, mock.Anything, mock.Anything).Return(test.Version, test.Error)

			var _, got = versionAPI.PredictVersion("0.0.0", mode)

			assert.Error(t, got)
			assert.Equal(t, test.Error, got, `want: "%s, got: "%s"`, test.Error, got)
		})
	}
}

func TestAPI_PushVersion(t *testing.T) {
	type Test struct {
		Mode    modes.Mode
		Name    string
		Prefix  string
		Suffix  string
		Version string
		Want    string
	}

	var tests = []Test{
		{Name: "PushWithoutPrefix", Mode: modes.NewPatchMode(), Prefix: "", Version: "0.0.1", Want: "0.0.1"},
		{Name: "PushWithPrefix", Mode: modes.NewPatchMode(), Prefix: "test/", Version: "0.0.1", Want: "test/0.0.1"},
		{Name: "PushWithoutSuffix", Mode: modes.NewPatchMode(), Prefix: "", Suffix: "", Version: "0.0.1", Want: "0.0.1"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var gitAPI = fakes.NewFakeGitAPI()
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var err = versionAPI.PushVersion(test.Version)

			var pushedTags = versionAPI.GitAPI.(*fakes.FakeGitAPI).PushedTags
			var got = pushedTags[len(pushedTags)-1]

			assert.NoError(t, err)
			assert.Equal(t, test.Want, got, `want: "%s, got: "%s"`, test.Want, got)
		})
	}

	type ErrorTest struct {
		Error   error
		Name    string
		Version string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnErrorOnGitApiError", Error: fmt.Errorf("some-error"), Version: "invalid"},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Run", mock.Anything, mock.Anything).Return(test.Error)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: "v", Suffix: "", GitAPI: gitAPI}

			var got = versionAPI.PushVersion("0.0.1")

			assert.Error(t, got)
			assert.Equal(t, test.Error, got, `want: "%s, got: "%s"`, test.Error, got)
		})
	}
}

func TestAPI_ReleaseVersion(t *testing.T) {
	type Test struct {
		Mode    modes.Mode
		Name    string
		Prefix  string
		Suffix  string
		Version string
		Want    string
	}

	var tests = []Test{
		{Name: "ReleaseWithoutPrefix", Mode: modes.NewPatchMode(), Prefix: "", Version: "0.0.1", Want: "0.0.1"},
		{Name: "ReleaseWithPrefix", Mode: modes.NewPatchMode(), Prefix: "test/", Version: "0.0.1", Want: "test/0.0.1"},
		{Name: "ReleaseWithoutSuffix", Mode: modes.NewPatchMode(), Prefix: "", Suffix: "", Version: "0.0.1", Want: "0.0.1"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var gitAPI = fakes.NewFakeGitAPI()
			var versionAPI = API{Prefix: test.Prefix, Suffix: test.Suffix, GitAPI: gitAPI}

			var err = versionAPI.ReleaseVersion(test.Version)

			var localTags = versionAPI.GitAPI.(*fakes.FakeGitAPI).LocalTags
			var got = localTags[len(localTags)-1]

			assert.NoError(t, err)
			assert.Equal(t, test.Want, got, `want: "%s, got: "%s"`, test.Want, got)
		})
	}

	type ErrorTest struct {
		Error   error
		Name    string
		Version string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnErrorOnGitApiError", Error: fmt.Errorf("some-error"), Version: "invalid"},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Run", mock.Anything, mock.Anything).Return(test.Error)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{Prefix: "v", Suffix: "", GitAPI: gitAPI}

			var got = versionAPI.ReleaseVersion("0.0.1")

			assert.Error(t, got)
			assert.Equal(t, test.Error, got, `want: "%s, got: "%s"`, test.Error, got)
		})
	}
}

func TestAPI_UpdateVersion(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		var gitAPI = fakes.NewFakeGitAPI()
		var versionAPI = API{GitAPI: gitAPI}

		var err = versionAPI.UpdateVersion()

		assert.NoError(t, err)
	})

	type ErrorTest struct {
		Error   error
		Name    string
		Version string
	}

	var errorTests = []ErrorTest{
		{Name: "ReturnErrorOnGitApiError", Error: fmt.Errorf("some-error"), Version: "invalid"},
	}

	for _, test := range errorTests {
		t.Run(test.Name, func(t *testing.T) {
			var cmder = mocks.NewMockCommander()
			cmder.On("Output", mock.Anything, mock.Anything).Return("", test.Error)

			var gitAPI = git.CLI{Commander: cmder}
			var versionAPI = API{GitAPI: gitAPI}

			var got = versionAPI.UpdateVersion()

			assert.Error(t, got)
			assert.Equal(t, test.Error, got, `want: "%s, got: "%s"`, test.Error, got)
		})
	}
}

func TestNewAPI(t *testing.T) {
	t.Run("ValidateState", func(t *testing.T) {
		var api = NewAPI("v", "")
		assert.NotNil(t, api.GitAPI)
	})
}
