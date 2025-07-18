package grf

import (
	"bytes"
	"fmt"
	"io"

	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/patchfile/des"
	"golang.org/x/text/encoding/charmap"
)

func NibbleSwap(src *[]byte, start, end int) {
	for idx := start; idx < end; idx++ {
		val := (*src)[idx]
		val = (val >> 4) | (val << 4)
		(*src)[idx] = val
	}
}

/**
 * Decodes encrypted filename from a version 01xx grf index.
 *
 * @param[in,out] buf The encrypted filename (decrypted in-place).
 * @param[in]     len The filename length.
 * @return A pointer to the decrypted filename.
 */
func grfio_decode_filename(buf []byte) (string, error) {
	for i := 0; i < len(buf); i += 8 {
		NibbleSwap(&buf, i, i+8)
		des.DesDecrypt(&buf, i, i+8)
	}

	buf2 := make([]byte, len(buf))
	copy(buf2, buf)

	for i := 0; i < len(buf2); i++ {
		if buf2[i] == 0 {
			buf2 = buf2[:i]
			break
		}
	}

	sr := bytes.NewReader(buf2)
	tr := charmap.Windows1252.NewDecoder().Reader(sr)
	buf2, err := io.ReadAll(tr)
	if err != nil {
		return "", fmt.Errorf("failed to decode filename: %w", err)
	}

	return string(buf2), nil
}
