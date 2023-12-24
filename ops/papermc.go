package ops

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	ProjectPaper = "paper"
	ProjectFolia = "folia"
)

func DownloadPaperMC(verbose bool, version string, projectId string, buildId float64, experimental bool, mojmap bool) ([]byte, error) {
	log(verbose, "attempting to download", version, "for", projectId, "using papermc api")

	var versionListRaw map[string]interface{}
	resp, err := http.Get("https://api.papermc.io/v2/projects/" + projectId)
	if err != nil {
		return nil, err
	}
	log(verbose, "retrieved", projectId, "version list json")
	defer func() { _ = resp.Body.Close() }()

	if er := json.NewDecoder(resp.Body).Decode(&versionListRaw); er != nil {
		return nil, err
	}
	versionList := versionListRaw["versions"].([]interface{})
	log(verbose, "decoded version list json")

	if version == "" || version == "latest" {
		version = versionList[len(versionList)-1].(string)
		log(verbose, "no version supplied, defaulting to latest available version:", version)
	} else {
		found := false
		for i := range versionList {
			if version == versionList[i] {
				found = true
				log(verbose, "found version in version list")
				break
			}
		}
		if !found {
			return nil, errors.New("could not locate version in version list")
		}
	}

	var buildListRaw map[string]interface{}
	resp, err = http.Get("https://api.papermc.io/v2/projects/" + projectId + "/versions/" + version + "/builds/")
	if err != nil {
		return nil, err
	}
	log(verbose, "retrieved", projectId, "build list json")
	defer func() { _ = resp.Body.Close() }()

	if er := json.NewDecoder(resp.Body).Decode(&buildListRaw); er != nil {
		return nil, err
	}
	buildList := buildListRaw["builds"].([]interface{})
	log(verbose, "decoded build list json")

	if buildId == 0 {
		for i := range buildList {
			id := buildList[i].(map[string]interface{})["build"].(float64)
			channel := buildList[i].(map[string]interface{})["channel"].(string)
			if experimental && channel == "experimental" {
				buildId = id
			} else if !experimental && channel == "default" {
				buildId = id
			}
		}
		if buildId == 0 {
			if experimental {
				return nil, errors.New("no 'experimental' build found")
			} else if !experimental {
				return nil, errors.New("no 'default' build found")
			}
		}
		if experimental && buildId != buildList[len(buildList)-1].(map[string]interface{})["build"].(float64) {
			fmt.Println("warning! the latest 'experimental' build is not the latest build")
		} else if !experimental && buildId != buildList[len(buildList)-1].(map[string]interface{})["build"].(float64) {
			log(verbose, "warning! the latest 'default' build is not the latest build")
		}
	}
	log(verbose, "using build id", buildId)

	var buildRaw map[string]interface{}
	resp, err = http.Get("https://api.papermc.io/v2/projects/" + projectId + "/versions/" + version + "/builds/" + strconv.FormatFloat(buildId, 'f', -1, 64))
	if err != nil {
		return nil, err
	}
	log(verbose, "retrieved", projectId, "build json")
	defer func() { _ = resp.Body.Close() }()

	if er := json.NewDecoder(resp.Body).Decode(&buildRaw); er != nil {
		return nil, err
	}
	downloads := buildRaw["downloads"].(map[string]interface{})
	jarName := downloads["application"].(map[string]interface{})["name"].(string)
	jarSha256 := downloads["application"].(map[string]interface{})["sha256"].(string)
	if mojmap && downloads["mojang-mappings"] != nil {
		jarName = downloads["mojang-mappings"].(map[string]interface{})["name"].(string)
		jarSha256 = downloads["mojang-mappings"].(map[string]interface{})["sha256"].(string)
	} else if mojmap {
		return nil, errors.New("no mojang-mappings jar is available for this version")
	}
	log(verbose, "decoded build json")

	resp, err = http.Get("https://api.papermc.io/v2/projects/" + projectId + "/versions/" + version + "/builds/" + strconv.FormatFloat(buildId, 'f', -1, 64) + "/downloads/" + jarName)
	if err != nil {
		return nil, err
	}
	log(verbose, "downloaded server jar")
	defer func() { _ = resp.Body.Close() }()

	jar, err := io.ReadAll(resp.Body)

	//log(verbose, "jarfile size:", len(jar))
	//log(verbose, "minecraft version number:", version)
	log(verbose, "sha256 checksum:", jarSha256)

	return jar, nil
}
