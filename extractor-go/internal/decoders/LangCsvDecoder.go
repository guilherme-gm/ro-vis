package decoders

import (
	"encoding/csv"
	"os"
	"strconv"

	base64Utils "github.com/guilherme-gm/ro-vis/extractor/internal/utils/base64Utils"
)

type LangCsvEntry struct {
	Id         string
	KoreanText string
	EnText     string
	PtBrText   string
	EsText     string
}

func DecodeLangCsv(filePath string) ([]LangCsvEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ','
	csvReader.Comment = '#'
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true

	lines, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var entries []LangCsvEntry
	for _, line := range lines {
		id, err := base64Utils.DecodeBase64ToUInt64(line[0])
		if err != nil {
			return nil, err
		}
		koreanText, err := base64Utils.DecodeBase64ToStr(line[1])
		if err != nil {
			return nil, err
		}
		engText, err := base64Utils.DecodeBase64ToStr(line[2])
		if err != nil {
			return nil, err
		}
		ptBrText, err := base64Utils.DecodeBase64ToStr(line[7])
		if err != nil {
			return nil, err
		}
		esText, err := base64Utils.DecodeBase64ToStr(line[9])
		if err != nil {
			return nil, err
		}

		entry := LangCsvEntry{
			Id:         strconv.FormatUint(id, 10),
			KoreanText: koreanText,
			EnText:     engText,
			PtBrText:   ptBrText,
			EsText:     esText,
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
