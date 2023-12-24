package ops

import "fmt"

func log(verbose bool, things ...any) {
	if verbose {
		fmt.Println(things...)
	}
}
