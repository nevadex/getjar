package cmd

import (
	"getjar/ops"

	"github.com/spf13/cobra"
)

// bukkitCmd represents the bukkit command
var bukkitCmd = &cobra.Command{
	Use:   "bukkit",
	Short: "Download a CraftBukkit server jar",
	Long: `Command to download a CraftBukkit minecraft server jar
Supports only versions provided by official CraftBukkit channels
Downloads and runs BuildTools from SpigotMC

All rights for the downloaded content belong to the appropriate persons/organizations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ops.RunBuildTools(VERBOSE, BUILDTOOLS_VERBOSE, ops.BuildCraftBukkit, VERSION, FILENAME, BUILDTOOLS_EXPERIMENTAL)
	},
}

func init() {
	rootCmd.AddCommand(bukkitCmd)

	bukkitCmd.Flags().BoolVar(&BUILDTOOLS_VERBOSE, "print-buildtools", false, "Output the stdout from BuildTools")
	bukkitCmd.Flags().BoolVarP(&BUILDTOOLS_EXPERIMENTAL, "experimental", "e", false, "Build from experimental source code")
}
