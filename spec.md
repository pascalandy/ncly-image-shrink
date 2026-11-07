# image-shrink spec

The manifest holds command syntax and result types. Each line below has a spec ID. A line without a passing scenario stays under Open lines.

## ncly image shrink

**Success means** one smaller image exists beside the input, in the same format. The input stays unchanged. The output name inserts `.shrunk` before the original extension, preserving its spelling and case.

**Effects.** The run writes one output and private staging beside it. It uses no network, key, or paid service. The dry run reads and validates the input and checks the program, then reports its destination. It creates no stage. Compression and publication remain pending until the run.

**Resources.** Before decoding, reject an input above 64 MiB or 40 million pixels. The vips probe and operation share a 30-second deadline. vips uses one processing thread and a 64 MiB cache. That cache is not a hard process-memory limit. Decoding and metadata checks use the bounded input in memory.

**Publication policy.** Staging uses an unpredictable directory with mode 0700 and files with mode 0600. A new shared copy has mode 0600. Replacement preserves the existing regular file's permission bits. Publication uses a hard link for create and a rename for replace. A failure before publication keeps the existing output. A failure after publication keeps the completed output and reports its path and the error. Symbolic-link outputs and outputs hard-linked to the input are refused even with force.

**Failures.**

None proven yet.

**Limits.** None.

## Acceptance

Tests generate all images. The real acceptance input is a generated PNG with colored quadrants, transparency and a private location block. An optional fixture replaces it and receives the same size, format, privacy and geometry checks. Independent Go image decoders inspect complete pixel content. Container readers inspect metadata. References come from the generated image's declared pixels, orientation and color data.

None proven yet.

## Case matrix

| Family | Source | Accepted | Refused |
|---|---|---|---|
| A readable local input | Input: Empty; Input: Truncated; Input: Unreadable | `pixels` | `input-invalid` |
| Only JPEG, PNG and WebP | Input: Of the wrong type; Input: With active content | `format` | `unsupported` |
| Bounded work | Input: Very large; Input: Large expansion | `geometry` | `budget` |
| One image, without animation | Input: With several parts | `format` | `multipart` |
| Safe file names | File name: Spaces; File name: Leading dash; File name: Shell characters | `names` | `input-invalid` |
| Existing output | Output: Exists, without --force | `replace` | `exists` |
| Link output | Output: Is a link to the input | `pixels` | `output-link` |
| Concurrent output | Output: Appears during the run, without --force | `publication` | `race` |
| Private staging | Output: Temporary file | `staging` | `worker-fails` |
| Descriptive metadata | Output: Keeps private data of the input | `privacy` | `invalid-result` |
| Orientation and color | Output: Keeps data that a reader needs; Output: Loses data that changes the meaning | `orientation`, `profile`, `png-color`, `alpha`, `depth` | `invalid-result` |
| Strict byte saving | Intake: output size | `smaller` | `not-smaller` |
| Interruption and cleanup | Interruption: SIGINT or SIGTERM during the work; Interruption: Timeout | `publication` | `interrupted`, `terminated`, `worker-fails` |
| Faithful dry run | Other: Dry run | `preview` | `input-invalid`, `exists` |
| Installed program | Other: Missing program | `pixels` | `no-vips`, `worker-fails` |
| Valid command options | Intake: bounded geometry and quality | `geometry`, `quality` | `options` |
| Publication errors | Publication: create and replace | `publication`, `replace` | `publish-fails` |

## Publication

| Command | Mode | Destination | Overwrite | Protected | Operation | Dry run |
|---|---|---|---|---|---|---|
| `ncly image shrink` | create | `output` | `overwrite` | `file` | `link` | pending |
| `ncly image shrink` | replace | `output` | `overwrite` | `file` | `rename` | pending |

## Open lines

| ID | Command | Part | Case | Code |
|---|---|---|---|---|
| `input-invalid` | `ncly image shrink` | Failures | Missing, empty, unreadable, truncated or malformed input is refused before vips starts, naming the input | `IMAGE_SHRINK_INPUT` |
| `unsupported` | `ncly image shrink` | Failures | A format outside JPEG, PNG or WebP, or an extension that disagrees with its bytes, is refused before vips starts | `IMAGE_SHRINK_FORMAT` |
| `budget` | `ncly image shrink` | Failures | Compressed input above 64 MiB, over 40 million pixels, or a dimension above 40000 is refused before full decoding | `IMAGE_SHRINK_LIMIT` |
| `multipart` | `ncly image shrink` | Failures | Animated PNG or WebP and multipicture JPEG are refused instead of discarding frames | `IMAGE_SHRINK_MULTIPART` |
| `options` | `ncly image shrink` | Failures | Invalid size, quality or timeout values are refused | `USAGE_INVALID` |
| `exists` | `ncly image shrink` | Failures | Existing output without force is refused before work and kept | `IMAGE_SHRINK_EXISTS` |
| `output-link` | `ncly image shrink` | Failures | A symbolic-link output, nonregular output or output hard-linked to the input is refused, with or without force | `IMAGE_SHRINK_OUTPUT` |
| `race` | `ncly image shrink` | Failures | Output appearing during create is kept unchanged and the run fails | `IMAGE_SHRINK_EXISTS` |
| `no-vips` | `ncly image shrink` | Failures | Missing or too-old vips and failed version probes receive the install hint | `PREREQ_MISSING` |
| `worker-fails` | `ncly image shrink` | Failures | Nonzero vips exit, missing output or timeout fails without publishing, keeping the old output and naming the input | `RUNTIME` |
| `invalid-result` | `ncly image shrink` | Failures | A damaged, wrong-format, wrong-size, metadata-bearing output or lost ICC is refused before publication | `IMAGE_SHRINK_RESULT` |
| `publish-fails` | `ncly image shrink` | Failures | Unsupported or failed publication fails closed and keeps any old output | `IMAGE_SHRINK_PUBLISH` |
| `not-smaller` | `ncly image shrink` | Acceptance | An optimized image whose candidate is not strictly smaller is refused with no new output, keeping an old one | `IMAGE_SHRINK_NOT_SMALLER` |
| `format` | `ncly image shrink` | Acceptance | JPEG, PNG and static WebP output keep the input format | - |
| `smaller` | `ncly image shrink` | Acceptance | The published copy has strictly fewer bytes than the input | - |
| `pixels` | `ncly image shrink` | Acceptance | An independent decoder reads the complete output and matches the reference image, allowing bounded error for lossy formats and resizing | - |
| `geometry` | `ncly image shrink` | Acceptance | The longest side fits the chosen bound, defaults to 2048, retains aspect ratio within pixel rounding, and never enlarges | - |
| `quality` | `ncly image shrink` | Acceptance | JPEG and WebP use the selected lossy quality, default 80, while PNG uses lossless encoding | - |
| `orientation` | `ncly image shrink` | Acceptance | All eight EXIF orientations are applied to the pixels before orientation metadata is removed | - |
| `profile` | `ncly image shrink` | Acceptance | The complete ICC color profile is retained | - |
| `png-color` | `ncly image shrink` | Acceptance | PNG gamma, chromaticity and sRGB declarations retain their display meaning | - |
| `alpha` | `ncly image shrink` | Acceptance | PNG and WebP retain transparency | - |
| `depth` | `ncly image shrink` | Acceptance | PNG retains its 16-bit sample precision | - |
| `privacy` | `ncly image shrink` | Acceptance | Location, author, camera and capture date in EXIF, XMP, IPTC and comments are removed. The output retains only color and pixel-decoding metadata | - |
| `names` | `ncly image shrink` | Acceptance | Paths with spaces, leading dashes, shell characters or loader-option characters are data, and output stays beside the input | - |
| `preview` | `ncly image shrink` | Acceptance | Dry run changes no file, shares input and conflict checks, reports output and overwrite, and marks compression and the publication operation pending | - |
| `staging` | `ncly image shrink` | Acceptance | Private exclusive staging cannot follow a planted temporary symlink and is removed on success, failure and signals | - |
| `publication` | `ncly image shrink` | Acceptance | Publication exposes only the completed output, keeps input bytes and permissions, and leaves no child process | - |
| `replace` | `ncly image shrink` | Acceptance | Force atomically replaces only the completed output and keeps the replaced regular file's permissions | - |
| `interrupted` | `ncly image shrink` | Failures | Ctrl-C names the input, keeps old output and removes staging | `INTERRUPTED` |
| `terminated` | `ncly image shrink` | Failures | SIGTERM names the input, keeps old output and removes staging | `TERMINATED` |
