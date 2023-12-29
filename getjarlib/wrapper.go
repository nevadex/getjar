package getjarlib

import (
	"getjar/ops"
	"os"
)

const (
	BANNER    = 1
	BUKKIT    = 2
	CATSERVER = 3
	FABRIC    = 4
	FOLIA     = 5
	MOHIST    = 6
	PAPER     = 7
	PURPUR    = 8
	SPIGOT    = 9
	VANILLA   = 10
)

//goland:noinspection GoUnusedExportedFunction
func ServerJar(version string, filename string, serverType int) error {
	if serverType != BUKKIT && serverType != SPIGOT {
		if version == "" {
			version = "latest"
		}

		var jar []byte
		var err error

		switch serverType {
		case BANNER:
			jar, _, err = ops.DownloadMohistMC(version, ops.ProjectBanner, 0)
		case CATSERVER:
			jar, _, err = ops.DownloadCatserver(version)
		case FABRIC:
			jar, _, err = ops.DownloadFabricMC(version, "latest", "latest", false)
		case FOLIA:
			jar, _, err = ops.DownloadPaperMC(version, ops.ProjectFolia, 0, false, false)
		case MOHIST:
			jar, _, err = ops.DownloadMohistMC(version, ops.ProjectMohist, 0)
		case PAPER:
			jar, _, err = ops.DownloadPaperMC(version, ops.ProjectPaper, 0, false, false)
		case PURPUR:
			jar, _, err = ops.DownloadPurpurMC(version, "latest")
		case VANILLA:
			jar, _, err = ops.DownloadVanilla(version)
		}

		if err != nil {
			return err
		}

		file, err := os.Create(filename)
		if err != nil {
			return err
		}

		_, err = file.Write(jar)
		if err != nil {
			return err
		}

		return file.Close()
	} else {
		var err error

		switch serverType {
		case BUKKIT:
			_, err = ops.RunBuildTools(false, false, ops.BuildCraftBukkit, version, filename, false)
		case SPIGOT:
			_, err = ops.RunBuildTools(false, false, ops.BuildSpigot, version, filename, false)
		}

		return err
	}
}
