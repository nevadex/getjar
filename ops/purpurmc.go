package ops

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func DownloadPurpurMC(version string, buildId string) ([]byte, string, error) {
	log("attempting to download", version, "using purpurmc api")

	slog("downloading version list")
	var versionListRaw map[string]interface{}
	resp, err := http.Get("https://api.purpurmc.org/v2/purpur/")
	if err != nil {
		return nil, "", err
	}
	log("retrieved version list json")
	defer func() { _ = resp.Body.Close() }()

	if err = json.NewDecoder(resp.Body).Decode(&versionListRaw); err != nil {
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
	var buildList map[string]interface{}
	resp, err = http.Get("https://api.purpurmc.org/v2/purpur/" + version)
	if err != nil {
		return nil, "", err
	}
	log("retrieved", version, "build list json")
	defer func() { _ = resp.Body.Close() }()

	if err = json.NewDecoder(resp.Body).Decode(&buildList); err != nil {
		return nil, "", err
	}
	log("decoded build list json")

	if buildId == "latest" {
		buildId = buildList["builds"].(map[string]interface{})["latest"].(string)
	}
	log("using build id", buildId)

	slog("downloading build", buildId)
	var build map[string]interface{}
	resp, err = http.Get("https://api.purpurmc.org/v2/purpur/" + version + "/" + buildId)
	if err != nil {
		return nil, "", err
	}
	log("retrieved build", buildId, "json")
	defer func() { _ = resp.Body.Close() }()

	if err = json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return nil, "", err
	}
	jarMd5 := build["md5"].(string)
	log("decoded build json")

	slog("downloading server jar")
	resp, err = http.Get("https://api.purpurmc.org/v2/purpur/" + version + "/" + buildId + "/download")
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

func GetVersionListPurpurMC() ([]string, error) {
	var versionListRaw map[string]interface{}
	resp, err := http.Get("https://api.purpurmc.org/v2/purpur/")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if err = json.NewDecoder(resp.Body).Decode(&versionListRaw); err != nil {
		return nil, err
	}
	versionList := versionListRaw["versions"].([]interface{})

	var list []string
	for i := range versionList {
		list = append(list, versionList[i].(string))
	}
	return list, nil
}
