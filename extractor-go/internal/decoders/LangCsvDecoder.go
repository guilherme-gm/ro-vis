package decoders

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/csv"
	"errors"
	"os"
)

type LangCsvEntry struct {
	Id         uint64
	KoreanText string
	EnText     string
	PtBrText   string
	EsText     string
}

func padBase64String(base64Str string) string {
	for len(base64Str)%4 != 0 {
		base64Str += "="
	}
	return base64Str
}

func decodeBase64ToUInt64(base64Str string) (uint64, error) {
	data, err := base64.StdEncoding.DecodeString(padBase64String(base64Str))
	if err != nil {
		return 0, err
	}

	if len(data) < 8 {
		// Padding with 0s - little endian should receive the 0 at start
		data = append(make([]byte, 8-len(data)), data...)
	}
	if len(data) > 8 {
		return 0, errors.New("data is too long")
	}

	return binary.LittleEndian.Uint64(data), nil
}

func decodeBase64ToStr(base64Str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(padBase64String(base64Str))
	if err != nil {
		return "", err
	}

	return string(data), nil
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
		id, err := decodeBase64ToUInt64(line[0])
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

		entry := LangCsvEntry{
			Id:         id,
			KoreanText: koreanText,
			EnText:     engText,
			PtBrText:   ptBrText,
			EsText:     esText,
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
