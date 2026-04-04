package notify

var DiscardHandler Handler = discardHandler{}

type discardHandler struct{}

func (discardHandler) Message(string) {}

func (discardHandler) Error(error) {}
