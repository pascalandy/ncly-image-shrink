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

## Human items

The human gives the UAT verdict and chooses whether a lossy sharing copy looks acceptable. This local command needs no key, grant or network. The human approves installation of a missing image program.

## Known limits

The spec lists accepted limits. A dry run leaves byte saving pending because it starts no image-writing operation.

## Test-kit accommodations

The binary example fixture ends in an existing newline inside its valid XMP chunk, so [Nicely #496](https://github.com/pascalandy/nicely/issues/496) adds no bytes. Remove that constraint and its transport test after the SDK preserves binary fixtures. The same fixture must continue to pass the actual shrink and independent image checks.

For [Nicely #497](https://github.com/pascalandy/nicely/issues/497), the unreadable input uses a readable reference, a permission assertion, restoration and a complete byte comparison instead of the kit's input snapshot. Restore `protect unreadable.png` and remove those manual steps when the kit supports that case. Other inputs use `protect`.

## Gotchas

The SDK moves HOME and XDG folders during scratch runs. `acceptance.sh` resolves program paths before that move. The independent reader lives in the scenario archive and is extracted for the real acceptance run.
