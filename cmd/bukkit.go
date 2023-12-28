package cmd

import (
	"getjar/ops"

	"github.com/spf13/cobra"
)

// bukkitCmd represents the bukkit command
var bukkitCmd = &cobra.Command{
	Use:   "bukkit",
	Short: "Compile a CraftBukkit server jar",
	Long: `Command to compile a CraftBukkit minecraft server jar
Downloads and runs BuildTools from SpigotMC
Supports only versions provided by official CraftBukkit channels

All rights for the downloaded content belong to the appropriate persons/organizations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ops.StartLog(VERBOSE)
		fver, err := ops.RunBuildTools(VERBOSE, BUILDTOOLS_VERBOSE, ops.BuildCraftBukkit, VERSION, FILENAME, BUILDTOOLS_EXPERIMENTAL)
		if err != nil {
			return err
		}
		ops.EndLog("saved craftbukkit", fver, "server jar at", FILENAME)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(bukkitCmd)

	bukkitCmd.Flags().BoolVar(&BUILDTOOLS_VERBOSE, "print-buildtools", false, "Output the stdout from BuildTools")
	bukkitCmd.Flags().BoolVarP(&BUILDTOOLS_EXPERIMENTAL, "experimental", "e", false, "Build from experimental source code")
}
