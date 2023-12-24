package ops

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func DownloadVanilla(verbose bool, version string) ([]byte, error) {
	log(verbose, "attempting to download vanilla", version)

	var versionManifest map[string]interface{}
	resp, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
	if err != nil {
		return nil, err
	}
	log(verbose, "downloaded version_manifest.json")
	defer func() { _ = resp.Body.Close() }()

	if err := json.NewDecoder(resp.Body).Decode(&versionManifest); err != nil {
		return nil, err
	}
	log(verbose, "decoded version_manifest.json")

	if version == "" || version == "latest" {
		version = versionManifest["latest"].(map[string]interface{})["release"].(string)
		log(verbose, "no version supplied, defaulting to latest:", version)
	}

	var versionUrl string
	var versionJson map[string]interface{}

	for _, v := range versionManifest["versions"].([]interface{}) {
		if v.(map[string]interface{})["id"].(string) == version {
			versionUrl = v.(map[string]interface{})["url"].(string)
			log(verbose, "located version in version_manifest.json")
			break
		}
	}
	if versionUrl == "" {
		return nil, errors.New("could not locate version in version_manifest.json")
	}

	resp, err = http.Get(versionUrl)
	if err != nil {
		return nil, err
	}
	log(verbose, "downloaded jar manifest")
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &versionJson)
	if err != nil {
		return nil, err
	}
	jarurl := versionJson["downloads"].(map[string]interface{})["server"].(map[string]interface{})["url"].(string)
	//jarsize := strconv.FormatFloat(versionJson["downloads"].(map[string]interface{})["server"].(map[string]interface{})["size"].(float64), 'f', -1, 64)
	jarsha1 := versionJson["downloads"].(map[string]interface{})["server"].(map[string]interface{})["sha1"].(string)
	//jarversiontype := versionJson["type"].(string)
	//jarjavaversion := strconv.FormatFloat(versionJson["javaVersion"].(map[string]interface{})["majorVersion"].(float64), 'f', -1, 64)
	log(verbose, "decoded jar manifest")

	jarresp, err := http.Get(jarurl)
	if err != nil {
		return nil, err
	}
	log(verbose, "downloaded server jar")
	defer func() { _ = resp.Body.Close() }()

	jar, err := io.ReadAll(jarresp.Body)

	//log(verbose, "jarfile size:", jarsize)
	//log(verbose, "minecraft version type:", jarversiontype)
	//log(verbose, "minecraft version number:", version)
	//log(verbose, "minimum java version:", jarjavaversion)
	log(verbose, "sha1 checksum:", jarsha1)

	return jar, nil
}
