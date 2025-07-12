package rgz

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	binutils "github.com/guilherme-gm/ro-vis/extractor/internal/binUtils"
)

type RgzFile struct {
	Path    string
	Entries []RgzEntry
}

type EntryType uint8

const (
	EntryType_File      EntryType = 'f'
	EntryType_Directory EntryType = 'd'
)

var EntryTypeName = map[EntryType]string{
	EntryType_File:      "File",
	EntryType_Directory: "Directory",
}

func (et EntryType) String() string {
	return EntryTypeName[et]
}

type RgzEntry struct {
	Name      string
	EntryType EntryType
	Data      []byte
}

func readDirectoryEntry(reader io.Reader) (*RgzEntry, error) {
	len := binutils.ReadByte(reader)
	fileName := string(binutils.ReadBytes(reader, int(len)))
	fileName = strings.Split(fileName, "\x00")[0]

	return &RgzEntry{
		Name:      fileName,
		EntryType: EntryType_Directory,
	}, nil
}

func readFileEntry(reader io.Reader) (*RgzEntry, error) {
	len := binutils.ReadByte(reader)
	fileName := string(binutils.ReadBytes(reader, int(len)))
	fileName = strings.Split(fileName, "\x00")[0]

	dataLen := binutils.ReadInt32(reader)
	data := binutils.ReadBytes(reader, int(dataLen))

	return &RgzEntry{
		Name:      fileName,
		EntryType: EntryType_File,
		Data:      data,
	}, nil
}

func decompress(file io.Reader) ([]byte, error) {
	decompressor, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer decompressor.Close()

	decompressedData, err := io.ReadAll(decompressor)
	if err != nil {
		return nil, err
	}

	return decompressedData, nil
}

func readEntries(RgzFile *RgzFile, reader io.Reader) error {
	isEndReached := false
	for !isEndReached {
		entryType := rune(binutils.ReadByte(reader))
		switch entryType {
		case 'f':
			entry, err := readFileEntry(reader)
			if err != nil {
				return err
			}
			RgzFile.Entries = append(RgzFile.Entries, *entry)

		case 'd':
			entry, err := readDirectoryEntry(reader)
			if err != nil {
				return err
			}
			RgzFile.Entries = append(RgzFile.Entries, *entry)

		case 'e':
			isEndReached = true

		default:
			return fmt.Errorf("invalid entry type: %c", entryType)
		}
	}

	return nil
}

func Open(filename string) (*RgzFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decompressedData, err := decompress(file)
	if err != nil {
		return nil, err
	}

	RgzFile := &RgzFile{
		Path:    filename,
		Entries: []RgzEntry{},
	}

	if err := readEntries(RgzFile, bytes.NewReader(decompressedData)); err != nil {
		return nil, err
	}

	return RgzFile, nil
}
