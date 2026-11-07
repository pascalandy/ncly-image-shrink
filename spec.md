# image-shrink spec

The manifest holds command syntax and result types. Each line below has a spec ID. A line without a passing scenario stays under Open lines.

## ncly image shrink

**Success means** one smaller image exists beside the input, in the same format. The input stays unchanged. The output name inserts `.shrunk` before the original extension, preserving its spelling and case.

**Effects.** The run writes one output and private staging beside it. It uses no network, key, or paid service. The dry run reads and validates the input and checks the program, then reports its destination. It creates no stage. Compression and publication remain pending until the run.

**Resources.** Before decoding, reject an input above 64 MiB or 40 million pixels. The vips probe and operation share a 30-second deadline. vips uses one processing thread and a 64 MiB cache. That cache is not a hard process-memory limit. Decoding and metadata checks use the bounded input in memory.

**Publication policy.** Staging uses an unpredictable directory with mode 0700 and files with mode 0600. A new shared copy has mode 0600. Replacement preserves the existing regular file's permission bits. Publication uses a hard link for create and a rename for replace. A failure before publication keeps the existing output. A failure after publication keeps the completed output and reports its path and the error. Symbolic-link outputs and outputs hard-linked to the input are refused even with force.

**Failures.**

| ID | Condition | Code |
|---|---|---|
| `input-invalid` | Missing, empty, unreadable, truncated or malformed input is refused before vips starts, naming the input | `IMAGE_SHRINK_INPUT` |
| `unsupported` | A format outside JPEG, PNG or WebP, or an extension that disagrees with its bytes, is refused before vips starts | `IMAGE_SHRINK_FORMAT` |
| `budget` | Compressed input above 64 MiB, over 40 million pixels, or a dimension above 40000 is refused before full decoding | `IMAGE_SHRINK_LIMIT` |
| `options` | Invalid size, quality or timeout values are refused | `USAGE_INVALID` |
| `exists` | Existing output without force is refused before work and kept | `IMAGE_SHRINK_EXISTS` |
| `output-link` | A symbolic-link output, nonregular output or output hard-linked to the input is refused, with or without force | `IMAGE_SHRINK_OUTPUT` |
| `race` | Output appearing during create is kept unchanged and the run fails | `IMAGE_SHRINK_EXISTS` |
| `no-vips` | Missing or too-old vips and failed version probes receive the install hint | `PREREQ_MISSING` |
| `worker-fails` | Nonzero vips exit, missing output or timeout fails without publishing, keeping the old output and naming the input | `RUNTIME` |
| `invalid-result` | A damaged, wrong-format, wrong-size, metadata-bearing output or lost ICC is refused before publication | `IMAGE_SHRINK_RESULT` |
| `publish-fails` | Unsupported or failed publication fails closed and keeps any old output | `IMAGE_SHRINK_PUBLISH` |
| `interrupted` | Ctrl-C names the input, keeps old output and removes staging | `INTERRUPTED` |
| `terminated` | SIGTERM names the input, keeps old output and removes staging | `TERMINATED` |
| `multipart` | Animated PNG or WebP and multipicture JPEG are refused instead of discarding frames | `IMAGE_SHRINK_MULTIPART` |

**Limits.** None.

## Acceptance

Tests generate all images. The real acceptance input is a generated PNG with colored quadrants, transparency and a private location block. An optional fixture replaces it and receives the same size, format, privacy and geometry checks. Independent Go image decoders inspect complete pixel content. Container readers inspect metadata. References come from the generated image's declared pixels, orientation and color data.

| ID | Property |
|---|---|
| `not-smaller` | An optimized image whose candidate is not strictly smaller is refused with no new output, keeping an old one |
| `format` | JPEG, PNG and static WebP output keep the input format |
| `smaller` | The published copy has strictly fewer bytes than the input |
| `geometry` | The longest side fits the chosen bound, defaults to 2048, retains aspect ratio within pixel rounding, and never enlarges |
| `orientation` | All eight EXIF orientations are applied to the pixels before orientation metadata is removed |
| `png-color` | PNG gamma, chromaticity and sRGB declarations retain their display meaning |
| `alpha` | PNG and WebP retain transparency |
| `depth` | PNG retains its 16-bit sample precision |
| `names` | Paths with spaces, leading dashes, shell characters or loader-option characters are data, and output stays beside the input |
| `preview` | Dry run changes no file, shares input and conflict checks, reports output and overwrite, and marks compression and the publication operation pending |
| `staging` | Private exclusive staging cannot follow a planted temporary symlink and is removed on success, failure and signals |
| `replace` | Force atomically replaces only the completed output and keeps the replaced regular file's permissions |
| `pixels` | An independent decoder reads the complete output and matches the reference image, allowing bounded error for lossy formats and resizing |
| `quality` | JPEG and WebP use the selected lossy quality, default 80, while PNG uses lossless encoding |
| `profile` | The complete ICC color profile is retained |
| `privacy` | Location, author, camera and capture date in EXIF, XMP, IPTC and comments are removed. The output retains only color and pixel-decoding metadata |
| `publication` | Publication exposes only the completed output, keeps input bytes and permissions, and leaves no child process |

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
| PNG display declarations | Intake: gAMA, cHRM, sRGB and cICP layouts, values, cardinality and conflicts | `png-color`, `privacy` | `input-invalid`, `invalid-result` |
| Compatible sample models and profiles | Intake: Gray, Gray-alpha at 8/16 bits, RGB and CMYK profiles, orientation and resize | `profile`, `pixels`, `alpha`, `depth` | `input-invalid`, `invalid-result` |
| JPEG display declarations | Intake: Adobe RGB/YCbCr/CMYK/YCCK and JFIF | `pixels`, `profile` | `input-invalid` |
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

None. The generic cases and intake families remain selected for whole-family regression proof in round 2.
