# image-shrink QA

Follow Nicely's [UAT protocol](https://github.com/pascalandy/nicely/blob/v0.27.0/docs/qa/protocol.md). The command's behavior is in [spec.md](spec.md).

## Prerequisites

Run `just acceptance` on the tag under review. The command needs vips 8.18.7 or later. Go and just build the extension and its independent image reader.

## Prompt for the fresh agent

"Use Nicely to shrink the image `<run>/photo.png`"

## Criteria

| Family | Criterion | Covered by |
|---|---|---|
| Result | The published file meets `smaller` and `format` | `just acceptance` and the generated-image scenario |
| Fidelity | The `pixels`, `orientation`, `profile`, `png-color`, `alpha` and `depth` properties hold | The independent Go reader in `acceptance.txtar` |
| Privacy | The `privacy` property holds for the location and descriptive block | `reads -private`, the container reader, and `just acceptance` |
| Determinism | Repeating the same call with force produces the same image content | `just acceptance` compares both results with the independent reference |
| Errors | Each declared failure and the optimized-image refusal is observed | The refusal, budget, program, publication and signal cases in `acceptance.txtar` |
| Uninstall | The installed command disappears from describe | `just acceptance` and the release receipt |

## Error cases

[acceptance.txtar](testdata/script/acceptance.txtar) owns error stimuli and expected answers. [shrink-first.txtar](testdata/script/shrink-first.txtar) retains the first public-boundary preview scenario. Its first failing commit is recorded in review.toml.

The independent reader checks every output pixel against the decoded input. Its box reference integrates every source pixel's fractional overlap in premultiplied RGBA, with orientation and centered JPEG chroma expansion. Per-pixel and whole-image bounds allow codec and resampling-kernel differences. Integer and fractional high-frequency reductions reject black or transparent substitutes. Separate bad outputs drop ICC or change gamma while retaining valid pixels and privacy metadata. Quality cases compare complete decoded JPEG/WebP outputs at 20 and 95. Multipart cases generate APNG/WebP frame sequences, MPO and concatenated JPEGs.

Transparency cases cover Gray tRNS at 1, 2, 4, 8 and 16 bits, RGB tRNS, palette tRNS and explicit alpha, with orientation and resizing. Profile cases distinguish absence from a declared empty payload in every format. JPEG profile fragments may arrive out of order, but missing, repeated or inconsistent parts are refused. PNG envelope errors and WebP profile flag/chunk mismatches are refused too.

Post-publication INT/create and TERM/replace tests use `survives`. An independent observer reads the complete new destination before it makes the signal marker. A normal scratch-file cleanup load provides that interval, without a production hook. The tests also force a cleanup error after publication and check the retained image and reported path.

## Human items

The human gives the UAT verdict and chooses whether a lossy sharing copy looks acceptable. This local command needs no key, grant or network. The human approves installation of a missing image program.

## Known limits

The spec lists accepted limits. A dry run leaves byte saving pending because it starts no image-writing operation.

## Test-kit accommodations

The binary example fixture ends in an existing newline inside its valid XMP chunk, so [Nicely #496](https://github.com/pascalandy/nicely/issues/496) adds no bytes. Remove that constraint and its transport test after the SDK preserves binary fixtures. The same fixture must continue to pass the actual shrink and independent image checks.

For [Nicely #497](https://github.com/pascalandy/nicely/issues/497), the unreadable input uses a readable reference, a permission assertion, restoration and a complete byte comparison instead of the kit's input snapshot. Restore `protect unreadable.png` and remove those manual steps when the kit supports that case. Other inputs use `protect`.

## Gotchas

The SDK moves HOME and XDG folders during scratch runs. `acceptance.sh` resolves program paths before that move. The independent reader lives in the scenario archive and is extracted for the real acceptance run.
