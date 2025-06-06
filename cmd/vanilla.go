package cmd

import (
	"github.com/nevadex/getjar/ops"
	"os"

	"github.com/spf13/cobra"
)

// vanillaCmd represents the vanilla command
var vanillaCmd = &cobra.Command{
	Use:   "vanilla",
	Short: "Download a Vanilla server jar",
	Long: `Command to download a Vanilla minecraft server jar
Supports any Minecraft: Java Edition version

All rights for the downloaded content belong to:
Mojang AB, Microsoft Corporation, or one of its local affiliates
EULA: https://aka.ms/MinecraftEULA`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ops.StartLog(VERBOSE, CHECKSUM)
		jar, fver, err := ops.DownloadVanilla(VERSION, VANILLA_SNAPSHOT)
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
		ops.EndLog("saved vanilla", fver, "server jar at", FILENAME)

		return file.Close()
	},
}

func init() {
	rootCmd.AddCommand(vanillaCmd)
	vanillaCmd.Flags().BoolVarP(&VANILLA_SNAPSHOT, "snapshot", "e", false, "Download latest snapshot version instead of latest release")
}
