package grf

import (
	"io"
	"os"
)

type GrfFile struct {
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
		Header:    header,
		FileTable: fileTable,
		grfReader: file,
	}, nil
}
