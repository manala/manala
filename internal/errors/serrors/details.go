package serrors

type ErrorDetails interface {
	ErrorDetails(ansi bool) string
}

func NewDetails() *Details {
	return &Details{}
}

type Details struct {
	details string
}

func (details *Details) SetDetails(str string) {
	details.details = str
}

func (details *Details) ErrorDetails(_ bool) string {
	return details.details
}
