package charm

import (
	"github.com/muesli/termenv"
	"manala/internal/ui/components"
)

func (adapter *Adapter) Error(err error) {
	renderer := adapter.errRenderer

	style := messageStyle.New(renderer)

	_, _ = renderer.Output().WriteString(
		style.Render(
			adapter.message(
				components.MessageFromError(
					err,
					renderer.ColorProfile() != termenv.Ascii,
				),
				renderer,
			),
		) + "\n",
	)
}
