package cmd

import (
	"github.com/nevadex/getjar/ops"
	"os"

	"github.com/spf13/cobra"
)

// bannerCmd represents the banner command
var bannerCmd = &cobra.Command{
	Use:   "banner",
	Short: "Download a Banner server jar",
	Long: `Command to download a Banner Minecraft server jar
Supports only versions provided by official MohistMC channels
Downloads directly from MohistMC's Downloads API

All rights for the downloaded content belong to the appropriate persons/organizations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ops.StartLog(VERBOSE, CHECKSUM)
		jar, fver, err := ops.DownloadMohistMC(VERSION, ops.ProjectBanner, float64(MOHISTMC_BUILD_ID))
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
		ops.EndLog("saved banner", fver, "server jar at", FILENAME)

		return file.Close()
	},
}

func init() {
	rootCmd.AddCommand(bannerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// paperCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// paperCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	bannerCmd.Flags().IntVarP(&MOHISTMC_BUILD_ID, "build-id", "b", 0, "Supply a specific build id instead of using the latest 'default' build")
}
