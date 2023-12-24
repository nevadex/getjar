package cmd

import (
	"fmt"
	"getjar/ops"
	"os"

	"github.com/spf13/cobra"
)

// vanillaCmd represents the vanilla command
var vanillaCmd = &cobra.Command{
	Use:   "vanilla",
	Short: "Download a Vanilla Minecraft server jar",
	Long: `Command to download a Vanilla Minecraft server jar
Supports any Minecraft: Java Edition version

All rights for the downloaded content belong to:
Mojang AB, Microsoft Corporation, or one of its local affiliates
EULA: https://aka.ms/MinecraftEULA`,
	RunE: func(cmd *cobra.Command, args []string) error {
		jar, err := ops.DownloadVanilla(VERBOSE, VERSION)
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
		fmt.Println("saved vanilla server jar at", FILENAME)

		return file.Close()
	},
}

func init() {
	rootCmd.AddCommand(vanillaCmd)
}
