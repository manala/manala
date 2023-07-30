package serrors

type ErrorArguments interface {
	ErrorArguments() []any
}

func NewArguments() *Arguments {
	return &Arguments{}
}

type Arguments struct {
	arguments []any
}

func (arguments *Arguments) ErrorArguments() []any {
	return arguments.arguments
}

func (arguments *Arguments) AppendArguments(args ...any) {
	arguments.arguments = append(arguments.arguments, args...)
}

func (arguments *Arguments) PrependArguments(args ...any) {
	arguments.arguments = append(args, arguments.arguments...)
}
