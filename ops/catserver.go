package ops

import (
	"encoding/json"
	"io"
	"net/http"
)

func DownloadCatserver(version string) ([]byte, string, error) {
	log("attempting to download catserver", version)

	if version == "latest" {
		version = "1.12.2" // "universal" or most often used bc mods
		log("no version supplied, defaulting to", version)
	}

	slog("downloading build")
	var buildRaw map[string]interface{}
	resp, err := http.Get("https://jenkins.rbqcloud.cn:30011/job/CatServer-" + version + "/lastSuccessfulBuild/api/json/")
	if err != nil {
		return nil, "", err
	}
	log("retrieved build json")
	defer func() { _ = resp.Body.Close() }()

	if err = json.NewDecoder(resp.Body).Decode(&buildRaw); err != nil {
		return nil, "", err
	}
	artifacts := buildRaw["artifacts"].([]interface{})
	jarPath := artifacts[0].(map[string]interface{})["relativePath"].(string)
	log("decoded build json")

	slog("downloading jar")
	resp, err = http.Get("https://jenkins.rbqcloud.cn:30011/job/CatServer-" + version + "/lastSuccessfulBuild/artifact/" + jarPath)
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = resp.Body.Close() }()
	jar, err := io.ReadAll(resp.Body)
	log("downloaded jar")

	return jar, version, err
}

func GetVersionListCatserver() ([]string, error) {
	var jobsRaw map[string]interface{}
	resp, err := http.Get("https://jenkins.rbqcloud.cn:30011/api/json/")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if er := json.NewDecoder(resp.Body).Decode(&jobsRaw); er != nil {
		return nil, er
	}
	jobs := jobsRaw["jobs"].([]interface{})

	var list []string
	for i := range jobs {
		list = append(list, versionRegex.FindStringSubmatch(jobs[i].(map[string]interface{})["name"].(string))[0])
	}

	return list, nil
}
