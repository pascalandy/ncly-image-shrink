Commit: f691b68d3bc619f38e62fb09c5685260bea1eb42
Reviewer: readiness-r3, delegation 11, Cook 1.1 brigade 461d; independent GPT-6.1-Sol via Codex, xhigh/default; wrote none of the implementation
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
- **Evidence**: Historical readiness-r2 evidence, retained with its original session-local wording: Preserves manic-2 F1/F4 attribution and the original block grade. a13b15e validates retained PNG display declarations before input use and candidate allowlisting; the preceding failing scenario supplies length/value/duplicate/conflict and privacy-reader controls. Subsequent d503d24/60fe121 complete declaration/model/placement handling. At frozen 27536df, all four preserved private-bearing gAMA/cHRM/sRGB/cICP seeds exit 2 with IMAGE_SHRINK_INPUT and no output. Fresh generated valid declarations succeed and pass independent pixels/display/privacy checks; length/value/duplicate variants refuse, and a readable malicious declaration fails the reader. Full committed scenario pairs include conflicts and legal post-IDAT EXIF. See readiness-r2/evidence/fresh-probes.json, second-pass.json and sources; both independent round-2 full-suite logs at 6f2 corroborate family execution. No new private-text leak was established by the empty ICC defect. Round-3 regression review at f691b68 retains this fixed state: fresh valid gAMA/cHRM/sRGB/cICP copies pass full readers and all twelve length/value/duplicate variants refuse. The extended transparency/profile admission stays fail closed. No new declaration privacy leak is demonstrated.
- **Command**: ncly image shrink <file.png>, including supplied-fixture acceptance
- **Consequence**: The original private-bearing PNG declaration leak and false certification no longer reproduce. Complete ICC retention and malformed ICC presence are assessed separately under R1-F2/R2-F2.
- **Owner**: Cook 1.1, extension parser and independent acceptance reader
- **Fix**: Implemented declaration validation and whole-family success/refusal/privacy controls. Retain them and the original privacy promise.

### R1-F2 Native sample-model and profile family still lacks complete proof

- **Severity**: `block`
- **State**: `fixed`
- **Spec**: `profile`
- **Family**: Orientation and color
- **Commit**: cc5725ee22829f52ffcedb0b83be665c4f830657
- **Evidence**: Preserves original manic-1 F1/manic-2 F2 attribution and block grade. d503d24 repaired native Gray/Gray-alpha and CMYK handling; 0faa261 preserves decoded Gray+tRNS alpha, and cc5725e completes declared-profile admission and whole-family controls. The final delta sessions independently replay the original model controls and repaired alpha/empty-ICC counterexamples. At frozen f691b68 this readiness session freshly replays preserved oriented/resized Gray and CMYK controls, all eight original Gray 8/16-bit type-0/type-4 orientation-1/2 cases with and without resize, and eleven generated Gray/RGB/palette tRNS representations with orientation 6 and resize. Full reader checks pass. Independent raw CRC/zlib/sample audits count 1024 transparent and 1024 opaque input samples in every generated representation and verify complete source/output ICC equality. Committed accepted/refused pairs cover Gray/Gray-alpha 8/16, Gray JPEG, RGB/CMYK, Adobe/JFIF, native precision, orientation, resize, malformed/incompatible profiles, empty presence, JPEG fragment assembly and WebP flag consistency. Final candidate validation still requires profile presence and every profile byte to match. Earlier LittleCMS audits establish the generated GRAY/CMYK profile controls are usable. The separate odd-factor reference false rejection remains R1-F3; no surviving model/profile counterexample is demonstrated. See readiness-r3/evidence/fresh-probes.json, raw-transparency-audit.json, second-pass.json and sources; both final delta raw audits remain separately attributed.
- **Command**: ncly image shrink <Gray PNG or CMYK JPEG>, with orientation or resizing as applicable
- **Consequence**: The original native-model refusals and the two concrete whole-family closure objections from round 2 are repaired and paired. This closes the broad model/profile finding without duplicating R1-F3 or narrowing the supported input promise.
- **Owner**: Cook 1.1, PNG preparation and JPEG processing
- **Fix**: Implemented native sample-model preservation, decoded Gray alpha and validation of every declared profile, with complete accepted/refused family pairs. Retain native precision, compatible ICC models and full final ICC equality.

### R1-F3 The fidelity reader rejects faithful odd-factor reductions

- **Severity**: `block`
- **State**: `open`
- **Spec**: `pixels`
- **Family**: Orientation and color
- **Evidence**: Preserves original manic-1 F2 attribution, related manic-2 fidelity evidence, M2-R2-F2 aliasing history and the existing block grade. 18247d4 now integrates every source pixel with fractional area in premultiplied RGBA and removes the earlier four-point aliasing. The old black/transparent substitutions no longer certify. M1-DELTA-F1 is a residual of this same finding, not a new ID or a severity reduction. At f691b68:testdata/script/acceptance.txtar:1489-1531 the box-relative mean limits reject faithful CatmullRom PNG reductions. Fresh default-size 14336x128 16-bit alpha input shrinks 34878 to 948 bytes at 2048x18; the current reader exits 2 at color mean difference 18.262/255 against 8/255. A separately authored exact-rational cubic polynomial predicts all 36864 output RGBA pixels with maximum difference zero in every channel, again after force. Independent complete-pixel audit measures box-relative mean absolute difference 18.262459721546694/255 in every channel, also exceeding resized-alpha 4/255. The 133x67 orientation-8 side-19 alpha and opaque cases fail at 17.619/255; matched 6:1 and 8:1 controls pass, while black/transparent substitutes still fail. Both final delta reports inspected the complete-area loops; manic-2 retained no new finding in its probes but did not exercise this odd-factor witness. This is acceptance false rejection, with no demonstrated product pixel damage or surviving total-corruption certification. See readiness-r3/evidence/fresh-probes.json, second-pass.json, raw-transparency-audit.json and probes/second_pass.py; manic-1-delta/evidence/independent-kernel-audit.json, its repeat, fraction-followup-summary.json and default-seven-no-override-reader.log.
- **Command**: Independent acceptance checks for ncly image shrink resized PNG, JPEG and WebP outputs
- **Consequence**: Mandatory independent primary-output validation rejects a faithful supported sharing copy, including at the default size. Supplied-fixture acceptance would fail at that reader. The unchanged complete-content promise remains open even though the earlier skipped-region and aliasing defects are eliminated.
- **Owner**: Cook 1.1, independent acceptance reader and references
- **Fix**: Add odd/even/fractional high-frequency accepted/refused pairs and a justified reference or error envelope covering permitted kernel differences. Keep complete source/destination coverage, orientation, native depth, premultiplied alpha, aggregate guards and matched valid-metadata corrupt-output controls. Do not delete the aggregate guard or fit one constant to this fixture. Keep pixels open.

### R1-F4 Multipart classifiers lack complete generated acceptance pairs

- **Severity**: `block`
- **State**: `fixed`
- **Spec**: `multipart`
- **Family**: One image, without animation
- **Commit**: 2cb63ea792272ca441509279fa8b3f049c4a0493
- **Evidence**: Historical readiness-r2 evidence, retained with its original session-local wording: Preserves both original reviewers attribution and block grade. The committed generator now makes actual APNG fcTL/fdAT frames, WebP ANIM/ANMF frame payloads, flag-only/chunk-only branches, MPO MP-entry targets and concatenated JPEGs, alongside still-image success cases. At frozen 27536df this session freshly generates the fixtures, ffprobe decodes two APNG frames, ImageMagick decodes two WebP frames, and all seven public direct-CLI multipart probes refuse with IMAGE_SHRINK_MULTIPART and absent output. Both independent round-2 suites at 6f2 pass the matching host-boundary families. Their independent TIFF audit confirms sizes 7928/7838 and actual JPEG targets at offsets 0/7928, with the second TIFF-relative offset 7918. See readiness-r2/evidence/second-pass.json and sources; manic-2-r2/evidence/second-pass-results.json and families-results.json. Round-3 regression inspection retains the real generated frame/MP-entry families and still-image pairs. Multipart classifiers remain fail closed, including changed WebP header handling. The final delta full-suite logs at f5546aa and their raw audits disclose no regression. Current f691b68 check fails on reopened pixels claims and supplies no green gate.
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
- **Evidence**: Historical readiness-r2 evidence, retained with its original session-local wording: Preserves both original reviewers attribution and block grade. Committed JPEG/WebP public calls at quality 20 and 95 independently compare every decoded pixel, profile/privacy and sizes. Fresh frozen-head direct CLI executions repeat JPEG MSE 984.016 versus 15.180 and WebP 1161.724 versus 868.137; higher quality is larger and closer. Identical-output ignored-override controls fail for both formats. Invalid options, default 80 and lossless PNG controls remain in the passing 6f2 suites. See readiness-r2/evidence/second-pass.json and second-quality-* logs; manic-2-r2/evidence/families-results.json. These are deterministic fixture metrics, not general visual-quality benchmarks. Round-3 regression inspection retains successful JPEG/WebP quality 20/95 complete decoded comparisons, ignored-override negatives, invalid-option and PNG controls. Quality parsing and Q forwarding are unchanged. Final delta full-suite logs at f5546aa disclose no regression; this is not a new visual-quality benchmark.
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
- **Evidence**: Historical readiness-r2 evidence, retained with its original session-local wording: Preserves readiness-r1 delegation 5 attribution and block grade. acceptance.txtar:458 onward commits create/INT and replace/TERM using survives and signal. The independent published observer fully reads the new, strictly smaller destination before writing the marker; the old replacement copy cannot satisfy that condition. Ordinary scratch cleanup provides the interval without a production hook; a missed interval makes the test fail. Input bytes/mode, completed destination bytes/mode, output path/error, stage cleanup and child exit are checked. A committed forced cleanup-failure case retains and reads the completed output with IMAGE_SHRINK_PUBLISH and its destination. Both independent fresh 6f2 full suites execute these cases successfully, matching Cook round2-postpublication.log and repeat log. Source and raw logs were independently inspected here; this session did not replay the installed host lifecycle suite. Native preparation signal measurements are separate and do not prove final-validation cancellation. Round-3 regression inspection retains the independent completed-destination observer, survives with create/INT and replace/TERM, permissions/path/error/cleanup/child assertions and real cleanup-error control. Publication implementation is unchanged. Final delta full-suite logs at f5546aa disclose no regression. Native final-validation cancellation remains unconfirmed and is not inferred from preparation measurements.
- **Command**: ncly image shrink <file>, create and force-replace publication lifecycle
- **Consequence**: The explicit post-publication acceptance gap is closed by committed real publication observations and survival checks. No signal or cleanup loss was reproduced.
- **Owner**: Cook 1.1, publication lifecycle acceptance family
- **Fix**: Implemented the required post-commit survival proof and cleanup-failure control. Retain the pre-commit controls and original publication policy.

### R2-F1 Gray PNG tRNS transparency is discarded during native-model preparation

- **Severity**: `block`
- **State**: `fixed`
- **Spec**: `alpha`
- **Family**: Orientation and color
- **Commit**: 0faa261b5b1d6bf6963c82df2440fd86f454d4e5
- **Evidence**: Preserves M1-R2-F1 critical attribution and existing Nicely block grade. At f691b68:png.go:70-76 decoded nonopacity enables Gray-alpha while retaining the original Gray model and native precision. Both independent final delta sessions replay all preserved 8/16-bit orientation-1/2 type-0+tRNS inputs and explicit type-4 controls with resize. This readiness session freshly reruns all eight original combinations, unchanged-size and side 32. The original 1024 transparent and 1024 opaque samples survive; resized boundary alpha remains. Eleven fresh Gray 1/2/4/8/16, RGB 8/16 and palette 1/2/4/8 representations at orientation 6 and side 31 pass the complete reader. Separate CRC/zlib/raw-sample audits verify actual transparency and complete profile equality, with 16-bit outputs remaining wide and Gray outputs using model 4. Malformed Gray/RGB/palette tRNS pairs refuse preview and run. An otherwise valid opaque Gray substitute fails specifically at channel 3 alpha difference 65535, bound 5140, in the resized control. The committed no-resize opaque control fails at alpha bound 8 in manic-2-delta/evidence/opaque-control.log. See readiness-r3/evidence/fresh-probes.json, raw-transparency-audit.json and both delta raw alpha/family audits.
- **Command**: ncly image shrink <Gray+tRNS PNG>, including orientation and resizing
- **Consequence**: The original successful publication with all Gray+tRNS transparency removed no longer reproduces. Complete family pairs distinguish transparency preservation from opaque output, without relying on an incompatible-profile refusal.
- **Owner**: Cook 1.1, native PNG preparation and Gray encoder
- **Fix**: Implemented decoded-opacity Gray-alpha selection with native precision and compatible profiles. Retain whole-family generated transparency and malformed-declaration pairs, explicit-alpha controls and final ICC equality.

### R2-F2 Empty declared ICC profiles bypass malformed-input refusal and reader validation

- **Severity**: `block`
- **State**: `fixed`
- **Spec**: `input-invalid`
- **Family**: Compatible sample models and profiles
- **Commit**: cc5725ee22829f52ffcedb0b83be665c4f830657
- **Evidence**: Preserves shared M1-R2-F2/M2-R2-F1 warning attribution and existing Nicely block grade. At f691b68:imageinfo.go:151 every present ICC reaches structural/model validation, independent of payload length. Presence and full bytes remain in final candidate equality at shrink.go:280. Both independent delta sessions replay the original empty PNG/JPEG/WebP declarations and inspect valid/absent-profile successes and malformed envelopes. This readiness session freshly replays each preserved empty declaration in preview/run: all exit 2 with IMAGE_SHRINK_INPUT, absent output and unchanged input bytes/mode. The current independent reader rejects each specifically as invalid declared ICC even against itself. Fresh legal ordered/reversed JPEG profile fragments succeed with full readers, while all generated PNG name/method/duplicate, JPEG missing/duplicate/inconsistent/zero/short and WebP flag/chunk/header refusal variants fail preview/run. Source inspection verifies sequence assembly, presence distinct from length and complete final equality. Independent manic-2 sentinel probes show these refusals precede vips, with a live legal-preview control. See readiness-r3/evidence/fresh-probes.json and sources; manic-2-delta/evidence/before-worker.json, generated-family-audit.json and probes-results.json. No new private-text leak, valid-profile conversion loss or full ICC semantic-conformance claim is made.
- **Command**: ncly image shrink <PNG, JPEG or WebP with an explicitly empty ICC declaration>
- **Consequence**: Explicitly empty ICC declarations can no longer pass malformed-input admission or certify themselves as absent profiles. The successful no-profile and compatible-profile paths remain supported.
- **Owner**: Cook 1.1, product container/profile boundary and independent reader
- **Fix**: Implemented separate profile presence, validation of every declared body and independent reader validation, with PNG/JPEG/WebP accepted/refused envelope families. Retain complete final ICC equality.

## Rounds

| Round | Head | Change review | Change commits | Readiness review | Finish |
|---|---|---|---|---|---|
| 1 | 662798674d8e6581592c430ae36982cb61eb157f | manic-1 delegation 3, node:delegated-task:command%3Amcp%3Aa997ddb2-d3c1-41a1-b56d-ae8286d664cf%3Adelegate-task%3A461d-1.1-manic-1-3; manic-2 delegation 4, node:delegated-task:command%3Amcp%3Aa997ddb2-d3c1-41a1-b56d-ae8286d664cf%3Adelegate-task%3A461d-1.1-manic-2-4; two independent completed GPT-6.1-Sol via Codex sessions, xhigh/default, serving both the first barrage pair and the independent change review | ea2ec9ce52ca6cb2322be75dbae3bb9887c163a0..68ee01e55ca0f08218607ff2ece346353239ddbf | readiness-r1, delegation 5, Cook 1.1 brigade 461d; independent GPT-6.1-Sol via Codex, xhigh/default; wrote none of the implementation | `fail` |
| 2 | 27536df4a96fb802629f749e1a900ca1371ac5c9 | manic-1 delegation 6, node:delegated-task:command%3Amcp%3Aa997ddb2-d3c1-41a1-b56d-ae8286d664cf%3Adelegate-task%3A461d-1.1-manic-1-6; manic-2 delegation 7, node:delegated-task:command%3Amcp%3Aa997ddb2-d3c1-41a1-b56d-ae8286d664cf%3Adelegate-task%3A461d-1.1-manic-2-7; two independent completed GPT-6.1-Sol via Codex sessions, xhigh/default, serving both the escargot second normal pair and nicely-extension round-2 independent change review, not readiness; complete base-to-head review with fresh concentration since original 68ee01e55ca0f08218607ff2ece346353239ddbf | 3f82e36f2916b6e1c8326971b14531924d4ef8cb..6f2c5eeba11c09a812f88b3502cb10028881bfce | readiness-r2, delegation 8, Cook 1.1 brigade 461d; independent GPT-6.1-Sol via Codex, xhigh/default; wrote none of the implementation | `fail` |
| 3 | f691b68d3bc619f38e62fb09c5685260bea1eb42 | manic-1 delegation 9, node:delegated-task:command%3Amcp%3Aa997ddb2-d3c1-41a1-b56d-ae8286d664cf%3Adelegate-task%3A461d-1.1-manic-1-9; manic-2 delegation 10, node:delegated-task:command%3Amcp%3Aa997ddb2-d3c1-41a1-b56d-ae8286d664cf%3Adelegate-task%3A461d-1.1-manic-2-10; two separate independent completed GPT-6.1-Sol via Codex sessions, xhigh/default, serving the required final delta and round-3 independent change review, not readiness; both cover 6f2c5ee..f5546aa and direct effects; supplementary just check passes included SDK-owned private scratch installs contrary to the old literal no-install brief, with no persistent/user install or shared-state change; standalone probes installed nothing | 5901656a332e682d84fb1c11d0237a4900c8238a..f5546aa365d2344ac46fb53cc11f7086bb08ad63 | readiness-r3, delegation 11, Cook 1.1 brigade 461d; independent GPT-6.1-Sol via Codex, xhigh/default; wrote none of the implementation | `pending` |

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
