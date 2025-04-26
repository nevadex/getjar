package ops

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

const (
	ProjectMohist = "mohist"
	ProjectBanner = "banner"
)

func DownloadMohistMC(version string, projectId string, buildId float64) ([]byte, string, error) {
	log("attempting to download", version, "using mohistmc api")

	slog("downloading version list")
	var versionListRaw map[string]interface{}
	resp, err := http.Get("https://mohistmc.com/api/v2/projects/" + projectId)
	if err != nil {
		return nil, "", err
	}
	log("retrieved version list json")
	defer func() { _ = resp.Body.Close() }()

	if er := json.NewDecoder(resp.Body).Decode(&versionListRaw); er != nil {
		return nil, "", err
	}
	versionList := versionListRaw["versions"].([]interface{})
	log("decoded version list json")

	if version == "latest" {
		version = versionList[len(versionList)-1].(string)
		log("no version supplied, defaulting to latest available version:", version)
	} else {
		found := false
		for i := range versionList {
			if version == versionList[i] {
				found = true
				log("found version in version list")
				break
			}
		}
		if !found {
			return nil, "", errors.New("could not locate version in version list")
		}
	}

stupidAF:
	slog("downloading build list")
	var buildListRaw map[string]interface{}
	resp, err = http.Get("https://mohistmc.com/api/v2/projects/" + projectId + "/" + version + "/builds/")
	if err != nil {
		return nil, "", err
	}
	log("retrieved", version, "build list json")
	defer func() { _ = resp.Body.Close() }()

	if er := json.NewDecoder(resp.Body).Decode(&buildListRaw); er != nil {
		return nil, "", err
	}
	buildList := buildListRaw["builds"].([]interface{})
	log("decoded build list json")

	// mohist sometimes puts a version up but makes no builds for it, address that here (STUPID)
	if len(buildList) == 0 {
		log("no builds available for version", version)
		versionList = versionList[:len(versionList)-1]
		version = versionList[len(versionList)-1].(string)
		log("trying again with version", version)
		goto stupidAF
	}

	// mohistmc alternates between "number" integer ids and "id" sha1 from git (HORRIBLE POS)
	realBuildIdentifier := ""
	if buildId == 0 {
		for i := range buildList {
			supposedBuildNumber := buildList[i].(map[string]interface{})["number"]
			if supposedBuildNumber != nil {
				realBuildIdentifier = strconv.FormatFloat(supposedBuildNumber.(float64), 'f', -1, 64)
				continue
			}

			supposedBuildId := buildList[i].(map[string]interface{})["id"]
			if supposedBuildId != nil {
				realBuildIdentifier = supposedBuildId.(string)
				continue
			}
		}
	} else {
		realBuildIdentifier = strconv.FormatFloat(buildId, 'f', -1, 64)
	}
	log("using build id", realBuildIdentifier)

	slog("downloading build", realBuildIdentifier)
	var buildRaw map[string]interface{}
	resp, err = http.Get("https://mohistmc.com/api/v2/projects/" + projectId + "/" + version + "/builds/" + realBuildIdentifier)
	if err != nil {
		return nil, "", err
	}
	log("retrieved build", realBuildIdentifier, "json")
	defer func() { _ = resp.Body.Close() }()

	if err = json.NewDecoder(resp.Body).Decode(&buildRaw); err != nil {
		return nil, "", err
	}
	downloads := buildRaw["build"].(map[string]interface{})
	jarUrl := downloads["originUrl"].(string)
	jarMd5 := downloads["fileMd5"].(string)
	jarSha256 := downloads["fileSha256"].(string)
	log("decoded build json")

	slog("downloading server jar")
	resp, err = http.Get(jarUrl)
	if err != nil {
		return nil, "", err
	}
	log("downloaded server jar")
	defer func() { _ = resp.Body.Close() }()

	jar, err := io.ReadAll(resp.Body)

	//log("jarfile size:", len(jar))
	//log("minecraft version number:", version)
	if printChecksum {
		post("md5 checksum:", jarMd5)
		post("sha256 checksum:", jarSha256)
	}

	return jar, version, nil
}

func GetVersionListMohistMC(projectId string) ([]string, error) {
	var raw map[string]interface{}
	resp, err := http.Get("https://mohistmc.com/api/v2/projects/" + projectId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if err = json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	versionList := raw["versions"].([]interface{})
	var list []string
	for i := range versionList {
		list = append(list, versionList[i].(string))
	}

	return list, nil
}
