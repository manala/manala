package log

import "github.com/manala/manala/internal/output"

var Discard = &Log{out: output.Discard}
