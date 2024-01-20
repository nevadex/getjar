package ops

import (
	"errors"
	"github.com/anaskhan96/soup"
	"io"
	"net/http"
	"strings"
)

func DownloadMinecraftForge(version string, forgeVersion string, experimental bool) ([]byte, string, error) {
	log("attempting to download", version, "using minecraft forge html parsing")

	slog("downloading minecraft version list")
	resp, err := soup.Get("https://files.minecraftforge.net/net/minecraftforge/forge/index.html")
	if err != nil {
		return nil, "", err
	}
	log("retrieved minecraft version list")

	slog("getting version lists")

	log("parsing minecraft version list")
	doc := soup.HTMLParse(resp)
	var versionIndexes []string
	elems := doc.FindAll("a")
	for _, link := range elems {
		href := link.Attrs()["href"]
		if strings.Contains(href, "index") && versionRegex.MatchString(href) {
			versionIndexes = append(versionIndexes, "https://files.minecraftforge.net/net/minecraftforge/forge/"+href)
		}
	}
	versionIndexes = append(versionIndexes, "https://files.minecraftforge.net/net/minecraftforge/forge/index_"+doc.Find("li", "class", "elem-active").Text()+".html")
	log("parsed all minecraft versions")

	var index string
	if version == "latest" {
		if experimental {
			index = versionIndexes[0]
		} else if !experimental {
			index = versionIndexes[len(versionIndexes)-1]
		}
	} else {
		found := false
		for i := range versionIndexes {
			if strings.Contains(versionIndexes[i], version) {
				found = true
				index = "https://files.minecraftforge.net/net/minecraftforge/forge/index_" + version + ".html"
				break
			}
		}
		if !found {
			return nil, "", errors.New("could not locate version in version list")
		}
	}
	version = versionRegex.FindStringSubmatch(index)[0]
	log("using version", version)

	log("downloading forge version list")
	resp, err = soup.Get(index)
	if err != nil {
		return nil, "", err
	}
	log("retrieved forge version list")

	log("parsing forge version list")
	doc = soup.HTMLParse(resp)
	indexOfVersionListIfLatestAndExperimental := 0
	if forgeVersion == "latest" && !experimental {
		elems = doc.FindAll("td", "class", "download-version")
		for i := range elems {
			if len(elems[i].Children()) == 3 {
				if strings.Contains(elems[i].Children()[1].Attrs()["class"], "recommended") {
					indexOfVersionListIfLatestAndExperimental = i
					break
				}
			}
		}
	}

	elems = doc.FindAll("div", "class", "info-tooltip")
	var jarSha1, jarUrl string
	//goland:noinspection GoUnreachableCode
	for i := 1; i < len(elems)-2; i += 3 {
		e := elems[i]
		sha1 := strings.Trim(e.Children()[6].HTML(), " \n")
		url := e.Children()[8].Attrs()["href"]

		if forgeVersion == "latest" {
			if experimental {
				jarUrl = url
				jarSha1 = sha1
				break
			} else {
				if i/3 == indexOfVersionListIfLatestAndExperimental {
					jarUrl = url
					jarSha1 = sha1
					break
				}
			}
		} else {
			if strings.Contains(url, forgeVersion) {
				jarUrl = url
				jarSha1 = sha1
			}
			break
		}
	}
	log("parsed all forge versions")

	slog("downloading server installer")
	res, err := http.Get(jarUrl)
	if err != nil {
		return nil, "", err
	}
	log("downloaded server installer")
	defer func() { _ = res.Body.Close() }()

	jar, err := io.ReadAll(res.Body)

	//log("jarfile size:", len(jar))
	//log("minecraft version number:", version)
	if printChecksum {
		post("sha1 checksum:", jarSha1)
	}

	post("the jar is an installer, not a server jar!")

	return jar, version, nil
}

func GetVersionListMinecraftForge() ([]string, error) {
	resp, err := soup.Get("https://files.minecraftforge.net/net/minecraftforge/forge/index.html")
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(resp)
	var versions []string
	elems := doc.FindAll("a")
	for _, link := range elems {
		href := link.Attrs()["href"]
		if strings.Contains(href, "index") && versionRegex.MatchString(href) {
			versions = append(versions, versionRegex.FindStringSubmatch(href)[0])
		}
	}
	versions = append([]string{doc.Find("li", "class", "elem-active").Text()}, versions...)

	return versions, nil
}
