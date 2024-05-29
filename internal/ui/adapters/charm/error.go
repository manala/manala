package charm

import (
	"manala/internal/ui/components"

	"github.com/muesli/termenv"
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
