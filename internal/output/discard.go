package output

import "io"

var Discard = Output{out: io.Discard}
