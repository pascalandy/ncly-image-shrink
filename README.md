# ncly-image-shrink

Make one JPEG, PNG or static WebP smaller for sharing. The copy stays beside the original as `<stem>.shrunk.<ext>`. The original stays unchanged.

The longest side fits 2048 pixels by default, without enlargement. JPEG and WebP use lossy quality 80. PNG uses lossless encoding. The copy applies orientation, retains its color profile and removes location, author, camera and capture date metadata. A candidate that saves no bytes is refused without publishing an output.

```sh
ncly image shrink photo.jpg --dry-run
ncly image shrink photo.jpg
```

Use `ncly describe image shrink --json` for flags, limits and result fields. A dry run leaves byte saving and publication pending. Inspect an existing copy before choosing `--force`.

## Install

Install [Nicely](https://github.com/pascalandy/nicely#install), Go, just and libvips 8.18.7 or later. On macOS use `brew install vips`. On Omarchy use `sudo pacman -S libvips`. `ncly doctor image shrink --json` reports the required program and its install command after installation.

Replace TAG with the highest release tag. `git ls-remote --tags --refs --sort=-v:refname https://github.com/pascalandy/ncly-image-shrink` lists the tags, highest first.

```sh
git clone --branch TAG https://github.com/pascalandy/ncly-image-shrink ncly-image-shrink
cd ncly-image-shrink
just install
ncly describe image shrink --json
```

An annotated tag can produce Git's warning that the tag `is not a commit!`, followed by a detached HEAD notice. The clone still holds the tagged commit.

`just install` replaces a previous install and grants no key. `just uninstall` removes it. This command uses no network or paid service.

## Let your agent find it

Add this line to your agent instructions, such as AGENTS.md or CLAUDE.md:

```text
When asked to use nicely, run ncly skill view nicely
```

The installed `image-shrink` skill teaches when to use the command. [AGENTS.md](AGENTS.md) gives the development loop, checks and release procedure.
