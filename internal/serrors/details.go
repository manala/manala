package serrors

type ErrorDetails interface {
	ErrorDetails(ansi bool) string
}
