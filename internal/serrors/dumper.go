package serrors

import "io"

type StringDumper string

func (s StringDumper) Dump(w io.Writer) {
	_, _ = io.WriteString(w, string(s))
}
