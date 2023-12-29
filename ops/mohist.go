package ops

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func DownloadMohist(version string, buildId float64) ([]byte, string, error) {
	log("attempting to download", version, "using mohist api")

	slog("downloading version list")
	var versionListRaw map[string]interface{}
	resp, err := http.Get("https://mohistmc.com/api/v2/projects/mohist/")
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

	slog("downloading build list")
	var buildListRaw map[string]interface{}
	resp, err = http.Get("https://mohistmc.com/api/v2/projects/mohist/" + version + "/builds/")
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

	if buildId == 0 {
		for i := range buildList {
			buildId = buildList[i].(map[string]interface{})["number"].(float64)
		}
	}
	log("using build id", buildId)

	slog("downloading build", buildId)
	var buildRaw map[string]interface{}
	resp, err = http.Get("https://mohistmc.com/api/v2/projects/mohist/" + version + "/builds/" + strconv.FormatFloat(buildId, 'f', -1, 64))
	if err != nil {
		return nil, "", err
	}
	log("retrieved build", buildId, "json")
	defer func() { _ = resp.Body.Close() }()

	if err = json.NewDecoder(resp.Body).Decode(&buildRaw); err != nil {
		return nil, "", err
	}
	downloads := buildRaw["build"].(map[string]interface{})
	jarUrl := downloads["originUrl"].(string)
	jarMd5 := downloads["fileMd5"].(string)
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
	log("md5 checksum:", jarMd5)

	return jar, version, nil
}
