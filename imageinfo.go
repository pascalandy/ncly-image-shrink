package image

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

var errImage = errors.New("invalid image container")
var errMultipart = errors.New("multiple images")

type metadata struct {
	icc         []byte
	orientation uint16
	colors      []chunk
	model       string
}

type chunk struct {
	kind string
	data []byte
}

// inspectImage validates container framing before a decoder or vips sees it.
// Unknown descriptive blocks are accepted on input and omitted from the copy.
func inspectImage(data []byte, format string) (metadata, error) {
	meta := metadata{orientation: 1}
	switch format {
	case "png":
		chunks, err := pngChunks(data)
		if err != nil {
			return meta, err
		}
		if err := pngDeclarations(chunks); err != nil {
			return meta, err
		}
		meta.model = "RGB "
		if len(data) > 25 && (data[25] == 0 || data[25] == 4) {
			meta.model = "GRAY"
		}
		for _, c := range chunks {
			switch c.kind {
			case "acTL", "fcTL", "fdAT":
				return meta, errMultipart
			case "eXIf":
				meta.orientation, err = orientation(c.data)
			case "iCCP":
				at := bytes.IndexByte(c.data, 0)
				if at < 1 || at+2 >= len(c.data) || c.data[at+1] != 0 || meta.icc != nil {
					return meta, errImage
				}
				reader, openErr := zlib.NewReader(bytes.NewReader(c.data[at+2:]))
				if openErr != nil {
					return meta, errImage
				}
				meta.icc, err = io.ReadAll(io.LimitReader(reader, (4<<20)+1))
				closeErr := reader.Close()
				if err != nil || closeErr != nil || len(meta.icc) > 4<<20 {
					return meta, errImage
				}
			case "gAMA", "cHRM", "sRGB", "cICP":
				meta.colors = append(meta.colors, c)
			}
			if err != nil {
				return meta, err
			}
		}
	case "jpeg":
		chunks, err := jpegChunks(data)
		if err != nil {
			return meta, err
		}
		meta.model, err = jpegDeclarations(chunks)
		if err != nil {
			return meta, err
		}
		profiles := map[byte][]byte{}
		total := byte(0)
		for _, c := range chunks {
			if c.kind == "\xe2" && bytes.HasPrefix(c.data, []byte("MPF\x00")) {
				return meta, errMultipart
			}
			if c.kind == "\xe1" && bytes.HasPrefix(c.data, []byte("Exif\x00\x00")) {
				meta.orientation, err = orientation(c.data[6:])
			}
			if err != nil {
				return meta, err
			}
			if c.kind == "\xe2" && bytes.HasPrefix(c.data, []byte("ICC_PROFILE\x00")) {
				if len(c.data) < 14 || c.data[12] < 1 || c.data[13] < c.data[12] || profiles[c.data[12]] != nil || total != 0 && total != c.data[13] {
					return meta, errImage
				}
				total = c.data[13]
				profiles[c.data[12]] = c.data[14:]
			}
		}
		for n := 1; n <= int(total); n++ {
			if profiles[byte(n)] == nil {
				return meta, errImage
			}
			meta.icc = append(meta.icc, profiles[byte(n)]...)
		}
	case "webp":
		meta.model = "RGB "
		chunks, err := webpChunks(data)
		if err != nil {
			return meta, err
		}
		for _, c := range chunks {
			switch c.kind {
			case "ANIM", "ANMF":
				return meta, errMultipart
			case "VP8X":
				if len(c.data) != 10 {
					return meta, errImage
				}
				if c.data[0]&2 != 0 {
					return meta, errMultipart
				}
			case "EXIF":
				exif := bytes.TrimPrefix(c.data, []byte("Exif\x00\x00"))
				meta.orientation, err = orientation(exif)
			case "ICCP":
				if meta.icc != nil {
					return meta, errImage
				}
				meta.icc = c.data
			}
			if err != nil {
				return meta, err
			}
		}
	default:
		return meta, errImage
	}
	if len(meta.icc) > 0 && !profileMatches(meta.icc, meta.model) {
		return meta, errImage
	}
	return meta, nil
}

func orientation(tiff []byte) (uint16, error) {
	if len(tiff) < 8 {
		return 0, errImage
	}
	var order binary.ByteOrder
	switch string(tiff[:2]) {
	case "II":
		order = binary.LittleEndian
	case "MM":
		order = binary.BigEndian
	default:
		return 0, errImage
	}
	if order.Uint16(tiff[2:4]) != 42 {
		return 0, errImage
	}
	offset := uint64(order.Uint32(tiff[4:8]))
	if offset+2 > uint64(len(tiff)) {
		return 0, errImage
	}
	count := uint64(order.Uint16(tiff[offset : offset+2]))
	if offset+2+12*count+4 > uint64(len(tiff)) {
		return 0, errImage
	}
	value := uint16(1)
	found := false
	for n := uint64(0); n < count; n++ {
		entry := tiff[offset+2+12*n : offset+14+12*n]
		if order.Uint16(entry[:2]) == 0x112 {
			if found || order.Uint16(entry[2:4]) != 3 || order.Uint32(entry[4:8]) != 1 {
				return 0, errImage
			}
			value, found = order.Uint16(entry[8:10]), true
			if value < 1 || value > 8 {
				return 0, errImage
			}
		}
	}
	return value, nil
}

func pngChunks(data []byte) ([]chunk, error) {
	if !bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")) {
		return nil, errImage
	}
	var result []chunk
	for pos := uint64(8); pos+12 <= uint64(len(data)); {
		size := uint64(binary.BigEndian.Uint32(data[pos : pos+4]))
		end := pos + 12 + size
		if end > uint64(len(data)) {
			return nil, errImage
		}
		if crc32.ChecksumIEEE(data[pos+4:pos+8+size]) != binary.BigEndian.Uint32(data[pos+8+size:end]) {
			return nil, errImage
		}
		c := chunk{string(data[pos+4 : pos+8]), data[pos+8 : pos+8+size]}
		result = append(result, c)
		if c.kind == "IEND" {
			if end != uint64(len(data)) || size != 0 {
				return nil, errImage
			}
			return result, nil
		}
		pos = end
	}
	return nil, errImage
}

func jpegChunks(data []byte) ([]chunk, error) {
	if len(data) < 4 || data[0] != 0xff || data[1] != 0xd8 || !bytes.HasSuffix(data, []byte{0xff, 0xd9}) {
		return nil, errImage
	}
	var result []chunk
	for pos := 2; pos < len(data); {
		start := pos
		if data[pos] != 0xff {
			return nil, errImage
		}
		for pos < len(data) && data[pos] == 0xff {
			pos++
		}
		if pos >= len(data) {
			return nil, errImage
		}
		kind := data[pos]
		pos++
		if kind == 0xd9 {
			if pos != len(data) {
				return nil, errMultipart
			}
			return append(result, chunk{"end", nil}), nil
		}
		if pos+2 > len(data) {
			return nil, errImage
		}
		size := int(binary.BigEndian.Uint16(data[pos : pos+2]))
		if size < 2 || pos+size > len(data) {
			return nil, errImage
		}
		if kind == 0xda {
			// Scan entropy until an unescaped marker. Further progressive scans
			// and metadata are parsed by the outer loop as ordinary segments.
			end := pos + size
			for end < len(data) {
				if data[end] != 0xff {
					end++
					continue
				}
				marker := end + 1
				for marker < len(data) && data[marker] == 0xff {
					marker++
				}
				if marker == len(data) {
					return nil, errImage
				}
				if data[marker] != 0 && (data[marker] < 0xd0 || data[marker] > 0xd7) {
					break
				}
				end = marker + 1
			}
			result = append(result, chunk{"scan", data[start:end]})
			pos = end
			continue
		}
		result = append(result, chunk{string([]byte{kind}), data[pos+2 : pos+size]})
		pos += size
	}
	return nil, errImage
}

func webpChunks(data []byte) ([]chunk, error) {
	if len(data) < 12 || string(data[:4]) != "RIFF" || string(data[8:12]) != "WEBP" || uint64(binary.LittleEndian.Uint32(data[4:8]))+8 != uint64(len(data)) {
		return nil, errImage
	}
	var result []chunk
	for pos := uint64(12); pos < uint64(len(data)); {
		if pos+8 > uint64(len(data)) {
			return nil, errImage
		}
		size := uint64(binary.LittleEndian.Uint32(data[pos+4 : pos+8]))
		end := pos + 8 + size
		if end+(size&1) > uint64(len(data)) {
			return nil, errImage
		}
		result = append(result, chunk{string(data[pos : pos+4]), data[pos+8 : end]})
		pos = end + (size & 1)
	}
	return result, nil
}

// cleanCandidate uses a format allowlist so optional descriptive blocks cannot
// survive even if a future vips saver starts writing another metadata field.
func cleanCandidate(data []byte, format string, original metadata) ([]byte, error) {
	var out bytes.Buffer
	// Validate before allowlisting: a known chunk name does not establish that
	// its payload is a display declaration rather than a private-data carrier.
	if _, err := inspectImage(data, format); err != nil {
		return nil, err
	}
	switch format {
	case "png":
		chunks, err := pngChunks(data)
		if err != nil {
			return nil, err
		}
		out.WriteString("\x89PNG\r\n\x1a\n")
		for _, c := range chunks {
			switch c.kind {
			case "IHDR", "PLTE", "IDAT", "IEND", "tRNS", "iCCP", "sBIT":
				writePNG(&out, c)
				if c.kind == "IHDR" {
					for _, color := range original.colors {
						writePNG(&out, color)
					}
				}
			}
		}
	case "jpeg":
		chunks, err := jpegChunks(data)
		if err != nil {
			return nil, err
		}
		out.Write([]byte{0xff, 0xd8})
		for _, c := range chunks {
			if c.kind == "end" {
				out.Write([]byte{0xff, 0xd9})
				continue
			}
			if c.kind == "scan" {
				out.Write(c.data)
				continue
			}
			kind := c.kind[0]
			if kind == 0xfe || kind >= 0xe0 && kind <= 0xef && !(kind == 0xe0 && bytes.HasPrefix(c.data, []byte("JFIF\x00"))) && !(kind == 0xee && bytes.HasPrefix(c.data, []byte("Adobe"))) && !(kind == 0xe2 && bytes.HasPrefix(c.data, []byte("ICC_PROFILE\x00"))) {
				continue
			}
			out.Write([]byte{0xff, kind})
			_ = binary.Write(&out, binary.BigEndian, uint16(len(c.data)+2))
			out.Write(c.data)
		}
	case "webp":
		chunks, err := webpChunks(data)
		if err != nil {
			return nil, err
		}
		out.WriteString("RIFF\x00\x00\x00\x00WEBP")
		for _, c := range chunks {
			switch c.kind {
			case "VP8X", "VP8 ", "VP8L", "ALPH", "ICCP":
				if c.kind == "VP8X" {
					if len(c.data) != 10 {
						return nil, errImage
					}
					c.data = bytes.Clone(c.data)
					c.data[0] &^= 0x0c // EXIF and XMP are absent from the new container.
				}
				out.WriteString(c.kind)
				_ = binary.Write(&out, binary.LittleEndian, uint32(len(c.data)))
				out.Write(c.data)
				if len(c.data)&1 != 0 {
					out.WriteByte(0)
				}
			}
		}
		binary.LittleEndian.PutUint32(out.Bytes()[4:8], uint32(out.Len()-8))
	default:
		return nil, errImage
	}
	return out.Bytes(), nil
}

func writePNG(out *bytes.Buffer, c chunk) {
	_ = binary.Write(out, binary.BigEndian, uint32(len(c.data)))
	out.WriteString(c.kind)
	out.Write(c.data)
	crc := crc32.NewIEEE()
	_, _ = crc.Write([]byte(c.kind))
	_, _ = crc.Write(c.data)
	_ = binary.Write(out, binary.BigEndian, crc.Sum32())
}
