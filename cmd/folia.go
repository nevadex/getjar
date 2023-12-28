package cmd

import (
	"getjar/ops"
	"os"

	"github.com/spf13/cobra"
)

// foliaCmd represents the folia command
var foliaCmd = &cobra.Command{
	Use:   "folia",
	Short: "Download a Folia Minecraft server jar",
	Long: `Command to download a Folia Minecraft server jar
Supports only versions provided by official PaperMC channels
Downloads directly from PaperMC's Downloads API

All rights for the downloaded content belong to the appropriate persons/organizations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ops.StartLog(VERBOSE)
		jar, fver, err := ops.DownloadPaperMC(VERSION, ops.ProjectFolia, float64(PAPERMC_BUILD_ID), PAPERMC_EXPERIMENTAL, PAPERMC_MOJANG_MAPPINGS)
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
		ops.EndLog("saved folia", fver, "server jar at", FILENAME)

		return file.Close()
	},
}

func init() {
	rootCmd.AddCommand(foliaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// foliaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// foliaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	foliaCmd.Flags().IntVarP(&PAPERMC_BUILD_ID, "build-id", "b", 0, "Supply a specific build id instead of using the latest 'default' build")
	foliaCmd.Flags().BoolVarP(&PAPERMC_EXPERIMENTAL, "experimental", "e", false, "Download latest 'experimental' build instead of latest 'default' build")
	foliaCmd.Flags().BoolVarP(&PAPERMC_MOJANG_MAPPINGS, "mojang-mappings", "m", false, "Download the mojmap jar instead of the regular jar")
}
