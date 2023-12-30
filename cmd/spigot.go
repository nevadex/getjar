package cmd

import (
	"github.com/nevadex/getjar/ops"

	"github.com/spf13/cobra"
)

// spigotCmd represents the spigot command
var spigotCmd = &cobra.Command{
	Use:   "spigot",
	Short: "Compile a SpigotMC server jar",
	Long: `Command to compile a SpigotMC minecraft server jar
Downloads and runs BuildTools from SpigotMC
Supports only versions provided by official SpigotMC channels

All rights for the downloaded content belong to the appropriate persons/organizations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ops.StartLog(VERBOSE)
		fver, err := ops.RunBuildTools(VERBOSE, BUILDTOOLS_VERBOSE, ops.BuildSpigot, VERSION, FILENAME, BUILDTOOLS_EXPERIMENTAL)
		if err != nil {
			return err
		}
		ops.EndLog("saved spigotmc", fver, "server jar at", FILENAME)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(spigotCmd)

	spigotCmd.Flags().BoolVar(&BUILDTOOLS_VERBOSE, "print-buildtools", false, "Output the stdout from BuildTools")
	spigotCmd.Flags().BoolVarP(&BUILDTOOLS_EXPERIMENTAL, "experimental", "e", false, "Build from experimental source code")
}
