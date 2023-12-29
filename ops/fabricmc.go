package ops

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func DownloadFabricMC(version string, fabricVersion string, installerVersion string, experimental bool) ([]byte, string, error) {
	log("attempting to download", version, "using fabricmc api")

	slog("downloading version list")
	var versionList []interface{}
	resp, err := http.Get("https://meta.fabricmc.net/v2/versions/game")
	if err != nil {
		return nil, "", err
	}
	log("retrieved version list json")
	defer func() { _ = resp.Body.Close() }()

	if err = json.NewDecoder(resp.Body).Decode(&versionList); err != nil {
		return nil, "", err
	}
	log("decoded version list json")

	if version == "latest" {
		for i := range versionList {
			if versionList[i].(map[string]interface{})["stable"].(bool) {
				version = versionList[i].(map[string]interface{})["version"].(string)
				break
			}
		}
		log("no version supplied, defaulting to latest available version:", version)
	} else {
		found := false
		for i := range versionList {
			if version == versionList[i].(map[string]interface{})["version"].(string) {
				found = true
				log("found version in version list")
				break
			}
		}
		if !found {
			return nil, "", errors.New("could not locate version in version list")
		}
	}

	slog("downloading installer version list")
	var installerVersionList []interface{}
	resp, err = http.Get("https://meta.fabricmc.net/v2/versions/installer")
	if err != nil {
		return nil, "", err
	}
	log("retrieved installer version list json")
	defer func() { _ = resp.Body.Close() }()

	if err = json.NewDecoder(resp.Body).Decode(&installerVersionList); err != nil {
		return nil, "", err
	}
	log("decoded installer version list json")

	if installerVersion == "latest" {
		for i := range installerVersionList {
			if installerVersionList[i].(map[string]interface{})["stable"].(bool) {
				installerVersion = installerVersionList[i].(map[string]interface{})["version"].(string)
				break
			}
		}
		log("no installer version supplied, defaulting to latest available installer version:", installerVersion)
	} else {
		found := false
		for i := range installerVersionList {
			if installerVersion == installerVersionList[i].(map[string]interface{})["version"].(string) {
				found = true
				log("found installer version in installer version list")
				break
			}
		}
		if !found {
			return nil, "", errors.New("could not locate installer version in installer version list")
		}
	}

	slog("downloading fabric version list")
	var fabricVersionList []interface{}
	resp, err = http.Get("https://meta.fabricmc.net/v2/versions/loader/" + version)
	if err != nil {
		return nil, "", err
	}
	log("retrieved fabric version list json")
	defer func() { _ = resp.Body.Close() }()

	if err = json.NewDecoder(resp.Body).Decode(&fabricVersionList); err != nil {
		return nil, "", err
	}
	log("decoded fabric version list json")

	if fabricVersion == "latest" {
		for i := range fabricVersionList {
			if fabricVersionList[i].(map[string]interface{})["loader"].(map[string]interface{})["stable"].(bool) && !experimental {
				fabricVersion = fabricVersionList[i].(map[string]interface{})["loader"].(map[string]interface{})["version"].(string)
				break
			} else if !fabricVersionList[i].(map[string]interface{})["loader"].(map[string]interface{})["stable"].(bool) && experimental {
				fabricVersion = fabricVersionList[i].(map[string]interface{})["loader"].(map[string]interface{})["version"].(string)
				break
			}
		}
		log("no fabric version supplied, defaulting to latest available fabric version:", fabricVersion)
	} else {
		found := false
		for i := range fabricVersionList {
			if fabricVersion == fabricVersionList[i].(map[string]interface{})["loader"].(map[string]interface{})["version"].(string) {
				found = true
				log("found fabric version in fabric version list")
				break
			}
		}
		if !found {
			return nil, "", errors.New("could not locate fabric version in fabric version list")
		}
	}

	slog("downloading server jar")
	resp, err = http.Get("https://meta.fabricmc.net/v2/versions/loader/" + version + "/" + fabricVersion + "/" + installerVersion + "/server/jar")
	if err != nil {
		return nil, "", err
	}
	log("downloaded server jar")
	defer func() { _ = resp.Body.Close() }()

	jar, err := io.ReadAll(resp.Body)

	return jar, version, nil
}

func GetVersionListFabricMC() ([]string, error) {
	var versionList []interface{}
	resp, err := http.Get("https://meta.fabricmc.net/v2/versions/game")
	if err != nil {
		return nil, err
	}
	log("retrieved version list json")
	defer func() { _ = resp.Body.Close() }()

	if err = json.NewDecoder(resp.Body).Decode(&versionList); err != nil {
		return nil, err
	}

	var list []string
	for i := range versionList {
		list = append(list, versionList[i].(map[string]interface{})["version"].(string))
	}

	return list, nil
}
