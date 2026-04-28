package serrors

import (
	"github.com/manala/manala/internal/output"
)

type StringDumper string

func (s StringDumper) Dump(_ output.Profile) string {
	return string(s)
}
