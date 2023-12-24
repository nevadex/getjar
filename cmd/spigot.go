package cmd

import (
	"getjar/ops"

	"github.com/spf13/cobra"
)

// spigotCmd represents the spigot command
var spigotCmd = &cobra.Command{
	Use:   "spigot",
	Short: "Download a SpigotMC server jar",
	Long: `Command to download a SpigotMC minecraft server jar
Supports only versions provided by official SpigotMC channels
Downloads and runs BuildTools from SpigotMC

All rights for the downloaded content belong to the appropriate persons/organizations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ops.RunBuildTools(VERBOSE, BUILDTOOLS_VERBOSE, ops.BuildSpigot, VERSION, FILENAME, BUILDTOOLS_EXPERIMENTAL)
	},
}

func init() {
	rootCmd.AddCommand(spigotCmd)

	spigotCmd.Flags().BoolVar(&BUILDTOOLS_VERBOSE, "print-buildtools", false, "Output the stdout from BuildTools")
	spigotCmd.Flags().BoolVarP(&BUILDTOOLS_EXPERIMENTAL, "experimental", "e", false, "Build from experimental source code")
}
