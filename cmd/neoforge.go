package cmd

import (
	"github.com/nevadex/getjar/ops"
	"os"

	"github.com/spf13/cobra"
)

// neoforgeCmd represents the catserver command
var neoforgeCmd = &cobra.Command{
	Use:   "neoforge",
	Short: "Download a Neoforge server installer",
	Long: `Command to download a Neoforge minecraft server installer
Supports only versions provided by official Neoforged channels
Downloads directly from Neoforged Maven

All rights for the downloaded content belong to the appropriate persons/organizations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ops.StartLog(VERBOSE, CHECKSUM)

		jar, fver, err := ops.DownloadNeoforged(VERSION, NEOFORGED_NEO_VERSION, NEOFORGED_EXPERIMENTAL)
		if err != nil {
			return err
		}

		file, err := os.Create(FILENAME)
		if err != nil {
			return err
		}

		_, err = file.Write(jar)
		if err != nil {
			return err
		}
		ops.EndLog("saved neoforge", fver, "server installer at", FILENAME)

		return file.Close()
	},
}

func init() {
	rootCmd.AddCommand(neoforgeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// neoforgeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// neoforgeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	neoforgeCmd.Flags().StringVarP(&NEOFORGED_NEO_VERSION, "neo-version", "b", "latest", "Supply a specific neo version instead of using the latest recommended version for this version of minecraft")
	neoforgeCmd.Flags().BoolVarP(&NEOFORGED_EXPERIMENTAL, "experimental", "e", false, "Download latest (possibly unstable) version instead of latest recommended version")
}
