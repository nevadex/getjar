package cmd

import (
	"github.com/nevadex/getjar/ops"
	"os"

	"github.com/spf13/cobra"
)

// forgeCmd represents the forge command
var forgeCmd = &cobra.Command{
	Use:   "forge",
	Short: "Download a Forge server installer",
	Long: `Command to download a Forge minecraft server installer
Supports only versions provided by official Forge channels
Downloads directly from MinecraftForge Maven

All rights for the downloaded content belong to the appropriate persons/organizations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ops.StartLog(VERBOSE, CHECKSUM)
		jar, fver, err := ops.DownloadMinecraftForge(VERSION, FORGE_FORGE_VERSION, FORGE_EXPERIMENTAL)
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
		ops.EndLog("saved forge", fver, "server installer at", FILENAME)

		return file.Close()
	},
}

func init() {
	rootCmd.AddCommand(forgeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exampleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exampleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	forgeCmd.Flags().StringVarP(&FORGE_FORGE_VERSION, "forge-version", "b", "latest", "Supply a specific forge version instead of using the latest recommended version")
	forgeCmd.Flags().BoolVarP(&FORGE_EXPERIMENTAL, "experimental", "e", false, "Download latest (possibly unstable) version instead of latest recommended version")
}
