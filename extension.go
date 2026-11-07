// Package image is the image-shrink extension of Nicely. Its one value,
// Extension, holds the manifest, the catalogs, the skill, and
// one handler per command. cmd/ncly-image-shrink runs it, and
// extension_test.go tests it.
package image

import (
	"embed"

	"github.com/pascalandy/nicely/sdk/extension"
)

//go:embed manifest.toml
var manifest []byte

//go:embed locales/*.toml
var locales embed.FS

//go:embed skill
var skill embed.FS

// Extension is the image-shrink extension. Each command of manifest.toml that
// runs has its handler here, keyed by its path.
var Extension = extension.Extension{
	Manifest: manifest,
	Locales:  locales,
	Skill:    skill,
	Handlers: map[string]extension.Handler{
		"image shrink": shrink,
	},
}
