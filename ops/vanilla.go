package ops

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
)

func DownloadVanilla(version string, snapshot bool) ([]byte, string, error) {
	log("attempting to download vanilla", version)

	slog("downloading version manifest")
	var versionManifest map[string]interface{}
	resp, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
	if err != nil {
		return nil, "", err
	}
	log("downloaded version_manifest.json")
	defer func() { _ = resp.Body.Close() }()

	err = json.NewDecoder(resp.Body).Decode(&versionManifest)
	if err != nil {
		return nil, "", err
	}
	log("decoded version_manifest.json")

	if version == "" || version == "latest" {
		if snapshot {
			version = versionManifest["latest"].(map[string]interface{})["snapshot"].(string)
		} else {
			version = versionManifest["latest"].(map[string]interface{})["release"].(string)
		}
		log("no version supplied, defaulting to latest:", version)
	}

	var versionUrl string
	var versionJson map[string]interface{}

	for _, v := range versionManifest["versions"].([]interface{}) {
		if v.(map[string]interface{})["id"].(string) == version {
			versionUrl = v.(map[string]interface{})["url"].(string)
			log("located version in version_manifest.json")
			break
		}
	}
	if versionUrl == "" {
		return nil, "", errors.New("could not locate version in version_manifest.json")
	}

	slog("downloading jar manifest")
	resp, err = http.Get(versionUrl)
	if err != nil {
		return nil, "", err
	}
	log("downloaded jar manifest")
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	err = json.Unmarshal(body, &versionJson)
	if err != nil {
		return nil, "", err
	}
	serverDownloadsSection := versionJson["downloads"].(map[string]interface{})["server"]
	if serverDownloadsSection == nil {
		return nil, "", errors.New("provided version does not have a server jar")
	}
	jarurl := serverDownloadsSection.(map[string]interface{})["url"].(string)
	//jarsize := strconv.FormatFloat(versionJson["downloads"].(map[string]interface{})["server"].(map[string]interface{})["size"].(float64), 'f', -1, 64)
	jarsha1 := versionJson["downloads"].(map[string]interface{})["server"].(map[string]interface{})["sha1"].(string)
	//jarversiontype := versionJson["type"].(string)
	//jarjavaversion := strconv.FormatFloat(versionJson["javaVersion"].(map[string]interface{})["majorVersion"].(float64), 'f', -1, 64)
	log("decoded jar manifest")

	slog("downloading jar")
	jarresp, err := http.Get(jarurl)
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = resp.Body.Close() }()
	jar, err := io.ReadAll(jarresp.Body)
	log("downloaded server jar")

	//log("jarfile size:", jarsize)
	//log("minecraft version type:", jarversiontype)
	//log("minecraft version number:", version)
	//log("minimum java version:", jarjavaversion)
	if printChecksum {
		post("sha1 checksum:", jarsha1)
	}

	return jar, version, nil
}

func GetVersionListVanilla() ([]string, error) {
	var versionManifest map[string]interface{}
	resp, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	err = json.NewDecoder(resp.Body).Decode(&versionManifest)
	if err != nil {
		return nil, err
	}
	versionList := versionManifest["versions"].([]interface{})

	var list []string
	for i := range versionList {
		list = append(list, versionList[i].(map[string]interface{})["id"].(string))
	}
	slices.Reverse(list)
	return list, nil
}
