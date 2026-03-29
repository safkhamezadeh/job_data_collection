package keywordextractor

import (
	"fmt"
	"strings"
)

type KeyWordFormat struct {
	Jobs     []string
	Keywords []string
}

func (k KeyWordFormat) String() string {
	return fmt.Sprintf(
		"titles: %s\nkeywords: %s",
		strings.Join(k.Jobs, ", "),
		strings.Join(k.Keywords, ", "),
	)
}

func (k KeyWordFormat) ToString(sep string) string {
	return fmt.Sprintf(
		"%s %s",
		strings.Join(k.Jobs, sep),
		strings.Join(k.Keywords, sep),
	)
}

func StoKeyWordFormat(input string) KeyWordFormat {
	lines := strings.Split(input, "\n")

	var result KeyWordFormat

	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2) // split into key and values
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		values := strings.Split(parts[1], ",")
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}

		switch key {
		case "titles":
			result.Jobs = values
		case "keywords":
			result.Keywords = values
		}
	}
	return result
}

func (f KeyWordFormat) IsValid() bool {
	if len(f.Jobs) == 5 && len(f.Keywords) == 5 {
		return true
	}
	return false
}
