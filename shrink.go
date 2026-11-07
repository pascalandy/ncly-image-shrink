package image

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pascalandy/nicely/sdk/command"
	"github.com/pascalandy/nicely/sdk/contract"
	"github.com/pascalandy/nicely/sdk/extension"
	"github.com/pascalandy/nicely/sdk/runner"
	"github.com/pascalandy/nicely/sdk/wrap"
	_ "golang.org/x/image/webp"
)

const maxInput = 64 << 20
const maxPixels = 40_000_000

type prepared struct {
	input, output, format string
	data                  []byte
	meta                  metadata
	config                image.Config
	side, quality         int
	timeout               time.Duration
	overwrite, force      bool
	mode                  os.FileMode
}

func problem(call extension.Call, code contract.Code, id string) contract.Outcome {
	fields := map[string]any{"File": call.Args[0]}
	return contract.Outcome{Errors: []contract.Problem{{
		Code: code, Message: call.Printer.T("image-shrink."+id, fields),
		Hint: call.Printer.T("image-shrink.hint", fields),
	}}}
}

func prepare(call extension.Call) (*prepared, contract.Outcome) {
	p := &prepared{side: 2048, quality: 80, timeout: 30 * time.Second, mode: 0o600, force: call.Flags["force"] == "true"}
	for flag, target := range map[string]*int{"max-side": &p.side, "quality": &p.quality} {
		if value, given := call.Flags[flag]; given {
			n, err := strconv.Atoi(value)
			if err != nil || n < 1 || flag == "quality" && n > 100 || flag == "max-side" && n > 40000 {
				return nil, problem(call, contract.UsageInvalid, "options")
			}
			*target = n
		}
	}
	if value, given := call.Flags["timeout"]; given {
		d, ok := command.Duration(value)
		if !ok || d > 30*time.Second {
			return nil, problem(call, contract.UsageInvalid, "options")
		}
		p.timeout = d
	}
	var err error
	p.input, err = filepath.Abs(call.Args[0])
	if err != nil {
		return nil, problem(call, "IMAGE_SHRINK_INPUT", "input")
	}
	ext := filepath.Ext(p.input)
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		p.format = "jpeg"
	case ".png":
		p.format = "png"
	case ".webp":
		p.format = "webp"
	default:
		return nil, problem(call, "IMAGE_SHRINK_FORMAT", "format")
	}
	info, err := os.Lstat(p.input)
	if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0o444 == 0 {
		return nil, problem(call, "IMAGE_SHRINK_INPUT", "input")
	}
	if info.Size() > maxInput {
		return nil, problem(call, "IMAGE_SHRINK_LIMIT", "limit")
	}
	file, err := os.Open(p.input)
	if err != nil {
		return nil, problem(call, "IMAGE_SHRINK_INPUT", "input")
	}
	p.data, err = io.ReadAll(io.LimitReader(file, maxInput+1))
	closeErr := file.Close()
	if err != nil || closeErr != nil {
		return nil, problem(call, "IMAGE_SHRINK_INPUT", "input")
	}
	if len(p.data) > maxInput {
		return nil, problem(call, "IMAGE_SHRINK_LIMIT", "limit")
	}
	actual := ""
	switch {
	case bytes.HasPrefix(p.data, []byte("\x89PNG\r\n\x1a\n")):
		actual = "png"
	case bytes.HasPrefix(p.data, []byte{0xff, 0xd8}):
		actual = "jpeg"
	case len(p.data) >= 12 && string(p.data[:4]) == "RIFF" && string(p.data[8:12]) == "WEBP":
		actual = "webp"
	}
	if len(p.data) != 0 && actual != p.format {
		return nil, problem(call, "IMAGE_SHRINK_FORMAT", "format")
	}
	p.meta, err = inspectImage(p.data, p.format)
	if err == errMultipart {
		return nil, problem(call, "IMAGE_SHRINK_MULTIPART", "multipart")
	}
	if err != nil {
		return nil, problem(call, "IMAGE_SHRINK_INPUT", "input")
	}
	var format string
	p.config, format, err = image.DecodeConfig(bytes.NewReader(p.data))
	if err != nil {
		return nil, problem(call, "IMAGE_SHRINK_INPUT", "input")
	}
	if format != p.format {
		return nil, problem(call, "IMAGE_SHRINK_FORMAT", "format")
	}
	w, h := p.config.Width, p.config.Height
	if w <= 0 || h <= 0 || w > 40000 || h > 40000 || int64(w)*int64(h) > maxPixels {
		return nil, problem(call, "IMAGE_SHRINK_LIMIT", "limit")
	}
	if _, _, err := image.Decode(bytes.NewReader(p.data)); err != nil {
		return nil, problem(call, "IMAGE_SHRINK_INPUT", "input")
	}
	p.output = strings.TrimSuffix(p.input, ext) + ".shrunk" + ext
	if output, err := os.Lstat(p.output); err == nil {
		if !output.Mode().IsRegular() || os.SameFile(info, output) {
			return nil, problem(call, "IMAGE_SHRINK_OUTPUT", "output")
		}
		if !p.force {
			return nil, problem(call, "IMAGE_SHRINK_EXISTS", "exists")
		}
		p.overwrite, p.mode = true, output.Mode().Perm()
	} else if !os.IsNotExist(err) {
		return nil, problem(call, "IMAGE_SHRINK_OUTPUT", "output")
	}
	return p, contract.Outcome{}
}

func shrink(ctx context.Context, call extension.Call) (outcome contract.Outcome) {
	p, failure := prepare(call)
	if p == nil {
		return failure
	}
	base := map[string]any{"output": p.output, "input_bytes": len(p.data)}
	if call.Flags["dry-run"] == "true" {
		return wrap.Run(ctx, call, wrap.Program{Name: "vips", Args: []string{"--version"}, Timeout: p.timeout, Input: call.Args[0], DiscardStdout: true}, func(end wrap.Ended) contract.Outcome {
			if end.Exit != 0 {
				return end.Failed()
			}
			base["overwrite"] = p.overwrite
			op := "link"
			if p.force {
				op = "rename"
			}
			base["pending_checks"] = []string{"output.smaller", "publication." + op}
			_, _ = fmt.Fprintln(call.Stdout, call.Printer.T("image-shrink.preview", map[string]any{"Output": p.output}))
			return contract.Outcome{Data: base}
		})
	}
	stage, err := os.MkdirTemp(filepath.Dir(p.output), ".ncly-image-shrink-")
	if err != nil {
		return problem(call, "IMAGE_SHRINK_PUBLISH", "publish")
	}
	defer func() {
		if err := os.RemoveAll(stage); err != nil {
			cleanup := problem(call, "IMAGE_SHRINK_PUBLISH", "cleanup")
			outcome.Errors = append(outcome.Errors, cleanup.Errors...)
		}
	}()
	suffix := "." + p.format
	if p.format == "jpeg" {
		suffix = ".jpg"
	}
	source, candidate := filepath.Join(stage, "input"+suffix), filepath.Join(stage, "output"+suffix)
	sourceData := p.data
	if p.format == "png" {
		sourceData, err = preparePNG(p)
		if err != nil {
			return problem(call, "IMAGE_SHRINK_RESULT", "result")
		}
	}
	if err := os.WriteFile(source, sourceData, 0o600); err != nil {
		return problem(call, "IMAGE_SHRINK_PUBLISH", "publish")
	}
	file, err := os.OpenFile(candidate, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return problem(call, "IMAGE_SHRINK_PUBLISH", "publish")
	}
	if err := file.Close(); err != nil {
		return problem(call, "IMAGE_SHRINK_PUBLISH", "publish")
	}
	options := "keep=icc"
	switch p.format {
	case "jpeg":
		options += ",optimize_coding,Q=" + strconv.Itoa(p.quality)
	case "webp":
		options += ",Q=" + strconv.Itoa(p.quality)
	case "png":
		options += ",compression=9"
	}
	program := wrap.Program{
		Name: "vips", Timeout: p.timeout, Input: call.Args[0], DiscardStdout: true,
		Args: []string{"thumbnail", source + "[fail_on=warning]", candidate + "[" + options + "]", strconv.Itoa(p.side), "--size", "down", "--vips-concurrency", "1", "--vips-cache-max-memory", "64m"},
		Env:  map[string]string{"VIPS_BLOCK_UNTRUSTED": "1"},
	}
	if p.format == "png" {
		program.Args = []string{"copy", source + "[fail_on=warning]", candidate + "[" + options + "]", "--vips-concurrency", "1", "--vips-cache-max-memory", "64m"}
	}
	if p.format == "jpeg" {
		// thumbnail converts CMYK to sRGB even without an ICC transform.
		// jpegload's autorotate and resize retain native Gray/RGB/CMYK bands.
		scale := math.Min(1, float64(p.side)/float64(max(p.config.Width, p.config.Height)))
		program.Args = []string{"resize", source + "[fail_on=warning,autorotate]", candidate + "[" + options + "]", strconv.FormatFloat(scale, 'g', -1, 64), "--vips-concurrency", "1", "--vips-cache-max-memory", "64m"}
	}
	return wrap.Run(ctx, call, program, func(end wrap.Ended) contract.Outcome {
		if end.Exit != 0 {
			return end.Failed()
		}
		data, err := readCandidate(candidate)
		if err != nil {
			return end.Malformed()
		}
		data, err = cleanCandidate(data, p.format, p.meta)
		if err != nil || !validCandidate(data, p) {
			return problem(call, "IMAGE_SHRINK_RESULT", "result")
		}
		if len(data) >= len(p.data) {
			return problem(call, "IMAGE_SHRINK_NOT_SMALLER", "not_smaller")
		}
		if ctx.Err() != nil {
			return contract.Outcome{Signal: contract.StopSignal(runner.SignalOf(ctx))}
		}
		if err := os.WriteFile(candidate, data, 0o600); err != nil {
			return problem(call, "IMAGE_SHRINK_PUBLISH", "publish")
		}
		if err := os.Chmod(candidate, p.mode); err != nil {
			return problem(call, "IMAGE_SHRINK_PUBLISH", "publish")
		}
		if p.force {
			err = os.Rename(candidate, p.output)
		} else {
			err = os.Link(candidate, p.output)
		}
		if os.IsExist(err) {
			return problem(call, "IMAGE_SHRINK_EXISTS", "exists")
		}
		if err != nil {
			return problem(call, "IMAGE_SHRINK_PUBLISH", "publish")
		}
		base["output_bytes"] = len(data)
		_, _ = fmt.Fprintln(call.Stdout, call.Printer.T("image-shrink.done", map[string]any{"Output": p.output, "Bytes": len(data)}))
		return contract.Outcome{Data: base}
	})
}

func readCandidate(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > maxInput {
		return nil, os.ErrInvalid
	}
	return os.ReadFile(path)
}

func validCandidate(data []byte, p *prepared) bool {
	if p.format == "png" && len(p.data) > 24 && p.data[24] == 16 && (len(data) < 25 || data[24] != 16) {
		return false
	}
	meta, err := inspectImage(data, p.format)
	if err != nil || !bytes.Equal(meta.icc, p.meta.icc) {
		return false
	}
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || format != p.format || config.Width < 1 || config.Height < 1 {
		return false
	}
	w, h := p.config.Width, p.config.Height
	if p.meta.orientation >= 5 {
		w, h = h, w
	}
	scale := math.Min(1, float64(p.side)/float64(max(w, h)))
	if config.Width > w || config.Height > h || max(config.Width, config.Height) > p.side || math.Abs(float64(config.Width)-float64(w)*scale) > 1 || math.Abs(float64(config.Height)-float64(h)*scale) > 1 {
		return false
	}
	_, _, err = image.Decode(bytes.NewReader(data))
	return err == nil
}
