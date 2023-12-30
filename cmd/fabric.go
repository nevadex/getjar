package cmd

import (
	"github.com/nevadex/getjar/ops"
	"os"

	"github.com/spf13/cobra"
)

// fabricCmd represents the fabric command
var fabricCmd = &cobra.Command{
	Use:   "fabric",
	Short: "Download a Fabric server jar",
	Long: `Command to download a Fabric minecraft server jar
Supports only versions provided by official FabricMC channels
Downloads directly from FabricMC's Meta API

All rights for the downloaded content belong to the appropriate persons/organizations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ops.StartLog(VERBOSE)
		jar, fver, err := ops.DownloadFabricMC(VERSION, FABRIC_FABRIC_VERSION, FABRIC_INSTALLER_VERSION, FABRIC_EXPERIMENTAL)
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
		ops.EndLog("saved fabric", fver, "server jar at", FILENAME)

		return file.Close()
	},
}

func init() {
	rootCmd.AddCommand(fabricCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// paperCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// paperCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	fabricCmd.Flags().StringVarP(&FABRIC_FABRIC_VERSION, "fabric-version", "b", "latest", "Supply a specific fabric version instead of using the latest stable version")
	fabricCmd.Flags().StringVar(&FABRIC_INSTALLER_VERSION, "installer-version", "latest", "Supply a specific fabric installer version instead of using the latest stable version")
	fabricCmd.Flags().BoolVarP(&FABRIC_EXPERIMENTAL, "unstable", "e", false, "Download latest unstable fabric version instead of latest stable fabric version")
}
