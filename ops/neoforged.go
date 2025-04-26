package ops

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func DownloadNeoforged(version string, neoVersion string, experimental bool) ([]byte, string, error) {
	log("attempting to download", neoVersion, "neoforged for mc", version, "using neoforged maven")

	slog("downloading version list")
	versionList, err := GetVersionListNeoforged()

	if version == "latest" && neoVersion == "latest" {
		for i := len(versionList) - 1; i >= 0; i-- {
			if !experimental && strings.Contains(versionList[i], "beta") {
				continue
			}

			neoVersion = versionList[i]
			neoSplit := strings.Split(neoVersion, ".")
			version = "1." + neoSplit[0] + "." + neoSplit[1]
			break
		}
	} else if version != "latest" && neoVersion == "latest" {
		found := false
		foundWithRightMinecraftVersion := false
		for i := len(versionList) - 1; i >= 0; i-- {
			if strings.Contains(versionList[i], version[2:]) {
				foundWithRightMinecraftVersion = true
				if !experimental && strings.Contains(versionList[i], "beta") {
					continue
				}
				neoVersion = versionList[i]
				found = true
				break
			}
		}
		if !found {
			if !foundWithRightMinecraftVersion {
				return nil, "", errors.New("could not find neo version for minecraft version " + version)
			} else {
				return nil, "", errors.New("no 'default' build found")
			}
		}
	} else {
		if version != "latest" {
			neoSplit := strings.Split(neoVersion, ".")
			if version != "1."+neoSplit[0]+"."+neoSplit[1] {
				return nil, "", errors.New("minecraft and neo versions are mismatched")
			}
		}

		if !slices.Contains(versionList, neoVersion) {
			return nil, "", errors.New("neo version not found in version list")
		}
	}

	//if neoVersion == "latest" {
	//	found := false
	//	for i := len(versionList) - 1; i >= 0; i-- {
	//		if version != "latest" && !strings.Contains(versionList[i], version[2:]) {
	//			continue
	//		}
	//
	//		if !experimental && strings.Contains(versionList[i], "beta") {
	//			continue
	//		}
	//
	//		version = versionList[i]
	//		found = true
	//		break
	//	}
	//	if !found {
	//
	//	}
	//} else {
	//	if version != "latest" {
	//		neoSplit := strings.Split(neoVersion, ".")
	//		if version != "1."+neoSplit[0]+"."+neoSplit[1] {
	//			return nil, "", errors.New("minecraft and neo versions are mismatched")
	//		}
	//	}
	//
	//	if !slices.Contains(versionList, neoVersion) {
	//		return nil, "", errors.New("neo version not found in version list")
	//	}
	//}
	log("using neo version", neoVersion)

	jarUrl := "https://maven.neoforged.net/releases/net/neoforged/neoforge/" + neoVersion + "/neoforge-" + neoVersion + "-installer.jar"

	slog("downloading server installer")
	res, err := http.Get(jarUrl)
	if err != nil {
		return nil, "", err
	}
	log("downloaded server installer")

	jar, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}
	err = res.Body.Close()
	if err != nil {
		return nil, "", err
	}

	if printChecksum {
		slog("downloading checksum")
		res, err = http.Get(jarUrl + ".sha1")
		if err != nil {
			return nil, "", err
		}
		jarSha1Bytes, er := io.ReadAll(res.Body)
		if er != nil {
			return nil, "", er
		}
		post("sha1 checksum:", string(jarSha1Bytes))
	}

	post("the jar is an installer, not a server jar!")

	return jar, version, nil
}

func GetVersionListNeoforged() ([]string, error) {
	resp, err := http.Get("https://maven.neoforged.net/releases/net/neoforged/neoforge/maven-metadata.xml")
	if err != nil {
		return nil, err
	}
	log("retrieved neoforged version list xml")
	defer func() { _ = resp.Body.Close() }()

	decoder := xml.NewDecoder(resp.Body)
	var list []string
	for {
		t, er := decoder.Token()
		if er != nil {
			break
		}

		switch token := t.(type) {
		case xml.StartElement:
			if token.Name.Local == "version" {
				var v string
				if e := decoder.DecodeElement(&v, &token); e != nil {
					return nil, e
				}
				list = append(list, v)
			}
		}
	}
	log("decoded version list xml")

	sort.SliceStable(list, func(i, j int) bool {
		vi := list[i]
		vj := list[j]

		if strings.Contains(vi, "-") {
			vi = strings.Split(vi, "-")[0]
		}
		if strings.Contains(vj, "-") {
			vj = strings.Split(vj, "-")[0]
		}

		iSplit := strings.Split(vi, ".")
		if len(iSplit) != 3 {
			return true
		}
		jSplit := strings.Split(vj, ".")
		if len(jSplit) != 3 {
			return false
		}

		iMajor, er := strconv.Atoi(iSplit[0])
		if er != nil {
			return false
		}
		jMajor, er := strconv.Atoi(jSplit[0])
		if er != nil {
			return true
		}
		if iMajor != jMajor {
			return iMajor < jMajor
		}

		iMinor, er := strconv.Atoi(iSplit[1])
		if er != nil {
			return false
		}
		jMinor, er := strconv.Atoi(jSplit[1])
		if er != nil {
			return true
		}
		if iMinor != jMinor {
			return iMinor < jMinor
		}

		iPatch, er := strconv.Atoi(iSplit[2])
		if er != nil {
			return false
		}
		jPatch, er := strconv.Atoi(jSplit[2])
		if er != nil {
			return true
		}
		return iPatch < jPatch
	})
	log("sorted versions")

	return list, nil
}
