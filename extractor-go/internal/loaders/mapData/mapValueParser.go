package mapData

import (
	"strings"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
)

func ParseMapValueTable(filePath string) (map[string]string, error) {
	table := make(map[string]string)

	lines, err := decoders.DecodeTokenTextTable(filePath, 0)
	if err != nil {
		return table, err
	}

	for i := 0; i < len(lines); i += 2 {
		mapId := lines[i]
		value := lines[i+1]

		// For some reason, map files has some maps with ";" prefix, which I believe are ignored
		if strings.HasPrefix(mapId, ";") {
			continue
		}

		table[mapId] = value
	}

	return table, nil
}
