package ops

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"time"
)

var (
	ch chan string
	vb bool
)

func StartLog(verbose bool) {
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
}

func log(things ...any) {
	if vb {
		fmt.Println(things...)
	}
}

func slog(things ...any) {
	if !vb {
		str := fmt.Sprintln(things...)
		ch <- str[:len(str)-1]
	} else {
		fmt.Println(things...)
	}
}

func AsyncSpinner() chan string {
	ch := make(chan string)
	if !terminal.IsTerminal(syscall.Stdout) {
		close(ch)
	}

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
