// binutils provides utility functions for reading binary data from a reader.
package binutils

import (
	"encoding/binary"
	"io"
)

// Read a slice of bytes from a reader. Advances the reader to the next byte after the slice.
func ReadBytes(r io.Reader, size int) []byte {
	buf := make([]byte, size)
	if _, err := io.ReadFull(r, buf); err != nil {
		panic(err)
	}
	return buf
}

// Read a string from a reader until a null byte is encountered. Advances the reader to the next byte after the null byte.
func ReadString(r io.Reader) string {
	var buf []byte
	for {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			return string(buf)
		}
		if b[0] == 0 {
			return string(buf)
		}
		buf = append(buf, b[0])
	}
}

// Read a byte from a reader. Advances the reader to the next byte.
func ReadByte(r io.Reader) byte {
	var buf [1]byte
	if _, err := r.Read(buf[:]); err != nil {
		panic(err)
	}
	return buf[0]
}

// Read an int32 from a reader. Advances the reader to the next byte after the int32.
func ReadInt32(r io.Reader) int32 {
	var buf [4]byte
	if _, err := r.Read(buf[:]); err != nil {
		panic(err)
	}
	return int32(binary.LittleEndian.Uint32(buf[:]))
}

// Read a uint32 from a reader. Advances the reader to the next byte after the uint32.
func ReadUint32(r io.Reader) uint32 {
	var buf [4]byte
	if _, err := r.Read(buf[:]); err != nil {
		panic(err)
	}
	return binary.LittleEndian.Uint32(buf[:])
}
