package v1

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/restechnica/semverbot/pkg/cli"
	"github.com/restechnica/semverbot/pkg/core"
)

// NewPushVersionCommand creates a new push version command.
// Returns the new spf13/cobra command.
func NewPushVersionCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:  "version",
		RunE: PushVersionCommandRunE,
	}

	return command
}

// PushVersionCommandRunE runs the command.
// Returns an error if the command fails.
func PushVersionCommandRunE(cmd *cobra.Command, args []string) (err error) {
	log.Debug().Str("command", "v1.push-version").Msg("starting run...")

	var options = &core.PushVersionOptions{
		DefaultVersion:  cli.DefaultVersion,
		Prefix:          viper.GetString(cli.GitTagsPrefixConfigKey),
		Suffix:          viper.GetString(cli.GitTagsSuffixConfigKey),
		SemVerParseMode: viper.GetString(cli.SemVerParseModeConfigKey),
		SemVerPerPrefix: viper.GetBool(cli.SemVerPerPrefixConfigKey),
	}

	log.Debug().
		Str("default", options.DefaultVersion).
		Str("prefix", options.Prefix).
		Str("suffix", options.Suffix).
		Str("semver-parse-mode", options.SemVerParseMode).
		Bool("semver-per-prefix", options.SemVerPerPrefix).
		Msg("options")

	if err = core.PushVersion(options); err != nil {
		err = cli.NewCommandError(err)
	}

	return err
}
