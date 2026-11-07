Commit: 27536df4a96fb802629f749e1a900ca1371ac5c9
Reviewer: readiness-r2, delegation 8, Cook 1.1 brigade 461d; independent GPT-6.1-Sol via Codex, xhigh/default; wrote none of the implementation
Verdict: Needs work

# Readiness review

`just records` writes this file from `review.toml`. Change that file, then run `just records` again.

## Findings

### R1-F1 Unvalidated PNG display declarations publish private text

- **Severity**: `block`
- **State**: `fixed`
- **Spec**: `privacy`
- **Family**: Descriptive metadata
- **Commit**: a13b15e229eb66e8558e071bc01cd3e2dbd4b6cd
- **Evidence**: Preserves manic-2 F1/F4 attribution and the original block grade. a13b15e validates retained PNG display declarations before input use and candidate allowlisting; the preceding failing scenario supplies length/value/duplicate/conflict and privacy-reader controls. Subsequent d503d24/60fe121 complete declaration/model/placement handling. At frozen 27536df, all four preserved private-bearing gAMA/cHRM/sRGB/cICP seeds exit 2 with IMAGE_SHRINK_INPUT and no output. Fresh generated valid declarations succeed and pass independent pixels/display/privacy checks; length/value/duplicate variants refuse, and a readable malicious declaration fails the reader. Full committed scenario pairs include conflicts and legal post-IDAT EXIF. See readiness-r2/evidence/fresh-probes.json, second-pass.json and sources; both independent round-2 full-suite logs at 6f2 corroborate family execution. No new private-text leak was established by the empty ICC defect.
- **Command**: ncly image shrink <file.png>, including supplied-fixture acceptance
- **Consequence**: The original private-bearing PNG declaration leak and false certification no longer reproduce. Complete ICC retention and malformed ICC presence are assessed separately under R1-F2/R2-F2.
- **Owner**: Cook 1.1, extension parser and independent acceptance reader
- **Fix**: Implemented declaration validation and whole-family success/refusal/privacy controls. Retain them and the original privacy promise.

### R1-F2 Native sample-model and profile family still lacks complete proof

- **Severity**: `block`
- **State**: `open`
- **Spec**: `profile`
- **Family**: Orientation and color
- **Evidence**: Preserves manic-1 F1/manic-2 F2 attribution and the original critical-to-block grade. d503d24 repairs the observed native-model failures: the preserved Gray orientation/resize and CMYK fixtures now succeed and pass the current full reader, as both independent passes observed and this session reran for the copied controls. Gray/Gray-alpha 8/16, Gray JPEG, RGB/CMYK, Adobe/JFIF and malformed/incompatible profile pairs are committed. Fresh LittleCMS opens and transforms both generated model profiles, and complete ICC equality remains enforced. Closure is withheld under the whole-family rule: legal Gray+tRNS is absent from the success family and loses alpha, including a fresh GRAY-ICC orientation-2 control whose full profile remains equal; empty declared ICC profiles also bypass the family refusal boundary and its reader. These concrete defects are detailed as R2-F1 and R2-F2, not as a claim that the original conversion refusals still reproduce. See readiness-r2/evidence/fresh-probes.json, second-pass.json, fresh-lcms.log; manic-1-r2/evidence/gray-trns-audit.jsonl; manic-2-r2/evidence/families-results.json. The unchanged profile acceptance row still needs reopening as specified in report.md.
- **Command**: ncly image shrink <Gray PNG or CMYK JPEG>, with orientation or resizing as applicable
- **Consequence**: The original samples/profile conversion repair works on the reported controls, but the promised whole-family repair proof is incomplete. Successful damaged Gray transparency and malformed profile acceptance prevent marking this broad earlier finding fixed. This record overlaps the two detailed new subjects and is not a fourth independent product defect.
- **Owner**: Cook 1.1, PNG preparation and JPEG processing
- **Fix**: Finish the existing orientation/color and compatible-model/profile family: preserve Gray+tRNS alpha without changing GRAY ICC compatibility, validate every declared ICC including zero length, retain full ICC equality, and add whole-family accepted/refused pairs plus independent content controls. Keep this finding open until those pairs prove the family. Keep the profile promise unchanged.

### R1-F3 The fidelity reference aliases source footprints and accepts gross corruption

- **Severity**: `block`
- **State**: `open`
- **Spec**: `pixels`
- **Family**: Orientation and color
- **Evidence**: Preserves original manic-1 F2 attribution, the fidelity portion of manic-2 F3 and the existing block grade; M2-R2-F2 is a residual of this ID. 2cb63ea removes skipped destination pixels and rejects the original 1008-pixel bad variant, but acceptance.txtar:1163-1191 averages only four source points at offsets .25/.75. At frozen 27536df a fresh 256x128 alternating-column source has mean 127.5, its real 64x32 output has mean 127.96875 with all 2048 pixels nonzero, yet the reader rejects it at pixel 0,0 (35209 versus bound 10280). The correctly sized all-black substitute passes and prints complete-pixel/profile/privacy success. Fresh orientation-6, 16-bit and alpha variants repeat the split, including fully transparent false certification. See readiness-r2/evidence/fresh-pixel-audit.log, stripes-reader-real.log, stripes-reader-corrupt.log, second-pass.json and retained fixtures; manic-2-r2/evidence/pixel-audit.json and second-pass-results.json. This is an acceptance defect, not observed encoder corruption.
- **Command**: Independent acceptance checks for ncly image shrink resized PNG, JPEG and WebP outputs
- **Consequence**: The mandatory independent primary-output proof can certify an entirely black or transparent substitute and reject a faithful real reduction. Visiting every destination pixel does not repair an aliased source reference.
- **Owner**: Cook 1.1, independent acceptance reader and references
- **Fix**: Replace the four-point estimate with an independently justified complete source-footprint reference, accounting for fractional edges, orientation and premultiplied alpha. Keep every-output-pixel and whole-image bounds with justified codec/kernel tolerances. Add generated high-frequency success/refusal pairs for reduction factors above two, orientation, depth and alpha. Keep pixels open.

### R1-F4 Multipart classifiers lack complete generated acceptance pairs

- **Severity**: `block`
- **State**: `fixed`
- **Spec**: `multipart`
- **Family**: One image, without animation
- **Commit**: 2cb63ea792272ca441509279fa8b3f049c4a0493
- **Evidence**: Preserves both original reviewers attribution and block grade. The committed generator now makes actual APNG fcTL/fdAT frames, WebP ANIM/ANMF frame payloads, flag-only/chunk-only branches, MPO MP-entry targets and concatenated JPEGs, alongside still-image success cases. At frozen 27536df this session freshly generates the fixtures, ffprobe decodes two APNG frames, ImageMagick decodes two WebP frames, and all seven public direct-CLI multipart probes refuse with IMAGE_SHRINK_MULTIPART and absent output. Both independent round-2 suites at 6f2 pass the matching host-boundary families. Their independent TIFF audit confirms sizes 7928/7838 and actual JPEG targets at offsets 0/7928, with the second TIFF-relative offset 7918. See readiness-r2/evidence/second-pass.json and sources; manic-2-r2/evidence/second-pass-results.json and families-results.json.
- **Command**: ncly image shrink <JPEG, PNG or WebP>
- **Consequence**: The earlier missing multipart regression family is now committed and demonstrated. No frame-loss product bug was claimed or discovered.
- **Owner**: Cook 1.1, generated fixtures and case matrix
- **Fix**: Implemented genuine still/multipart acceptance and refusal pairs. Retain the independent fixture validity checks and fail-closed classifiers.

### R1-F5 Lossy quality overrides have no successful independent proof

- **Severity**: `block`
- **State**: `fixed`
- **Spec**: `quality`
- **Family**: Valid command options
- **Commit**: 2cb63ea792272ca441509279fa8b3f049c4a0493
- **Evidence**: Preserves both original reviewers attribution and block grade. Committed JPEG/WebP public calls at quality 20 and 95 independently compare every decoded pixel, profile/privacy and sizes. Fresh frozen-head direct CLI executions repeat JPEG MSE 984.016 versus 15.180 and WebP 1161.724 versus 868.137; higher quality is larger and closer. Identical-output ignored-override controls fail for both formats. Invalid options, default 80 and lossless PNG controls remain in the passing 6f2 suites. See readiness-r2/evidence/second-pass.json and second-quality-* logs; manic-2-r2/evidence/families-results.json. These are deterministic fixture metrics, not general visual-quality benchmarks.
- **Command**: ncly image shrink <JPEG or WebP> --quality <1..100>
- **Consequence**: The missing accepted lossy-quality proof is supplied. Hardcoding a single quality no longer passes these pairs.
- **Owner**: Cook 1.1, lossy-quality acceptance family
- **Fix**: Implemented materially different valid-quality calls, complete decoded metrics and ignored-override negatives, alongside option refusals. Retain the selected-quality promise.

### R1-F6 Post-publication signals lack the required survival proof

- **Severity**: `block`
- **State**: `fixed`
- **Spec**: `publication`
- **Family**: Publication errors
- **Commit**: 115e68dddafb43464a5d89fb5f970764a7ed9303
- **Evidence**: Preserves readiness-r1 delegation 5 attribution and block grade. acceptance.txtar:458 onward commits create/INT and replace/TERM using survives and signal. The independent published observer fully reads the new, strictly smaller destination before writing the marker; the old replacement copy cannot satisfy that condition. Ordinary scratch cleanup provides the interval without a production hook; a missed interval makes the test fail. Input bytes/mode, completed destination bytes/mode, output path/error, stage cleanup and child exit are checked. A committed forced cleanup-failure case retains and reads the completed output with IMAGE_SHRINK_PUBLISH and its destination. Both independent fresh 6f2 full suites execute these cases successfully, matching Cook round2-postpublication.log and repeat log. Source and raw logs were independently inspected here; this session did not replay the installed host lifecycle suite. Native preparation signal measurements are separate and do not prove final-validation cancellation.
- **Command**: ncly image shrink <file>, create and force-replace publication lifecycle
- **Consequence**: The explicit post-publication acceptance gap is closed by committed real publication observations and survival checks. No signal or cleanup loss was reproduced.
- **Owner**: Cook 1.1, publication lifecycle acceptance family
- **Fix**: Implemented the required post-commit survival proof and cleanup-failure control. Retain the pre-commit controls and original publication policy.

### R2-F1 Gray PNG tRNS transparency is discarded during native-model preparation

- **Severity**: `block`
- **State**: `open`
- **Spec**: `alpha`
- **Family**: Orientation and color
- **Evidence**: M1-R2-F1, critical, newly reproduced at 27536df; regression in d503d24. png.go:70 selects alpha only when the original color type is 4, and :125 omits decoded alpha for legal type-0 tRNS. Four 8/16-bit orientation-1/2 inputs contain 1024 transparent and 1024 opaque samples; every fresh published output has zero transparent and 2048 opaque samples. Equivalent explicit Gray-alpha type-4 controls retain transparency. Max-side 32 reproduces the split. CRC/decompression audits independently count every raw tRNS sample. A fresh GRAY-ICC control retains its complete compatible profile but still loses alpha. The current reader rejects the damaged output; manic-1 actual supplied-fixture acceptance intentionally fails. See readiness-r2/evidence/fresh-alpha-audit.log, raw-gray-audit.json, fresh-probes.json and second-pass.json; manic-1-r2/evidence/gray-trns-audit.jsonl, gray-resize-audit.jsonl and trns-fixture-acceptance.log.
- **Command**: ncly image shrink <Gray+tRNS PNG>, including orientation and resizing
- **Consequence**: The command successfully publishes a sharing copy with every originally transparent pixel made opaque, violating the unchanged transparency promise. The current independent reader detects this product damage.
- **Owner**: Cook 1.1, native PNG preparation and Gray encoder
- **Fix**: Choose Gray-alpha encoding from decoded transparency or validated tRNS presence as well as original type 4. Preserve native precision and GRAY ICC compatibility. Add generated success/refusal controls across PNG transparency representations, supported depths, orientation and resize, with complete pixels and profiles. Retain the final ICC guard.

### R2-F2 Empty declared ICC profiles bypass malformed-input refusal and reader validation

- **Severity**: `block`
- **State**: `open`
- **Spec**: `input-invalid`
- **Family**: Compatible sample models and profiles
- **Evidence**: Consolidates M1-R2-F2 and M2-R2-F1, originally warning in barrage terms and block under the broken Nicely contract. imageinfo.go:139 skips profileMatches when assembled length is zero; the declared minimum at declarations.go:132 never runs. PNG zlib streams expanding to zero, empty numbered JPEG ICC_PROFILE payloads and WebP ICCP chunks pass preview and run and then disappear. At 27536df fresh PNG shrinks 8532 to 145 bytes, JPEG 7370 to 497, and the manic-1 WebP seed 7572 to 174; all readers succeed. Independent container/CRC/zlib audits find one present zero-length declaration before and none after. The distinct manic-2 WebP has 6704 input bytes from different XMP padding, also producing 174; these observations are not interchangeable. Changing only the PNG profile body to one byte causes IMAGE_SHRINK_INPUT in preview/run. Reader metadataOf and acceptance.txtar:1130 repeat empty-versus-absent equality. See readiness-r2/evidence/empty-profile-container-audit.json, fresh-probes.json and one-byte-* logs; manic-2-r2/evidence/empty-profile-webp-and-audit.json, second-pass-results.json and actual just-empty-profile-acceptance.log.
- **Command**: ncly image shrink <PNG, JPEG or WebP with an explicitly empty ICC declaration>
- **Consequence**: Malformed declared profiles receive a successful preview and publication despite the promised pre-worker refusal. Independent acceptance falsely certifies this boundary. No new private-text leak or valid-profile color-conversion failure is claimed.
- **Owner**: Cook 1.1, product container/profile boundary and independent reader
- **Fix**: Track profile presence separately from assembled byte length and validate every present ICC, including zero bytes, before effects. Make the independent reader retain and validate that distinction. Add PNG/JPEG/WebP empty-declaration refusals beside absent-profile and valid compatible-profile successes; retain complete final ICC equality.

## Rounds

| Round | Head | Change review | Change commits | Readiness review | Finish |
|---|---|---|---|---|---|
| 1 | 662798674d8e6581592c430ae36982cb61eb157f | manic-1 delegation 3, node:delegated-task:command%3Amcp%3Aa997ddb2-d3c1-41a1-b56d-ae8286d664cf%3Adelegate-task%3A461d-1.1-manic-1-3; manic-2 delegation 4, node:delegated-task:command%3Amcp%3Aa997ddb2-d3c1-41a1-b56d-ae8286d664cf%3Adelegate-task%3A461d-1.1-manic-2-4; two independent completed GPT-6.1-Sol via Codex sessions, xhigh/default, serving both the first barrage pair and the independent change review | ea2ec9ce52ca6cb2322be75dbae3bb9887c163a0..68ee01e55ca0f08218607ff2ece346353239ddbf | readiness-r1, delegation 5, Cook 1.1 brigade 461d; independent GPT-6.1-Sol via Codex, xhigh/default; wrote none of the implementation | `fail` |
| 2 | 27536df4a96fb802629f749e1a900ca1371ac5c9 | manic-1 delegation 6, node:delegated-task:command%3Amcp%3Aa997ddb2-d3c1-41a1-b56d-ae8286d664cf%3Adelegate-task%3A461d-1.1-manic-1-6; manic-2 delegation 7, node:delegated-task:command%3Amcp%3Aa997ddb2-d3c1-41a1-b56d-ae8286d664cf%3Adelegate-task%3A461d-1.1-manic-2-7; two independent completed GPT-6.1-Sol via Codex sessions, xhigh/default, serving both the escargot second normal pair and nicely-extension round-2 independent change review, not readiness; complete base-to-head review with fresh concentration since original 68ee01e55ca0f08218607ff2ece346353239ddbf | 3f82e36f2916b6e1c8326971b14531924d4ef8cb..6f2c5eeba11c09a812f88b3502cb10028881bfce | readiness-r2, delegation 8, Cook 1.1 brigade 461d; independent GPT-6.1-Sol via Codex, xhigh/default; wrote none of the implementation | `pending` |

## Order of round 1

| Step | Commit |
|---|---|
| Intake | ea2ec9ce52ca6cb2322be75dbae3bb9887c163a0 |
| Spec | 8bc556b6d3d03a253c202ad1dbd3ba7527be87a0 |
| Scenario | ac1c49562770a5028792949c8534923930047445 |
| Handler | 7a0f36ef0e3791fcf4278acc256eaccb7dbf1727 |

## First failing run

- **Commit**: ac1c49562770a5028792949c8534923930047445
- **Command**: `just scenario shrink-first`
- **Exit**: 1
- **Assertion**: ncly image shrink photo.png --dry-run --json
