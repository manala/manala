package notifier

import "github.com/gen2brain/beeep"

type Beeep struct {
	title string
}

func NewBeeep(title string) *Beeep {
	return &Beeep{
		title: title,
	}
}

func (notifier *Beeep) Message(message string) {
	_ = beeep.Notify(notifier.title, message, "")
}

func (notifier *Beeep) Error(err error) {
	_ = beeep.Alert(notifier.title, err.Error(), "")
}
