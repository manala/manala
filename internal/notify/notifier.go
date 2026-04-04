package notify

type Handler interface {
	Message(message string)
	Error(err error)
}

type Notifier struct {
	handler Handler
}

func New(h Handler) *Notifier {
	return &Notifier{handler: h}
}

func (n *Notifier) Message(message string) {
	n.handler.Message(message)
}

func (n *Notifier) Error(err error) {
	n.handler.Error(err)
}
