package decoders

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/csv"
	"os"
)

type ROCsvEntry struct {
	Id         int
	KoreanText string
	EngText    string
	PtBrText   string
	EsText     string
}

func padBase64String(base64Str string) string {
	for len(base64Str)%4 != 0 {
		base64Str += "="
	}
	return base64Str
}

func decodeBase64ToUInt32(base64Str string) (uint32, error) {
	data, err := base64.StdEncoding.DecodeString(padBase64String(base64Str))
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(data), nil
}

func decodeBase64ToStr(base64Str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(padBase64String(base64Str))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func DecodeLangCsv(filePath string) ([]ROCsvEntry, error) {
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

	var entries []ROCsvEntry
	for _, line := range lines {
		id, err := decodeBase64ToUInt32(line[0])
		if err != nil {
			return nil, err
		}
		koreanText, err := decodeBase64ToStr(line[1])
		if err != nil {
			return nil, err
		}
		engText, err := decodeBase64ToStr(line[2])
		if err != nil {
			return nil, err
		}
		ptBrText, err := decodeBase64ToStr(line[7])
		if err != nil {
			return nil, err
		}
		esText, err := decodeBase64ToStr(line[9])
		if err != nil {
			return nil, err
		}

		entry := ROCsvEntry{
			Id:         int(id),
			KoreanText: koreanText,
			EngText:    engText,
			PtBrText:   ptBrText,
			EsText:     esText,
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
