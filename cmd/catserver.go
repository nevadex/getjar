package cmd

import (
	"getjar/ops"
	"os"

	"github.com/spf13/cobra"
)

// catserverCmd represents the catserver command
var catserverCmd = &cobra.Command{
	Use:   "catserver",
	Short: "Download a Catserver server jar",
	Long: `Command to download a Catserver minecraft server jar
Supports only versions provided by official Catserver channels
Downloads directly from Catserver Jenkins

All rights for the downloaded content belong to the appropriate persons/organizations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ops.StartLog(VERBOSE)

		jar, fver, err := ops.DownloadCatserver(VERSION)
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
		ops.EndLog("saved catserver", fver, "server jar at", FILENAME)

		return file.Close()
	},
}

func init() {
	rootCmd.AddCommand(catserverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catserverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// catserverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
