Commit: 662798674d8e6581592c430ae36982cb61eb157f
Reviewer: readiness-r1, delegation 5, Cook 1.1 brigade 461d; independent GPT-6.1-Sol via Codex, xhigh/default; wrote none of the implementation
Verdict: Needs work

# Readiness review

`just records` writes this file from `review.toml`. Change that file, then run `just records` again.

## Findings

### R1-F1 Unvalidated PNG display declarations publish private text

- **Severity**: `block`
- **State**: `open`
- **Spec**: `privacy`
- **Family**: Descriptive metadata
- **Evidence**: At 662798674d8e6581592c430ae36982cb61eb157f, imageinfo.go:56 accepts gAMA/cHRM/sRGB/cICP payloads opaquely; png.go:85 and imageinfo.go:295 copy them. Fresh generated CRC-correct inputs each contain the complete seeded author/location/camera/date block in a 194-byte declaration instead of the required 4/32/1/4 bytes. All four calls exit 0 and publish 574-byte copies from 8953-byte inputs with the full block retained. The fixture reader at acceptance.txtar:786 and :873 accepts all four. Actual just acceptance with malformed-gAMA.png exits 0 and reports privacy success. Ordinary PNG control drops the block. See evidence/readiness-r1/output/reproduced.json, second-pass-artifacts.json, just-malformed-fixture-acceptance.log, and probes/malformed-*.png. Consolidates manic-2 F1/F4; fresh output hashes match its preserved artifacts.
- **Command**: ncly image shrink <file.png>, including supplied-fixture acceptance
- **Consequence**: A sharing copy discloses private author, location, camera and capture date while both the command and acceptance report success. This breaks the privacy and validated-display-metadata contract for untrusted local images.
- **Owner**: Cook 1.1, extension parser and independent acceptance reader
- **Fix**: Validate exact lengths, values, uniqueness and conflicts across the complete PNG display-declaration family before use, and serialize only validated fields. Independently validate retained declarations and reject private-bearing bad variants in the reader. Add generated successful and refused pairs for the whole declaration family, including malformed and conflicting declarations, with primary-output privacy checks. Cover applicable Adobe/CMYK declarations with R1-F2. Keep the existing privacy promise.

### R1-F2 Pixel-model conversion makes valid Gray and CMYK profiles incompatible

- **Severity**: `block`
- **State**: `open`
- **Spec**: `profile`
- **Family**: Orientation and color
- **Evidence**: png.go:26 and :57 allocate RGB models for oriented/resized Gray PNGs; shrink.go:214 uses thumbnail for JPEG; shrink.go:274 correctly refuses ICC loss. A fresh 10894-byte valid Gray/GRAY-ICC control succeeds unchanged, producing 380 bytes. Its otherwise identical orientation-2 variant fails with IMAGE_SHRINK_RESULT, exit 1; the orientation-1 input also fails with max-side 32. Second-pass chunk comparison confirms the pair differs only in EXIF orientation. A valid 188626-byte CMYK/CMYK-ICC JPEG passes dry-run and fails the real call with the same code. Preserved direct-vips controls and v8.18.7 thumbnail.c:792-807/foreign.c:1771-1786 corroborate the incompatible-model cause. See evidence/readiness-r1/output/reproduced.json, second-pass-artifacts.json, probes/gray-*.png, probes/cmyk.jpg, and evidence/manic-1/color-model-control.log. Consolidates manic-1 F1 and manic-2 F2; retains the highest original critical finding as Nicely block.
- **Command**: ncly image shrink <Gray PNG or CMYK JPEG>, with orientation or resizing as applicable
- **Consequence**: Supported valid images cannot obtain the promised smaller sharing copy. The final guard prevents publishing a profile-damaged image, but no user accepted a Gray or CMYK exclusion.
- **Owner**: Cook 1.1, PNG preparation and JPEG processing
- **Fix**: Preserve the relationship between sample model and profile during orientation, resizing and encoding. Keep the final complete ICC comparison. Add independent complete-content/profile checks and successful/refused pairs for the full orientation/color declaration family, including Gray, Gray plus alpha, CMYK, Adobe interpretation, compatible profiles and invalid/incompatible declarations. Do not narrow support to RGB or remove the guard.

### R1-F3 The fidelity reader skips pixels and accepts gross corruption

- **Severity**: `block`
- **State**: `open`
- **Spec**: `pixels`
- **Family**: Orientation and color
- **Evidence**: testdata/script/acceptance.txtar:611-614 skips broad boundary and central bands for lossy formats and resize checks. Freshly changing 1008 of 2048 pixels in the valid 64x32 PNG to opaque black, while preserving its valid ICC and gamma and omitting private metadata, still makes imageproof check 1 border-bad.png 64 32 exit 0 with 'complete pixels, orientation, alpha, ICC and privacy match'. The unmodified control passes; the same corrupt image without resize arguments fails at pixel 0,0, exit 2. A separate full-image Go comparison confirms exactly 1008 changed pixels. See evidence/readiness-r1/output/reproduced.json, second-pass-pixels.json, probes/photo.shrunk.png and probes/border-bad.png. Consolidates manic-1 F2 and the fidelity-control portion of manic-2 F3.
- **Command**: Independent acceptance checks for ncly image shrink resized PNG, JPEG and WebP outputs
- **Consequence**: Required independent full-content proof is missing and can certify a badly corrupted sharing copy. This finding establishes an acceptance defect, not a demonstrated encoder corruption.
- **Owner**: Cook 1.1, independent acceptance reader and references
- **Fix**: Compare every output pixel with an independently derived orientation/resampling reference and justified format-specific tolerances, including all boundaries and alpha. Add fidelity-specific bad variants with valid privacy/profile metadata that fail because their pixels or display declarations are wrong. Keep the full-content promise.

### R1-F4 Multipart classifiers lack complete generated acceptance pairs

- **Severity**: `block`
- **State**: `open`
- **Spec**: `multipart`
- **Family**: One image, without animation
- **Evidence**: The complete committed scenario scope at the readiness head contains only acceptance.txtar:137-139 and :743-746 for multipart refusal, adding an APNG acTL declaration. No committed stimulus generates animated WebP or multipicture JPEG. imageinfo.go:71 and :101-108 classify MPF and WebP animation separately; jpegChunks also classifies concatenated images. Both independent reports identify these omissions. Their successful MPF/VP8X marker probes are preserved evidence of those probes, not genuine animation/multipicture fixture pairs or committed regression coverage. See spec.md Case matrix, intake.toml display-declaration fact, the complete testdata/script files, and evidence/manic-2/output/variant-results.json. Consolidates multipart portions of manic-1 F3 and manic-2 F3.
- **Command**: ncly image shrink <JPEG, PNG or WebP>
- **Consequence**: The mandatory round-1 whole-family contract is unmet. The existing suite can pass if unexercised multipart classifiers regress; no observed frame-loss defect is claimed here.
- **Owner**: Cook 1.1, generated fixtures and case matrix
- **Fix**: Commit generated single/static and genuine multiple/animated-image acceptance/refusal pairs for JPEG, PNG and WebP. Cover applicable animation flag/chunk branches, MPF and concatenated JPEG classification, and independently read successful still-image outputs. Keep multipart open until the full family is proved.

### R1-F5 Lossy quality overrides have no successful independent proof

- **Severity**: `block`
- **State**: `open`
- **Spec**: `quality`
- **Family**: Valid command options
- **Evidence**: Complete inspection of testdata/script/*.txtar and acceptance.sh finds only default-quality successful JPEG/WebP calls. acceptance.txtar:80 uses quality 1 on lossless PNG; :143 refuses quality 101. Neither proves an accepted nondefault lossy quality. shrink.go:204-207 passes Q, but inspecting that argument is no primary-output acceptance proof. See the quality claims at acceptance.txtar:8 and :87 and the preserved manic-1 F3/manic-2 F3 reports.
- **Command**: ncly image shrink <JPEG or WebP> --quality <1..100>
- **Consequence**: The selected-quality promise lacks its required test. Hardcoding JPEG/WebP to Q=80 could pass all existing quality stimuli; no actual quality-option implementation failure is claimed.
- **Owner**: Cook 1.1, lossy-quality acceptance family
- **Fix**: Add generated JPEG and WebP public-boundary successful calls using materially different valid quality requests. Independently decode and judge their primary outputs so that ignoring an override fails for the intended reason, alongside the invalid-option refusal. Retain default 80 and lossless PNG checks.

### R1-F6 Post-publication signals lack the required survival proof

- **Severity**: `block`
- **State**: `open`
- **Spec**: `publication`
- **Family**: Publication errors
- **Evidence**: The pinned .agents/skills/nicely-extension/replace-a-file.md, 'Fail before or after the commit', explicitly requires proving that a signal after the commit keeps the completed output with survives and signal. At the readiness head, both signal cases in testdata/script/acceptance.txtar use STUB=hang, which writes only the staged candidate and its marker before vips returns. The output is protected as absent or old. No scenario uses survives or synchronizes a signal after link/rename. A fresh separate cleanup-failure probe confirms the create path can keep its valid 368-byte completed image and return its output path with IMAGE_SHRINK_PUBLISH, exit 1; that successful scratch observation does not prove the required post-publication signal case. See shrink.go:175-179 and :244-258, evidence/readiness-r1/output/after-publication.*, and second-pass-artifacts.json. This is a new readiness finding. The publication spec ID is still under Acceptance at this frozen head and needs the Open-lines row identified in report.md.
- **Command**: ncly image shrink <file>, create and force-replace publication lifecycle
- **Consequence**: The required acceptance proof for retaining and reporting a completed output after a signal is absent. No post-publication loss or cleanup-error reporting defect was reproduced.
- **Owner**: Cook 1.1, publication lifecycle acceptance family
- **Fix**: Keep the publication promise open and add reliable public-boundary post-commit INT/TERM scenarios using survives and signal, with complete-output bytes and permissions, destination/error reporting, staging cleanup and child-exit checks. Cover create and replacement alongside the existing pre-publication failure/signal controls. Commit an after-publication cleanup-failure control as part of that family. No new product limit is accepted.

## Rounds

| Round | Head | Change review | Change commits | Readiness review | Finish |
|---|---|---|---|---|---|
| 1 | 662798674d8e6581592c430ae36982cb61eb157f | manic-1 delegation 3, node:delegated-task:command%3Amcp%3Aa997ddb2-d3c1-41a1-b56d-ae8286d664cf%3Adelegate-task%3A461d-1.1-manic-1-3; manic-2 delegation 4, node:delegated-task:command%3Amcp%3Aa997ddb2-d3c1-41a1-b56d-ae8286d664cf%3Adelegate-task%3A461d-1.1-manic-2-4; two independent completed GPT-6.1-Sol via Codex sessions, xhigh/default, serving both the first barrage pair and the independent change review | ea2ec9ce52ca6cb2322be75dbae3bb9887c163a0..68ee01e55ca0f08218607ff2ece346353239ddbf | readiness-r1, delegation 5, Cook 1.1 brigade 461d; independent GPT-6.1-Sol via Codex, xhigh/default; wrote none of the implementation | `pending` |

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
