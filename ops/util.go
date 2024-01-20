package ops

import (
	"fmt"
	"regexp"
	"time"
)

var (
	ch       chan string
	vb       bool
	logging  bool
	postLogs []string

	versionRegex = regexp.MustCompile(`1\.\d+-pre\d+|1\.\d+\.\d+[-_]pre\d+|1\.\d+-rc\d+|1\.\d+\.\d+-rc\d+|1\.\d+.\d+|1\.\d+`)
)

func StartLog(verbose bool) {
	logging = true
	vb = verbose
	if !vb {
		ch = AsyncSpinner()
	}
}

func EndLog(things ...any) {
	if !vb {
		str := fmt.Sprintln(things...)
		ch <- str[:len(str)-1]
		close(ch)
		time.Sleep(time.Millisecond)
	} else {
		fmt.Println(things...)
	}

	for i := range postLogs {
		fmt.Println(postLogs[i])
	}
}

func log(things ...any) {
	if !logging {
		return
	}
	if vb {
		fmt.Println(things...)
	}
}

func slog(things ...any) {
	if !logging {
		return
	}
	if !vb {
		str := fmt.Sprintln(things...)
		ch <- str[:len(str)-1]
	} else {
		fmt.Println(things...)
	}
}

func post(s string) {
	if !logging {
		return
	}
	postLogs = append(postLogs, s)
}

func AsyncSpinner() chan string {
	ch := make(chan string)
	// because for some fucking reason syscall.Stdout is a different type on windows and idgaf to fix it
	//if !terminal.IsTerminal(syscall.Stdout) {
	//	close(ch)
	//}

	go func() {
		chars := []string{"|", "/", "-", "\\"}
		var title string
		i := 0
		for {
			if i+1 < len(chars) {
				i++
			} else {
				i = 0
			}
			fmt.Printf("\033[2K\n\033[1A%v %v", chars[i], title)

			select {
			case rx, open := <-ch:
				if open {
					title = rx
					continue
				}
				fmt.Printf("\033[2K\n\033[1A%v %v\n", "✓", title)
				return
			case <-time.After(150 * time.Millisecond):
			}
		}
	}()
	return ch
}
