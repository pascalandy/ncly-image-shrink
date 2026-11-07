# ncly-image-shrink

`ncly-image-shrink` is an extension of [Nicely](https://github.com/pascalandy/nicely), which adds `ncly image shrink`.
`vision.md` states its job, and `spec.md` states what each command proves.

## Install

A release installs with `ncly`, Go, and `just`.
To install `ncly`, follow Install in [Nicely's README](https://github.com/pascalandy/nicely#install).
Replace TAG with the highest release tag. `git ls-remote --tags --refs --sort=-v:refname https://github.com/pascalandy/ncly-image-shrink` lists the tags, highest first.

```sh
git clone --branch TAG https://github.com/pascalandy/ncly-image-shrink ncly-image-shrink
cd ncly-image-shrink
just install
ncly describe image shrink --json
```

For an annotated tag, Git can warn that the tag `is not a commit!`, then note a detached HEAD: the clone holds the tagged commit.
`ncly describe` confirms the install, with `image shrink` and the source `extension`.
`just install` replaces a previous install, and grants no key. `just uninstall` in the clone removes it.
An agent reads `ncly skill view image-shrink` to learn when to use the extension.

## Let your agent find it

Add one line to your agent instructions, such as `AGENTS.md` or `CLAUDE.md`, to teach an agent to use Nicely:

```text
When asked to use nicely, run ncly skill view nicely
```

## Change the extension

`AGENTS.md` gives the reading order, the change steps, the recipes, and the release steps.
