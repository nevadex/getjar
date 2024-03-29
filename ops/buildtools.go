package ops

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	BuildSpigot      = 1
	BuildCraftBukkit = 2

	BuildToolsDir = ".getjar-buildtools"
	BuildToolsJar = "BuildTools.jar"
)

func RunBuildTools(verbose bool, btverbose bool, buildType int, version string, filename string, experimental bool) (string, error) {
	slog("initializing buildtools")
	err := getBuildTools()
	if err != nil {
		return "", err
	}
	err = os.Chdir(BuildToolsDir)
	if err != nil {
		return "", err
	}
	log("changed working directory to buildtools directory")

	cmd := exec.Command("java", "-jar", BuildToolsJar, "--compile", "NONE")
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	if err != nil {
		return "", err
	}
	log("buildtools ran for first time")

	if buildType == BuildSpigot {
		if !experimental {
			cmd = exec.Command("java", "-jar", BuildToolsJar, "--rev", version, "--compile", "SPIGOT")
		} else if experimental {
			cmd = exec.Command("java", "-jar", BuildToolsJar, "--rev", version, "--compile", "SPIGOT", "--experimental")
		}
	} else if buildType == BuildCraftBukkit {
		if !experimental {
			cmd = exec.Command("java", "-jar", BuildToolsJar, "--rev", version, "--compile", "CRAFTBUKKIT")
		} else if experimental {
			cmd = exec.Command("java", "-jar", BuildToolsJar, "--rev", version, "--compile", "CRAFTBUKKIT", "--experimental")
		}
	}

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	log("created pipes")

	slog("compiling with buildtools")
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	log("started process")

	go func() {
		if btverbose {
			_, _ = io.Copy(os.Stdout, cmdStdout)
		} else {
			_, _ = io.Copy(io.Discard, cmdStdout)
		}
	}()

	go func() {
		if verbose {
			_, _ = io.Copy(os.Stderr, cmdStderr)
		} else {
			_, _ = io.Copy(io.Discard, cmdStderr)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		if !verbose && !btverbose {
			bs, _ := os.ReadFile("BuildTools.log.txt")
			s := string(bs)
			lines := strings.Split(s, "\n")

			fmt.Println("buildtools ran into an error")
			fmt.Println("last five lines of BuildTools.log.txt")
			fmt.Println(lines[len(lines)-1])
			fmt.Println(lines[len(lines)-2])
			fmt.Println(lines[len(lines)-3])
			fmt.Println(lines[len(lines)-4])
			fmt.Println(lines[len(lines)-5])
		}
		return "", err
	}
	log("successfully compiled jar with buildtools")

	slog("renaming jar")
	outBytes, err := os.ReadFile("BuildTools.log.txt")
	if err != nil {
		return "", err
	}
	out := string(outBytes)
	var regex *regexp.Regexp
	if buildType == BuildSpigot {
		regex = regexp.MustCompile(`(spigot-1\.\d\d.\d\.jar)`)
	} else if buildType == BuildCraftBukkit {
		regex = regexp.MustCompile(`(craftbukkit-1\.\d\d.\d\.jar)`)
	}
	matches := regex.FindStringSubmatch(out)
	autoFileName := matches[len(matches)-1]
	version = regexp.MustCompile(`(1\.\d\d.\d)`).FindStringSubmatch(autoFileName)[0]
	log("extracted autogenerated filename")

	err = os.Chdir("..")
	if err != nil {
		return "", err
	}
	log("changed working directory to original directory")

	bytes, err := os.ReadFile(filepath.Join(BuildToolsDir, autoFileName))
	if err != nil {
		return "", err
	}
	err = os.WriteFile(filename, bytes, os.ModePerm)
	log("moved file and renamed")

	if printChecksum {
		post("no checksum provided")
	}

	return version, nil
}

func getBuildTools() error {
	if _, er := os.Stat(BuildToolsDir); os.IsNotExist(er) {
		log("buildtools missing, getting buildtools")

		err := os.MkdirAll(BuildToolsDir, os.ModePerm)
		if err != nil {
			return err
		}
		log("created directory for buildtools:", BuildToolsDir)

		resp, err := http.Get("https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar")
		if err != nil {
			return err
		}
		log("downloaded BuildTools.jar")
		defer func() { _ = resp.Body.Close() }()

		file, err := os.Create(filepath.Join(BuildToolsDir, BuildToolsJar))
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
		log("finished getting buildtools")

		//if !(os.PathSeparator == '\\' && os.PathListSeparator == ';') {
		//	cmd := exec.Command("git", "config", "--global", "--unset", "core.autocrlf")
		//	err = cmd.CombinedOutput()
		//	if err != nil {
		//		return err
		//	}
		//	log("set git config according to buildtools docs (non-windows only)")
		//}

		return nil
	}
	return nil
}

func GetVersionListBuildTools() ([]string, error) {
	resp, err := http.Get("https://hub.spigotmc.org/versions/")
	if err != nil {
		return nil, err
	}
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	upper := versionRegex.FindAllStringSubmatch(string(html), -1)

	noDupes := make(map[string]bool)
	for i := range upper {
		noDupes[upper[i][0]] = true
	}

	var lower []string
	for k := range noDupes {
		lower = append(lower, k)
	}

	return lower, nil
}
