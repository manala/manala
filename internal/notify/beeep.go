package notify

import "github.com/gen2brain/beeep"

type BeeepHandler struct {
	title string
}

func NewBeeepHandler(title string) *BeeepHandler {
	return &BeeepHandler{
		title: title,
	}
}

func (h *BeeepHandler) Message(message string) {
	_ = beeep.Notify(h.title, message, "")
}

func (h *BeeepHandler) Error(err error) {
	_ = beeep.Alert(h.title, err.Error(), "")
}
