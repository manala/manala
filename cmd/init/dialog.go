package init

import (
	"github.com/manala/manala/app"
	"github.com/manala/manala/cmd"
	"github.com/manala/manala/internal/output"

	"codeberg.org/tslocum/cview"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type (
	DialogVariant       any
	DialogSingleVariant struct {
		Recipe app.Recipe
	}
	DialogMultiVariant struct {
		Recipes []app.Recipe
	}
)

type DialogOutcome struct {
	Recipe app.Recipe
	Vars   map[string]any
}

func RunDialog(title string, variant DialogVariant, profile output.Profile) (*DialogOutcome, error) {
	var outcome *DialogOutcome

	// Dialog
	dialog, panels := NewDialog(title, profile)
	defer dialog.HandlePanic()

	// Form
	form, formPanel := NewDialogForm("Configure recipe", profile)
	form.SetErroredFunc(func(err error) { dialog.Fatal(err) })
	form.SetAppliedFunc(func() { dialog.Stop() })

	switch variant := variant.(type) {
	case DialogSingleVariant:
		outcome = &DialogOutcome{variant.Recipe, variant.Recipe.Vars()}

		options := outcome.Recipe.Options()
		if len(options) == 0 {
			return outcome, nil
		}

		// Form
		form.SetCancelFunc(func() { dialog.Cancel() })
		if err := form.Build(options, &outcome.Vars); err != nil {
			return nil, err
		}

		panels.AddPanel("form", formPanel, true, true)
	case DialogMultiVariant:
		// List
		list, listPanel := NewDialogList("Select a recipe", profile)
		list.SetDoneFunc(func() { dialog.Cancel() })
		list.SetSelectedFunc(func(recipe app.Recipe) {
			outcome = &DialogOutcome{recipe, recipe.Vars()}

			options := outcome.Recipe.Options()
			if len(options) == 0 {
				dialog.Stop()

				return
			}

			if err := form.Build(options, &outcome.Vars); err != nil {
				dialog.Fatal(err)

				return
			}

			panels.SetCurrentPanel("form")
		})

		list.Build(variant.Recipes)

		// Form
		form.SetCancelFunc(func() { panels.SetCurrentPanel("list") })

		panels.AddPanel("list", listPanel, true, true)
		panels.AddPanel("form", formPanel, true, false)
	}

	if err := dialog.Run(); err != nil {
		return nil, err
	}

	return outcome, nil
}

type Dialog struct {
	*cview.Application

	title   string
	profile output.Profile
	panels  *cview.Panels
	err     error
}

func NewDialog(title string, profile output.Profile) (*Dialog, *cview.Panels) {
	application := &Dialog{
		Application: cview.NewApplication(),
		profile:     profile,
	}
	application.EnableMouse(true)
	application.EnableBracketedPaste(true)
	application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			application.Cancel()
			return nil
		}

		return event
	})

	// Title
	application.title = title

	// Panels
	application.panels = cview.NewPanels()

	return application, application.panels
}

func (dialog *Dialog) Run() error {
	// Ensure we're running in a terminal before cview tries to use it
	if tty, err := tcell.NewDevTty(); err != nil {
		return &cmd.TerminalNotFoundError{}
	} else {
		_ = tty.Close()
	}

	if err := dialog.Init(); err != nil {
		return err
	}

	// Screen
	screen := dialog.GetScreen()
	screen.SetTitle(dialog.title)
	screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock, dialog.profile.EmphasisColor())

	// Panels
	dialog.SetRoot(dialog.panels, true)

	if err := dialog.Application.Run(); err != nil {
		return err
	}

	return dialog.err
}

func (dialog *Dialog) Cancel() {
	dialog.err = &cmd.CancelError{}
	dialog.Stop()
}

func (dialog *Dialog) Fatal(err error) {
	dialog.err = err
	dialog.Stop()
}

func NewDialogPanel(title string, item cview.Primitive, profile output.Profile) *cview.Flex {
	panel := cview.NewFlex()
	panel.SetDirection(cview.FlexRow)

	// Header
	header := cview.NewTextView()
	header.SetBackgroundColor(color.Default)
	header.SetText(title)
	header.SetTextColor(profile.Color())
	header.SetPadding(0, 0, 2, 0)
	panel.AddItem(header, 1, 0, false)

	// Separator
	panel.AddItem(NewDialogPanelSeparator(profile), 1, 0, false)

	// Item
	panel.AddItem(item, 0, 1, true)

	return panel
}

type DialogPanelSeparator struct {
	*cview.Box

	profile output.Profile
}

func NewDialogPanelSeparator(profile output.Profile) *DialogPanelSeparator {
	return &DialogPanelSeparator{
		Box:     cview.NewBox(),
		profile: profile,
	}
}

func (separator *DialogPanelSeparator) Draw(screen tcell.Screen) {
	x, y, w, _ := separator.GetInnerRect()
	for i := range w {
		screen.SetContent(
			x+i,
			y,
			'─',
			nil,
			tcell.StyleDefault.Foreground(separator.profile.MutedColor()),
		)
	}
}
