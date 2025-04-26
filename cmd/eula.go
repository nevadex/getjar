package cmd

import (
	"github.com/nevadex/getjar/ops"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// eulaCmd represents the eula command
var eulaCmd = &cobra.Command{
	Use:   "eula",
	Short: "Accepts the Minecraft EULA",
	Long: `Command to accept the Minecraft EULA in the current directory on the user's behalf
If eula.txt doesn't exist, it will automatically create it

Before using this, please read through the EULA:
https://aka.ms/MinecraftEULA`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if FILENAME != "server.jar" {
			return ops.AcceptEula(FILENAME)
		}

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		return ops.AcceptEula(filepath.Join(dir, "eula.txt"))
	},
}

func init() {
	rootCmd.AddCommand(eulaCmd)
}
