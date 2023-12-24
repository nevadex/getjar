package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "getjar",
	Short: "Quickly get popular minecraft server jars",
	Long: `A tool to automatically download or compile Minecraft server jars and proxies
Uses official channels from development teams to acquire files
No middleman API is used!

https://github.com/nevadex/getjar/`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	VERSION  string
	FILENAME string

	VERBOSE bool

	BUILDTOOLS_VERBOSE      bool
	BUILDTOOLS_EXPERIMENTAL bool

	PAPERMC_BUILD_ID        int
	PAPERMC_EXPERIMENTAL    bool
	PAPERMC_MOJANG_MAPPINGS bool
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.getjar.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringVarP(&VERSION, "version", "v", "latest", "The version of minecraft jar")
	rootCmd.PersistentFlags().StringVarP(&FILENAME, "filename", "f", "server.jar", "The name of the jarfile to write")

	rootCmd.PersistentFlags().BoolVarP(&VERBOSE, "verbose", "p", false, "Output as much detail as possible")
}
