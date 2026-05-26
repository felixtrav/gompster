// Package samples embeds static sample files served by gompster.
package samples

import _ "embed"

// UTF8Demo is the UTF-8 demo text from Markus Kuhn, University of Cambridge.
// Source: https://www.cl.cam.ac.uk/~mgk25/ucs/examples/UTF-8-demo.txt
// Licensed under CC BY.
//
//go:embed UTF8.txt
var UTF8Demo []byte
