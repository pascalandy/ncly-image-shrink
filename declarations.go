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
	if len(chunks) == 0 || chunks[0].kind != "IHDR" || len(chunks[0].data) != 13 {
		return errImage
	}
	depth, model := chunks[0].data[8], chunks[0].data[9]
	for _, c := range chunks {
		if c.kind == "IDAT" {
			afterPixels = true
		}
		switch c.kind {
		case "gAMA", "cHRM", "sRGB", "cICP", "iCCP", "sBIT", "eXIf":
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
		case "sBIT":
			count := map[byte]int{0: 1, 2: 3, 3: 3, 4: 2, 6: 4}[model]
			if len(c.data) != count || count == 0 {
				return errImage
			}
			limit := depth
			if model == 3 {
				limit = 8
			}
			for _, n := range c.data {
				if n == 0 || n > limit {
					return errImage
				}
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

func jpegDeclarations(chunks []chunk) (string, error) {
	components := byte(0)
	var adobe []byte
	for _, c := range chunks {
		switch c.kind {
		case "\xc0", "\xc1", "\xc2":
			if len(c.data) < 6 || components != 0 {
				return "", errImage
			}
			components = c.data[5]
		case "\xee":
			if bytes.HasPrefix(c.data, []byte("Adobe")) {
				if adobe != nil || len(c.data) != 12 || c.data[11] > 2 {
					return "", errImage
				}
				adobe = c.data
			}
		case "\xe0":
			if bytes.HasPrefix(c.data, []byte("JFIF\x00")) {
				if len(c.data) < 14 || c.data[5] != 1 || c.data[6] > 2 || c.data[7] > 2 || len(c.data) != 14+3*int(c.data[12])*int(c.data[13]) {
					return "", errImage
				}
			}
		}
	}
	model := map[byte]string{1: "GRAY", 3: "RGB ", 4: "CMYK"}[components]
	if model == "" {
		return "", errImage
	}
	if adobe != nil {
		t := adobe[11]
		if components == 1 && t != 0 || components == 3 && t > 1 || components == 4 && t != 0 && t != 2 {
			return "", errImage
		}
	}
	return model, nil
}

func profileMatches(profile []byte, model string) bool {
	if len(profile) < 132 || len(profile) > 4<<20 || uint64(binary.BigEndian.Uint32(profile)) != uint64(len(profile)) || string(profile[36:40]) != "acsp" || string(profile[16:20]) != model {
		return false
	}
	count := uint64(binary.BigEndian.Uint32(profile[128:132]))
	tableEnd := 132 + 12*count
	if tableEnd > uint64(len(profile)) {
		return false
	}
	for i := uint64(0); i < count; i++ {
		at := 132 + 12*i
		offset, size := uint64(binary.BigEndian.Uint32(profile[at+4:])), uint64(binary.BigEndian.Uint32(profile[at+8:]))
		if offset < tableEnd || offset%4 != 0 || size < 8 || offset+size > uint64(len(profile)) {
			return false
		}
	}
	return true
}
