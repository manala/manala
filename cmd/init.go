package cmd

import (
	"code.rocketnine.space/tslocum/cview"
	"fmt"
	"github.com/apex/log"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
	"manala/app"
	"manala/internal/config"
	"manala/loaders"
	"manala/models"
	"manala/validator"
)

type InitCmd struct{}

func (cmd *InitCmd) Command(conf *config.Config, logger *log.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:     "init [dir]",
		Aliases: []string{"in"},
		Short:   "Init project",
		Long: `Init (manala init) will init a project.

Example: manala init -> resulting in a project init in a directory (default to the current directory)`,
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		RunE: func(command *cobra.Command, args []string) error {
			// App
			_ = conf.BindPFlags(command.PersistentFlags())
			manala := app.New(conf, logger)

			// Command
			flags := config.New()
			_ = flags.BindPFlags(command.Flags())
			return manala.Init(
				cmd.runRecipeListApplication,
				cmd.runRecipeOptionsFormApplication,
				append(args, ".")[0],
				flags.GetString("recipe"),
			)
		},
	}

	// Persistent flags
	pFlags := command.PersistentFlags()
	pFlags.StringP("repository", "o", "", "use repository")

	// Flags
	flags := command.Flags()
	flags.StringP("recipe", "i", "", "use recipe")

	return command
}

func (cmd *InitCmd) runRecipeListApplication(recipeLoader loaders.RecipeLoaderInterface, repo models.RepositoryInterface) (models.RecipeInterface, error) {
	// Application
	application := cview.NewApplication()
	application.EnableMouse(true)

	var err error

	// List
	list := cview.NewList()
	list.SetPadding(0, 0, 1, 0)
	list.SetScrollBarVisibility(cview.ScrollBarAlways)
	list.SetDoneFunc(func() {
		err = fmt.Errorf("operation cancelled")
		application.Stop()
	})

	var recipe models.RecipeInterface

	// Walk into recipes
	if err2 := recipeLoader.Walk(repo, func(rec models.RecipeInterface) {
		listItem := cview.NewListItem(" " + rec.Name() + " ")
		listItem.SetSecondaryText("   " + rec.Description())
		listItem.SetSelectedFunc(func() {
			recipe = rec
			application.Stop()
		})
		list.AddItem(listItem)
	}); err2 != nil {
		return nil, err2
	}

	frame := cview.NewFrame(list)
	frame.SetBorders(1, 1, 1, 1, 1, 1)
	frame.AddText("Please, select a recipe...", true, cview.AlignLeft, tcell.ColorAqua)

	application.SetRoot(frame, true)
	application.SetFocus(frame)
	if err2 := application.Run(); err2 != nil {
		return nil, err2
	}

	if err != nil {
		return nil, err
	}

	if recipe == nil {
		return nil, fmt.Errorf("unknown error")
	}

	return recipe, nil
}

func (cmd *InitCmd) runRecipeOptionsFormApplication(rec models.RecipeInterface, vars map[string]interface{}) error {
	// Application
	application := cview.NewApplication()
	application.EnableMouse(true)

	applicationPanels := cview.NewPanels()

	var err error

	// Form panel
	form := cview.NewForm()
	form.SetPadding(0, 0, 1, 0)
	form.SetItemPadding(0)
	form.SetCancelFunc(func() {
		err = fmt.Errorf("operation cancelled")
		application.Stop()
	})

	frame := cview.NewFrame(form)
	frame.SetBorders(1, 1, 1, 1, 1, 1)
	frame.AddText("Please, enter \""+rec.Name()+"\" recipe options...", true, cview.AlignLeft, tcell.ColorAqua)

	applicationPanels.AddPanel("form", frame, true, true)

	// Modal panel
	modal := cview.NewModal()
	modal.SetBorderColor(tcell.ColorRed)
	modal.SetTextColor(tcell.ColorWhite)
	modal.AddButtons([]string{"Ok"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		applicationPanels.HidePanel("modal")
	})

	applicationPanels.AddPanel("modal", modal, false, false)

	// Recipe form binder
	binder, err2 := NewRecipeFormBinder(rec)
	if err2 != nil {
		return err2
	}

	binder.BindForm(form)

	form.AddButton("Apply", func() {
		// Validate
		valid := true
		for _, bnd := range binder.Binds() {
			err2 := validator.ValidateValue(bnd.Value, bnd.Option.Schema)
			if err2 != nil {
				if err3, ok := err2.(*validator.ValueValidationError); ok {
					valid = false
					modal.SetText(bnd.Option.Label + err3.Error())
					applicationPanels.ShowPanel("modal")
					form.SetFocus(bnd.ItemIndex)
				} else {
					err = err2
					application.Stop()
				}
				break
			}
		}
		if valid && (err == nil) {
			// Apply values
			_ = binder.Apply(vars)
			application.Stop()
		}
	})

	application.SetRoot(applicationPanels, true)
	application.SetFocus(applicationPanels)
	if err3 := application.Run(); err3 != nil {
		return err3
	}

	if err != nil {
		return err
	}

	return nil
}

func init() {
	cview.Styles = cview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorBlack,
		ContrastBackgroundColor:     tcell.ColorBlack,
		MoreContrastBackgroundColor: tcell.ColorWhite,
		BorderColor:                 tcell.ColorRed,
		TitleColor:                  tcell.ColorAqua,
		GraphicsColor:               tcell.ColorWhite,
		PrimaryTextColor:            tcell.ColorLime,
		SecondaryTextColor:          tcell.ColorWhite,
		TertiaryTextColor:           tcell.ColorWhite,
		InverseTextColor:            tcell.ColorYellow,
		ContrastSecondaryTextColor:  tcell.ColorDarkCyan,
		ScrollBarColor:              tcell.ColorWhite,
	}
}
