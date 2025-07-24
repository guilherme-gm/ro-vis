package grf

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type GrfFile struct {
	Path      string
	Header    GrfHeader
	FileTable FileTable
	grfReader io.Reader
}

// Open opens a GRF or GPF file and returns a GrfFile object.
func Open(path string) (*GrfFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	header, err := NewGrfHeader(file)
	if err != nil {
		return nil, err
	}

	file.Seek(int64(header.FileTableOffset+headerSize), 0)

	fileTable, err := NewFileTable(file, header)
	if err != nil {
		return nil, err
	}

	return &GrfFile{
		Path:      path,
		Header:    header,
		FileTable: fileTable,
		grfReader: file,
	}, nil
}

func (grfFile *GrfFile) Extract(filePath string, rootFolder string) error {
	toPath := path.Join(rootFolder, filePath)
	for _, file := range grfFile.FileTable.Files {
		if file.Flags == EntryType_File && strings.EqualFold(file.FileName, filePath) {
			return file.Extract(grfFile.Path, toPath)
		}
	}

	return fmt.Errorf("file %s not found", filePath)
}
