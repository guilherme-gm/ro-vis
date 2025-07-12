package grf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strings"

	binUtils "github.com/guilherme-gm/ro-vis/extractor/internal/binUtils"

	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/patchfile/des"
)

type EntryType uint8

const (
	EntryType_Directory EntryType = 0 // entry is a directory
	EntryType_File      EntryType = 1 // entry is a file
)

var EntryTypeName = map[EntryType]string{
	EntryType_Directory: "Directory",
	EntryType_File:      "File",
}

func (et EntryType) String() string {
	return EntryTypeName[et]
}

type EncryptionType uint8

const (
	Encryption_None   EncryptionType = 0 // no encryption
	Encryption_Mixed  EncryptionType = 1 // encryption mode 0 (header DES + periodic DES/shuffle)
	Encryption_Header EncryptionType = 2 // encryption mode 1 (header DES only)
)

var EncryptionTypeName = map[EncryptionType]string{
	Encryption_None:   "None",
	Encryption_Mixed:  "Mixed",
	Encryption_Header: "Header",
}

func (et EncryptionType) String() string {
	return EncryptionTypeName[et]
}

type FileTableEntry struct {
	CompressedSize          int   ///< compressed size (srclen)
	CompressedSizeAlignment int   ///< compressed size alignment (srclenaligned)
	OriginalSize            int   ///< original size (declen)
	ExactOffset             int64 ///< position of entry in grf (srcpos)
	Flags                   EntryType
	FileName                string ///< file name
	Encryption              EncryptionType
}

func (ft *FileTableEntry) readFromV1(r io.Reader) error {
	size := binUtils.ReadInt32(r)

	binUtils.ReadBytes(r, 2)

	encryptedName := binUtils.ReadBytes(r, int(size)-6)

	binUtils.ReadBytes(r, 4)

	ft.FileName = strings.Trim(grfio_decode_filename(encryptedName), "\x00")

	// +ofs2

	lenEnd := binUtils.ReadInt32(r)
	srcLenAligned := binUtils.ReadInt32(r) - 37579
	decLen := binUtils.ReadInt32(r)
	fileType := binUtils.ReadByte(r)
	srcPos := binUtils.ReadInt32(r) + 46

	srcLen := lenEnd - decLen - 715

	ft.CompressedSize = int(srcLen)
	ft.CompressedSizeAlignment = int(srcLenAligned)
	ft.OriginalSize = int(decLen)
	ft.ExactOffset = int64(srcPos)
	ft.Flags = EntryType(fileType)
	ft.Encryption = Encryption_None
	if ft.Flags == EntryType_File && len(ft.FileName) > 4 {
		extension := ft.FileName[len(ft.FileName)-4:]
		if extension == ".gnd" || extension == ".gat" || extension == ".act" || extension == ".str" {
			ft.Encryption = Encryption_Header
		} else {
			ft.Encryption = Encryption_Mixed
		}
	}

	if strings.ContainsRune(ft.FileName, '?') {
		return fmt.Errorf("invalid filename (encoding issue?): %s", ft.FileName)
	}

	return nil
}

func (ft *FileTableEntry) readFromV2(r io.Reader) error {
	ft.FileName = binUtils.ReadString(r)
	ft.CompressedSize = int(binUtils.ReadInt32(r))
	ft.CompressedSizeAlignment = int(binUtils.ReadInt32(r))
	ft.OriginalSize = int(binUtils.ReadInt32(r))
	ft.Flags = EntryType(binUtils.ReadByte(r))
	ft.ExactOffset = int64(binUtils.ReadInt32(r) + headerSize)

	return nil
}

/**
 * (De)-obfuscates data.
 *
 * Substitutes some specific values for others, leaves rest intact.
 * NOTE: Operation is symmetric (calling it twice gives back the original input).
 */
func grf_substitution(in uint8) uint8 {
	var out uint8

	grfSubstitutionMap := map[uint8]uint8{
		0x00: 0x2B,
		0x2B: 0x00,
		0x6C: 0x80,
		0x01: 0x68,
		0x68: 0x01,
		0x48: 0x77,
		0x60: 0xFF,
		0x77: 0x48,
		0xB9: 0xC0,
		0xC0: 0xB9,
		0xFE: 0xEB,
		0xEB: 0xFE,
		0x80: 0x6C,
		0xFF: 0x60,
	}
	out, ok := grfSubstitutionMap[in]
	if !ok {
		out = in
	}

	return out
}

func grf_shuffle_dec(buf *[]byte, start int, end int) {
	// copy start:end from buf to src, a new variable
	var src [8]byte
	copy(src[:], (*buf)[start:end])

	(*buf)[start] = src[3]
	(*buf)[start+1] = src[4]
	(*buf)[start+2] = src[6]
	(*buf)[start+3] = src[0]
	(*buf)[start+4] = src[1]
	(*buf)[start+5] = src[2]
	(*buf)[start+6] = src[5]
	(*buf)[start+7] = grf_substitution(src[7])
}

/**
 * Decodes fully encrypted grf data
 *
 * @param[in,out] buf   Data to decode (in-place).
 * @param[in]     len   Length of the data.
 * @param[in]     cycle The current decoding cycle.
 */
func grf_decode_full(buf *[]byte, len int, cycle int) {
	nblocks := len / 8
	var dcycle, scycle int
	var i, j int

	// first 20 blocks are all des-encrypted
	for i := 0; i < 20 && i < nblocks; i++ {
		if err := des.DesDecryptBlock_Super(buf, i*8, (i+1)*8); err != nil {
			panic(err)
		}
	}

	// after that only one of every 'dcycle' blocks is des-encrypted
	dcycle = cycle

	// and one of every 'scycle' plaintext blocks is shuffled (starting from the 0th but skipping the 0th)
	scycle = 7

	// so decrypt/de-shuffle periodically
	j = -1 // 0, adjusted to fit the ++j step
	for i = 20; i < nblocks; i++ {
		if i%dcycle == 0 {
			// decrypt block
			if err := des.DesDecryptBlock_Super(buf, i*8, (i+1)*8); err != nil {
				panic(err)
			}
			continue
		}

		j++
		if j%scycle == 0 && j != 0 {
			// de-shuffle block
			grf_shuffle_dec(buf, i*8, (i+1)*8)
			continue
		}

		// plaintext, do nothing.
	}
}

/**
 * Decodes header-encrypted grf data.
 *
 * @param[in,out] buf   Data to decode (in-place).
 * @param[in]     len   Length of the data.
 */
func grf_decode_header(buffer *[]byte) {
	nblocks := len(*buffer) / 8

	// first 20 blocks are all des-encrypted
	for i := 0; i < 20 && i < nblocks; i++ {
		des.DesDecryptBlock_Super(buffer, i*8, (i+1)*8)
	}

	// the rest is plaintext, done.
}

func grf_decode(buffer *[]byte, ft *FileTableEntry) {
	switch ft.Encryption {
	case Encryption_None:
		return
	case Encryption_Header:
		// header encrypted
		grf_decode_header(buffer)
	case Encryption_Mixed:
		// fully encrypted

		// compute number of digits of the entry length
		digits := 1
		for i := 10; i <= ft.CompressedSize; i *= 10 {
			digits++
		}

		// choose size of gap between two encrypted blocks
		// digits:  0  1  2  3  4  5  6  7  8  9 ...
		//  cycle:  1  1  1  4  5 14 15 22 23 24 ...
		cycle := 1
		if digits < 3 {
			cycle = 1
		} else if digits < 5 {
			cycle = digits + 1
		} else if digits < 7 {
			cycle = digits + 9
		} else {
			cycle = digits + 15
		}

		grf_decode_full(buffer, ft.CompressedSizeAlignment, cycle)

		return
	default:
		panic("unknown encryption type")
	}
}

func (ft *FileTableEntry) Extract(path string) error {
	if ft.Flags == EntryType_Directory {
		panic("directory not supported yet")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Seek(ft.ExactOffset, 0); err != nil {
		return err
	}

	buffer := make([]byte, ft.CompressedSizeAlignment)
	if _, err := io.ReadFull(file, buffer); err != nil {
		return err
	}

	grf_decode(&buffer, ft)

	data, err := zlib.NewReader(bytes.NewReader(buffer))
	if err != nil {
		return err
	}

	defer data.Close()

	// read all data
	dataBytes, err := io.ReadAll(data)
	if err != nil {
		return err
	}

	outFile, err := os.Create("out_file.dat")
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := outFile.Write(dataBytes); err != nil {
		return err
	}

	return nil
}
