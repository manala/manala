package notifier

type Notifier interface {
	Message(message string)
	Error(err error)
}
