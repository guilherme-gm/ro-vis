package grf

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/patchfile/des"
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
func grfio_decode_filename(buf []byte) string {
	for i := 0; i < len(buf); i += 8 {
		NibbleSwap(&buf, i, i+8)
		des.DesDecrypt(&buf, i, i+8)
	}
	return string(buf)
}
