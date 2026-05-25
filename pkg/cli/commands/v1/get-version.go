package v1

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/menathan/semverbot/pkg/cli"
	"github.com/menathan/semverbot/pkg/core"
)

// NewGetVersionCommand creates a new get version command.
// Returns the new spf13/cobra command.
func NewGetVersionCommand() *cobra.Command {
	var command = &cobra.Command{
		Use: "version",
		Run: GetVersionCommandRun,
	}

	return command
}

// GetVersionCommandRun runs the command.
func GetVersionCommandRun(cmd *cobra.Command, args []string) {
	log.Debug().Str("command", "v1.get-version").Msg("starting run...")

	var options = &core.GetVersionOptions{
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

	var version = core.GetVersion(options)
	fmt.Println(version)
}
