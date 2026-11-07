package image

import (
	"bytes"
	"compress/zlib"
	"image"
	"image/png"
	"math"

	"golang.org/x/image/draw"
)

// thumbnail converts PNG samples to eight bits. Prepare its pixels in Go and
// use the lossless vips PNG saver so sixteen-bit inputs keep their precision.
func preparePNG(p *prepared) ([]byte, error) {
	input, _, err := image.Decode(bytes.NewReader(p.data))
	if err != nil {
		return nil, err
	}
	w, h := input.Bounds().Dx(), input.Bounds().Dy()
	if p.meta.orientation != 1 {
		ow, oh := w, h
		if p.meta.orientation >= 5 {
			ow, oh = h, w
		}
		var oriented draw.Image = image.NewNRGBA(image.Rect(0, 0, ow, oh))
		if p.data[24] == 16 {
			oriented = image.NewNRGBA64(image.Rect(0, 0, ow, oh))
		}
		for y := 0; y < oh; y++ {
			for x := 0; x < ow; x++ {
				sx, sy := x, y
				switch p.meta.orientation {
				case 2:
					sx = w - 1 - x
				case 3:
					sx, sy = w-1-x, h-1-y
				case 4:
					sy = h - 1 - y
				case 5:
					sx, sy = y, x
				case 6:
					sx, sy = y, h-1-x
				case 7:
					sx, sy = w-1-y, h-1-x
				case 8:
					sx, sy = w-1-y, x
				}
				oriented.Set(x, y, input.At(sx, sy))
			}
		}
		input, w, h = oriented, ow, oh
	}
	if max(w, h) > p.side {
		scale := float64(p.side) / float64(max(w, h))
		nw, nh := max(1, int(math.Round(float64(w)*scale))), max(1, int(math.Round(float64(h)*scale)))
		resized := image.NewNRGBA64(image.Rect(0, 0, nw, nh))
		draw.CatmullRom.Scale(resized, resized.Bounds(), input, input.Bounds(), draw.Src, nil)
		input = resized
	}
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, input); err != nil {
		return nil, err
	}
	chunks, err := pngChunks(encoded.Bytes())
	if err != nil {
		return nil, err
	}
	var output bytes.Buffer
	output.WriteString("\x89PNG\r\n\x1a\n")
	for _, c := range chunks {
		writePNG(&output, c)
		if c.kind == "IHDR" {
			if len(p.meta.icc) != 0 {
				var profile bytes.Buffer
				compressed := zlib.NewWriter(&profile)
				if _, err := compressed.Write(p.meta.icc); err != nil {
					return nil, err
				}
				if err := compressed.Close(); err != nil {
					return nil, err
				}
				writePNG(&output, chunk{"iCCP", append([]byte("Color\x00\x00"), profile.Bytes()...)})
			}
			for _, color := range p.meta.colors {
				writePNG(&output, color)
			}
		}
	}
	return output.Bytes(), nil
}
