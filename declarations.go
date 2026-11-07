package image

import (
	"bytes"
	"encoding/binary"
)

// PNG colour declarations have fixed layouts. Retain their validated bytes,
// including fallback declarations, without accepting arbitrary opaque text.
func pngDeclarations(chunks []chunk) error {
	seen := map[string][]byte{}
	afterPixels := false
	for _, c := range chunks {
		if c.kind == "IDAT" {
			afterPixels = true
		}
		switch c.kind {
		case "gAMA", "cHRM", "sRGB", "cICP", "iCCP":
			if _, exists := seen[c.kind]; exists || afterPixels {
				return errImage
			}
			seen[c.kind] = c.data
		default:
			continue
		}
		switch c.kind {
		case "gAMA":
			if len(c.data) != 4 || binary.BigEndian.Uint32(c.data) == 0 {
				return errImage
			}
		case "cHRM":
			if len(c.data) != 32 {
				return errImage
			}
			for at := 0; at < 32; at += 8 {
				x, y := binary.BigEndian.Uint32(c.data[at:]), binary.BigEndian.Uint32(c.data[at+4:])
				if x > 100000 || y == 0 || y > 100000 || uint64(x)+uint64(y) > 100000 {
					return errImage
				}
			}
		case "sRGB":
			if len(c.data) != 1 || c.data[0] > 3 {
				return errImage
			}
		case "cICP":
			// PNG stores RGB or Gray samples, hence identity matrix and full
			// range. Reserved H.273 primaries and transfer values are invalid.
			if len(c.data) != 4 || !bytes.Contains([]byte{1, 2, 4, 5, 6, 7, 8, 9, 10, 11, 12, 22}, c.data[:1]) || !bytes.Contains([]byte{1, 2, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}, c.data[1:2]) || c.data[2] != 0 || c.data[3] != 1 {
				return errImage
			}
		}
	}
	if _, srgb := seen["sRGB"]; srgb {
		if _, icc := seen["iCCP"]; icc {
			return errImage
		}
		if gamma, ok := seen["gAMA"]; ok && binary.BigEndian.Uint32(gamma) != 45455 {
			return errImage
		}
		if chroma, ok := seen["cHRM"]; ok {
			for i, value := range []uint32{31270, 32900, 64000, 33000, 30000, 60000, 15000, 6000} {
				if binary.BigEndian.Uint32(chroma[4*i:]) != value {
					return errImage
				}
			}
		}
	}
	return nil
}
