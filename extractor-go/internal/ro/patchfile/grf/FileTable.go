package grf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"

	binUtils "github.com/guilherme-gm/ro-vis/extractor/internal/binUtils"
)

type FileTable struct {
	Files []FileTableEntry
}

func NewFileTable(r io.Reader, header GrfHeader) (FileTable, error) {
	var fileTable FileTable

	if header.IsCompatibleWith(2, 0) {
		fileTable.readFromV2(r, header.RealFilesCount) // for some reason, they don't give the right count
	} else if header.IsCompatibleWith(1, 0) {
		fileTable.readFromV1(r, header.RealFilesCount)
		// } else if header.IsCompatibleWith(0, 18) {
		// 	// fileTable.readFromV0(r)
	} else {
		return fileTable, fmt.Errorf("GRF Version %d.%d is not supported", header.MajorVersion, header.MinorVersion)
	}

	return fileTable, nil
}

func (ft *FileTable) readFromV1(r io.Reader, fileCount uint32) error {
	for range fileCount {
		var fileEntry FileTableEntry
		if err := fileEntry.readFromV1(r); err != nil {
			return err
		}
		ft.Files = append(ft.Files, fileEntry)
	}

	return nil
}

func (ft *FileTable) readFromV2(r io.Reader, fileCount uint32) error {
	TableSizeCompressed := binUtils.ReadUint32(r)
	TableSize := binUtils.ReadUint32(r)

	if TableSizeCompressed == 0 || TableSize == 0 {
		return fmt.Errorf("table size is 0 (compressed: %d, size: %d)", TableSizeCompressed, TableSize)
	}

	compressedData := binUtils.ReadBytes(r, int(TableSizeCompressed))
	data, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return err
	}

	defer data.Close()

	bufferPosition := 0
	for range fileCount {
		var fileEntry FileTableEntry
		if err = fileEntry.readFromV2(data); err != nil {
			return err
		}
		ft.Files = append(ft.Files, fileEntry)
		bufferPosition += fileEntry.CompressedSize
	}

	return nil
}
