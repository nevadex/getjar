package cmd

import (
	"fmt"
	"github.com/nevadex/getjar/ops"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "getjar",
	Short: "Quickly get popular minecraft server jars",
	Long: `A tool to automatically download or compile Minecraft server jars
Uses official channels from development teams to acquire files
No middleman API is used!

https://github.com/nevadex/getjar/`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if QUIET {
			os.Stdout = nil
		}

		if LIST_VERSIONS {
			var versionList []string
			var err error

			ops.StartLog(VERBOSE, false)
			switch cmd {
			case bannerCmd:
				versionList, err = ops.GetVersionListMohistMC(ops.ProjectBanner)
			case bukkitCmd:
				versionList, err = ops.GetVersionListBuildTools()
			case catserverCmd:
				versionList, err = ops.GetVersionListCatserver()
			case fabricCmd:
				versionList, err = ops.GetVersionListFabricMC()
			case foliaCmd:
				versionList, err = ops.GetVersionListPaperMC(ops.ProjectFolia)
			case forgeCmd:
				versionList, err = ops.GetVersionListMinecraftForge()
			case mohistCmd:
				versionList, err = ops.GetVersionListMohistMC(ops.ProjectMohist)
			case neoforgeCmd:
				versionList, err = ops.GetVersionListNeoforged()
			case paperCmd:
				versionList, err = ops.GetVersionListPaperMC(ops.ProjectPaper)
			case purpurCmd:
				versionList, err = ops.GetVersionListPurpurMC()
			case spigotCmd:
				versionList, err = ops.GetVersionListBuildTools()
			case vanillaCmd:
				versionList, err = ops.GetVersionListVanilla()
			}

			if FILENAME == "server.jar" {
				ops.EndLog("downloaded version list")

				for i := range versionList {
					fmt.Println(versionList[i])
				}

				fmt.Println("(total", len(versionList), "versions)")
			} else {
				var str string
				for i := range versionList {
					str += fmt.Sprintln(versionList[i])
				}

				err = os.WriteFile(FILENAME, []byte(str), os.ModePerm)
				ops.EndLog("saved version list to", FILENAME)
			}

			os.Exit(0)
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if PRINT_VERSION {
			fmt.Println("1.1.0")
			return nil
		} else {
			return cmd.Help()
		}
	},
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
	PRINT_VERSION bool

	VERSION  string
	FILENAME string

	VERBOSE       bool
	QUIET         bool
	LIST_VERSIONS bool
	CHECKSUM      bool

	BUILDTOOLS_VERBOSE      bool
	BUILDTOOLS_EXPERIMENTAL bool

	PAPERMC_BUILD_ID        int
	PAPERMC_EXPERIMENTAL    bool
	PAPERMC_MOJANG_MAPPINGS bool

	MOHISTMC_BUILD_ID int

	FABRIC_FABRIC_VERSION    string
	FABRIC_INSTALLER_VERSION string
	FABRIC_EXPERIMENTAL      bool

	PURPURMC_BUILD_ID string

	FORGE_FORGE_VERSION string
	FORGE_EXPERIMENTAL  bool

	NEOFORGED_NEO_VERSION  string
	NEOFORGED_EXPERIMENTAL bool

	VANILLA_SNAPSHOT bool
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
	rootCmd.PersistentFlags().BoolVarP(&LIST_VERSIONS, "list", "l", false, "Output a list of minecraft versions supported by a server")
	rootCmd.PersistentFlags().BoolVarP(&QUIET, "quiet", "q", false, "Output nothing to stdout (sets stdout to nil)")
	rootCmd.PersistentFlags().BoolVarP(&CHECKSUM, "checksum", "c", false, "Output the jar checksum if available")

	rootCmd.Flags().BoolVarP(&PRINT_VERSION, "version", "v", false, "Output the version of this build of getjar")
}
