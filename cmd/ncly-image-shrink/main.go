// Command ncly-image-shrink runs the image-shrink extension of Nicely. Its main
// starts nothing before extension.Main, so --help and --version stay safe.
// The import name ext never collides with the domain, such as fmt.
package main

import (
	"github.com/pascalandy/nicely/sdk/extension"

	ext "github.com/pascalandy/ncly-image-shrink"
)

func main() { extension.Main(ext.Extension) }
