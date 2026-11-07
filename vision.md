# image-shrink vision

This file is a form filled for the example command, `ncly image shrink`.
It inspects a file until you replace it.
Replace the text of each part, then delete each line that starts with "Fill:".

Fill: one sentence that states the job of the extension, and for whom.

`ncly image shrink` gives the size and the line count of a file, so an agent knows what it will read. The extension meets "What a great extension does", in the contract that `just contract` prints.

## Wraps

Fill: the source and the revision of each program or script that a command wraps, such as a Git URL and a commit.

- `worker/inspect.py`, which the extension template wrote

## Intake

`just records` writes this part from `intake.toml`. Change that file, then run `just records` again.

- Program: vips thumbnail reads one local image and writes the path selected by its output suffix. CLI options follow the operation. It does not ask a question. Source: https://github.com/libvips/libvips/blob/v8.18.7/doc/using-the-cli.md, at v8.18.7. Read: Optional arguments and Implicit file format conversion, including bracketed loader and saver options
- Geometry: thumbnail fits inside a square, applies EXIF orientation by default, and size=down prevents enlargement. Without an output_profile it does not request an ICC transform. Source: https://github.com/libvips/libvips/blob/v8.18.7/libvips/resample/thumbnail.c, at v8.18.7. Read: vips_thumbnail documentation and vips_thumbnail_build, especially needs_icc_transform, size, and auto_rotate
- Privacy: keep=icc drops EXIF, XMP, IPTC and optional comments while retaining ICC. An incompatible ICC can be silently dropped, so the adapter must verify preservation before publication. Source: https://github.com/libvips/libvips/blob/v8.18.7/libvips/foreign/foreign.c, at v8.18.7. Read: vips_foreign_get_keep, vips_foreign_save_remove_metadata and vips__foreign_update_metadata
- JPEG quality is Q, with optimized Huffman coding available. The output is lossy. Source: https://github.com/libvips/libvips/blob/v8.18.7/libvips/foreign/jpegsave.c, at v8.18.7. Read: vips_jpegsave documentation, Q and optimize_coding
- PNG uses lossless compression from 0 to 9. Palette quantization is optional and must remain disabled. RGB, alpha, monochrome and 16-bit samples are supported. Source: https://github.com/libvips/libvips/blob/v8.18.7/libvips/foreign/pngsave.c, at v8.18.7. Read: vips_pngsave documentation, compression, palette and bitdepth
- WebP defaults to lossy compression controlled by Q. Alpha quality defaults to 100. Animated input can have several frames, but thumbnail defaults to only the first frame, so the adapter must refuse animation. Source: https://github.com/libvips/libvips/blob/v8.18.7/libvips/foreign/webpsave.c and libvips/resample/thumbnail.c, at v8.18.7. Read: vips_webpsave documentation, Q, alpha_q and lossless, and vips_thumbnail documentation on first-frame loading
- Exit codes: vips exits 0 on success and nonzero on error. Its loaders are permissive by default. fail_on=warning requests a refusal instead of publishing damaged input. Source: https://github.com/libvips/libvips/blob/v8.18.7/man/vips.1 and libvips/resample/thumbnail.c, at v8.18.7. Read: RETURN VALUE and vips_thumbnail documentation on fail_on
- The wrapper needs no key, login, network, or payment. It accepts only validated JPEG, PNG and static WebP bytes, uses fixed private snapshot names, and never relays an operation or loader option from a user path. It removes the snapshot after the call. Inference: no source gives it
- Resource policy: cap compressed input at 64 MiB and expanded input at 40 million pixels before decoding. Limit the vips call, including its version probe, to 30 seconds. One processing thread and a 64 MiB operation cache reduce memory use, but that cache is not an operating-system memory ceiling. Inference: no source gives it
- Display declarations to cover are EXIF orientation, ICC color profiles, PNG gamma and chromaticity, PNG bit depth and transparency, JPEG CMYK and Adobe markers, and WebP alpha. Multipage JPEG and animated PNG or WebP must be refused instead of losing parts. Inference: no source gives it

**Experiments.**

- `vips --version`: Exit 0, vips-8.18.7 on om1. stderr also warns that the optional openslide plugin cannot load libopenslide.so.1. The selected JPEG, PNG and WebP operations do not use that plugin.

## Refuses

Fill: what the extension will not do, and why.

- It never changes a file: it only reads
