package cmd

import (
	"getjar/ops"
	"os"

	"github.com/spf13/cobra"
)

// mohistCmd represents the paper command
var mohistCmd = &cobra.Command{
	Use:   "mohist",
	Short: "Download a Mohist Minecraft server jar",
	Long: `Command to download a Mohist Minecraft server jar
Supports only versions provided by official Mohist channels
Downloads directly from Mohist's Downloads API

All rights for the downloaded content belong to the appropriate persons/organizations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ops.StartLog(VERBOSE)
		jar, fver, err := ops.DownloadMohist(VERSION, float64(MOHIST_BUILD_ID))
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
		ops.EndLog("saved mohist", fver, "server jar at", FILENAME)

		return file.Close()
	},
}

func init() {
	rootCmd.AddCommand(mohistCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// paperCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// paperCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	mohistCmd.Flags().IntVarP(&MOHIST_BUILD_ID, "build-id", "b", 0, "Supply a specific build id instead of using the latest 'default' build")
}
