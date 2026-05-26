// Package images embeds sample image files served by the /image endpoints.
// Replace each stub file with a real image of the matching format, then rebuild.
package images

import _ "embed"

//go:embed sample.jpeg
var JPEG []byte

//go:embed sample.png
var PNG []byte

//go:embed sample.svg
var SVG []byte

//go:embed sample.webp
var WebP []byte

//go:embed sample.avif
var AVIF []byte

//go:embed sample.gif
var GIF []byte

//go:embed sample.apng
var APNG []byte
