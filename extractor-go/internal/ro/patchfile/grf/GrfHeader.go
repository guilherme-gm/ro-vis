package grf

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

const headerSize = 46

type GrfHeader struct {
	Magic           string
	Key             string
	FileTableOffset uint32
	Seed            uint32
	FileCount       uint32
	Version         uint32
	MajorVersion    uint8
	MinorVersion    uint8
	RealFilesCount  uint32
}

func NewGrfHeader(r io.Reader) (GrfHeader, error) {
	var header GrfHeader

	buf := make([]byte, headerSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return header, err
	}

	header.Magic = string(buf[:16])
	if strings.ToLower(header.Magic) != "master of magic\x00" {
		// Maybe alpha grf -- not supported for now
		return header, fmt.Errorf("invalid grf header: '%s'", strings.ToLower(header.Magic))
	}

	header.Key = string(buf[16:30])
	header.FileTableOffset = binary.LittleEndian.Uint32(buf[30:34])
	header.Seed = binary.LittleEndian.Uint32(buf[34:38])
	header.FileCount = binary.LittleEndian.Uint32(buf[38:42])
	header.Version = binary.LittleEndian.Uint32(buf[42:46])
	header.MajorVersion = uint8(header.Version >> 8)
	header.MinorVersion = uint8(header.Version & 0x000000FF)
	header.RealFilesCount = header.FileCount - header.Seed - 7

	return header, nil
}

func (h *GrfHeader) IsCompatibleWith(major uint8, minor uint8) bool {
	return (h.MajorVersion == major && h.MinorVersion >= minor) || (h.MajorVersion > major)
}
