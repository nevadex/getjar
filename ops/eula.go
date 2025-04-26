package ops

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func AcceptEula(filename string) error {
	allBytes, err := os.ReadFile(filename)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error opening file: %v", err)
	}

	if errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(filename, []byte("eula=true\n"), 0666)
	}

	all := string(allBytes)

	if len(strings.TrimSpace(all)) == 0 {
		return os.WriteFile(filename, []byte("eula=true\n"), 0666)
	}

	all = strings.ReplaceAll(all, "eula=false", "eula=true")

	return os.WriteFile(filename, []byte(all), 0666)
}
