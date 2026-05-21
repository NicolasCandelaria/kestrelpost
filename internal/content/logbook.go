package content

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed logbook/*.txt
var logbookFS embed.FS

func SeededLogbookEntry(night int) (string, error) {
	name := fmt.Sprintf("logbook/night%02d.txt", night)
	raw, err := logbookFS.ReadFile(name)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(raw)), nil
}
