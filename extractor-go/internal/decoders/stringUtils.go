package decoders

import (
	"bytes"
	"io"

	"golang.org/x/net/html/charset"
)

type StringReencoder func(str string) string

func ConvertEucKrToUtf8(str string) string {
	strBytes := []byte(str)
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel("euc-kr", byteReader)
	strBytes, _ = io.ReadAll(reader)
	return string(strBytes)
}

func ConvertWin1252ToUtf8(str string) string {
	strBytes := []byte(str)
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel("windows-1252", byteReader)
	strBytes, _ = io.ReadAll(reader)
	return string(strBytes)
}

func ConvertNoop(str string) string {
	return str
}
