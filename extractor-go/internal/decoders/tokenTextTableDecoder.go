package decoders

import (
	"bufio"
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

func DecodeTokenTextTable(filePath string, readType int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	koreanReader := transform.NewReader(file, korean.EUCKR.NewDecoder())
	reader := bufio.NewReader(koreanReader)

	var lines []string
	nextVal := ""
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		if err != nil && err != io.EOF {
			return nil, err
		}

		// Note: in a real client, it actually consumes \n as part of the text and breaks on \r
		//       but chances are this will fail more often than not, so we just go with ignoring \r and using \n
		switch r {
		case '\r':
			continue

		case '\n':
			// Commented line, discard it
			if strings.HasPrefix(nextVal, "//") {
				nextVal = ""
				continue
			}

			if readType > 0 && nextVal != "" {
				lines = append(lines, nextVal)
			}

			nextVal = ""
			continue

		case '#':
			// Commented line, continue consuming as part of nextVal (will be discarded later)
			if strings.HasPrefix(nextVal, "//") {
				break
			}

			lines = append(lines, nextVal)
			if readType == 2 {
				lines = append(lines, "")
			}

			nextVal = ""
			continue
		}

		nextVal += string(r)
	}

	// nextVal after the last # is intentionally discarded
	// lines = append(lines, nextVal)

	return lines, nil
}
