#!/usr/bin/env bash
# Run the sharing-copy command in the SDK's owned scratch install. All default
# images and the independent reader come from the acceptance scenario archive.
set -Eeuo pipefail

fct_die() { printf 'acceptance: %s\n' "$*" >&2; exit 1; }
fct_sum() {
	if command -v sha256sum >/dev/null; then sha256sum "$1"; else shasum -a 256 "$1"; fi | cut -d ' ' -f 1
}
fixture="${1:-}"
if [[ -n "$fixture" ]]; then
	[[ "$fixture" == /* ]] || fixture="$PWD/$fixture"
	[[ -f "$fixture" ]] || fct_die "fixture is not a regular file"
fi
cd "$(dirname "${BASH_SOURCE[0]}")"
scratch="$(mktemp -d "${TMPDIR:-/tmp}/ncly-image-shrink-acceptance.XXXXXX")"
trap 'rm -rf -- "$scratch"' EXIT
mkdir -p "$scratch/bin" "$scratch/work"
path="$scratch/bin"
for program in go just vips; do
	found="$(mise which "$program" 2>/dev/null || command -v "$program")" || fct_die "$program is not on PATH"
	path="$path:${found%/*}"
done
export PATH="$path:$PATH"
sed '1,/^-- proof.go --$/d' testdata/script/acceptance.txtar > "$scratch/proof.go"
go build -o "$scratch/bin/imageproof" "$scratch/proof.go"
go build -o "$scratch/bin/ncly" github.com/pascalandy/nicely/cmd/ncly
XDG_DATA_HOME="$scratch/data" just install
export XDG_DATA_HOME="$scratch/data" XDG_CONFIG_HOME="$scratch/config" XDG_STATE_HOME="$scratch/state" XDG_CACHE_HOME="$scratch/cache"
cd "$scratch/work"
imageproof make
input="$PWD/photo.png"
if [[ -n "$fixture" ]]; then
	before="$(fct_sum "$fixture")"
	cp -- "$fixture" "$PWD/"
	input="$PWD/${fixture##*/}"
fi
input_before="$(fct_sum "$input")"
# spec: preview
ncly image shrink "$input" --dry-run --json --no-input
output="${input%.*}.shrunk.${input##*.}"
[[ ! -e "$output" ]] || fct_die "dry run published an output"
# spec: smaller pixels format privacy profile png-color alpha publication
ncly image shrink "$input" --json --no-input
[[ -f "$output" ]] || fct_die "the output is absent"
[[ "$(wc -c < "$output")" -lt "$(wc -c < "$input")" ]] || fct_die "the output saves no bytes"
imageproof check 1 "$output"
[[ "$(fct_sum "$input")" == "$input_before" ]] || fct_die "the input changed"
if [[ -n "$fixture" ]]; then
	[[ "$(fct_sum "$fixture")" == "$before" ]] || fct_die "the fixture changed"
fi
# spec: not-smaller
status=0
ncly image shrink tiny.png --json --no-input > "$scratch/refusal.out" 2> "$scratch/refusal.err" || status=$?
[[ "$status" == 2 ]] || fct_die "the optimized PNG was not refused"
grep -q '"code":"IMAGE_SHRINK_NOT_SMALLER"' "$scratch/refusal.err" || fct_die "wrong refusal code"
[[ ! -e tiny.shrunk.png ]] || fct_die "the refused PNG left an output"
# spec: replace
ncly image shrink "$input" --force --json --no-input
imageproof check 1 "$output"
imageproof clean
cd - > /dev/null
just uninstall > /dev/null
status=0
ncly describe image shrink --json > /dev/null 2> "$scratch/uninstall.err" || status=$?
if [[ "$status" != 2 ]] || ! grep -q '"code":"NOT_FOUND"' "$scratch/uninstall.err"; then
	fct_die "uninstall left the extension discoverable"
fi
printf 'acceptance: sharing copy, privacy, fidelity, refusal, replacement and uninstall passed\n'
