package serrors

type Dumper interface {
	Dump(ansi bool) string
}

type StringDumper string

func (dumper StringDumper) Dump(_ bool) string {
	return string(dumper)
}

type ErrorDumper interface {
	ErrorDump(ansi bool) string
}
