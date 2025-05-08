package decoders

import (
	"bytes"
	"io"

	"golang.org/x/net/html/charset"
)

func convertToUTF8(str string) string {
	strBytes := []byte(str)
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel("euc-kr", byteReader)
	strBytes, _ = io.ReadAll(reader)
	return string(strBytes)
}
