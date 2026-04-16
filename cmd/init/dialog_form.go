package init

import (
	"fmt"
	"slices"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe/option"
	"github.com/manala/manala/internal/accessor"
	"github.com/manala/manala/internal/output"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"

	"codeberg.org/tslocum/cview"
	"github.com/gdamore/tcell/v3/color"
)

type DialogForm struct {
	*cview.Form

	profile output.Profile
	errored func(error)
	applied func()
}

func NewDialogForm(title string, profile output.Profile) (*DialogForm, *cview.Flex) {
	// Form
	form := &DialogForm{
		Form:    cview.NewForm(),
		profile: profile,
	}
	form.SetWrapAround(true)
	form.SetPadding(1, 1, 2, 1)
	form.SetItemPadding(0)
	form.SetButtonsAlign(cview.AlignLeft)
	form.SetBackgroundColor(color.Default)
	form.SetLabelColor(profile.MutedColor())
	form.SetLabelColorFocused(profile.Color())
	form.SetFieldBackgroundColor(profile.Color())
	form.SetFieldBackgroundColorFocused(profile.MutedColor())
	form.SetFieldTextColor(profile.Color())
	form.SetFieldTextColorFocused(profile.ReverseColor())
	form.SetButtonTextColor(profile.ReverseColor())
	form.SetButtonTextColorFocused(profile.Color())
	form.SetButtonBackgroundColor(profile.MutedColor())
	form.SetButtonBackgroundColorFocused(profile.MutedColor())

	return form, NewDialogPanel(title, form, profile)
}

func (form *DialogForm) SetErroredFunc(handler func(error)) {
	form.errored = handler
}

func (form *DialogForm) SetAppliedFunc(handler func()) {
	form.applied = handler
}

func (form *DialogForm) Build(options []app.RecipeOption, vars *map[string]any) error {
	form.Clear(true)

	// Items
	var items []DialogFormItem
	for _, opt := range options {
		switch opt := opt.(type) {
		case *option.Enum:
			item, err := NewSelectFormItem(opt, vars, form.errored, form.profile)
			if err != nil {
				return serrors.New("invalid recipe option").
					With("label", opt.Label()).
					WithErrors(err)
			}
			items = append(items, item)
			form.AddFormItem(item)
		case *option.String:
			item, err := NewDialogTextFormItem(opt, vars, form.errored, form.profile)
			if err != nil {
				return serrors.New("invalid recipe option").
					With("label", opt.Label()).
					WithErrors(err)
			}
			items = append(items, item)
			form.AddFormItem(item)
		default:
			return serrors.New("unknown recipe option").
				With("label", opt.Label())
		}
	}

	// Apply
	form.AddButton("Apply", func() {
		applied := true
		for _, item := range items {
			applied = item.Apply() && applied
		}
		if applied && form.applied != nil {
			form.applied()
		}
	})

	return nil
}

/*********/
/* Items */
/*********/

type DialogFormItem interface {
	Apply() bool
}

type DialogTextFormItem struct {
	*cview.InputField

	option    *option.String
	accessor  accessor.Accessor
	validator *schema.Validator
	profile   output.Profile
	errored   func(error)
}

func NewDialogTextFormItem(
	option *option.String,
	vars *map[string]any,
	errored func(error),
	profile output.Profile,
) (*DialogTextFormItem, error) {
	// Accessor
	itemAccessor := option.Accessor(vars)

	// Item
	item := &DialogTextFormItem{
		InputField: cview.NewInputField(),
		option:     option,
		accessor:   itemAccessor,
		validator:  option.Validator(),
		errored:    errored,
		profile:    profile,
	}

	// Input field
	item.SetFieldNoteTextColor(profile.MutedColor())
	item.SetLabel(option.Label())
	item.SetChangedFunc(func(_ string) {
		item.Apply()
	})

	// Initial value
	value, err := itemAccessor.Get()
	if err != nil {
		return nil, err
	}
	if value, ok := value.(string); ok {
		item.SetText(value)
	}

	return item, nil
}

func (item *DialogTextFormItem) Apply() bool {
	value := item.GetText()

	// Validation
	violations, err := item.validator.Validate(value)
	if err != nil {
		item.errored(serrors.New("validation error").
			With("label", item.option.Label()).
			WithErrors(err),
		)

		return false
	}
	if violations != nil {
		if errs := violations.Errors(); len(errs) > 0 {
			item.SetFieldNoteTextColor(item.profile.ErrorColor())
			item.SetFieldNote(errs[0].Error())
		}

		return false
	}

	item.SetFieldNote(item.option.Help())
	item.SetFieldNoteTextColor(item.profile.MutedColor())

	// Accession
	if err := item.accessor.Set(value); err != nil {
		item.errored(serrors.New("accession error").
			With("label", item.option.Label()).
			WithErrors(err),
		)

		return false
	}

	return true
}

type DialogSelectFormItem struct {
	*cview.DropDown

	option   *option.Enum
	values   []any
	accessor accessor.Accessor
	errored  func(error)
}

func NewSelectFormItem(
	option *option.Enum,
	vars *map[string]any,
	errored func(error),
	profile output.Profile,
) (*DialogSelectFormItem, error) {
	// Accessor
	itemAccessor := option.Accessor(vars)

	// Item
	item := &DialogSelectFormItem{
		DropDown: cview.NewDropDown(),
		option:   option,
		accessor: itemAccessor,
		values:   option.Values(),
		errored:  errored,
	}

	// Dropdown
	item.SetDropDownTextColor(profile.ReverseColor())
	item.SetDropDownBackgroundColor(profile.MutedColor())
	item.SetDropDownSelectedTextColor(profile.EmphasisColor())
	item.SetDropDownSelectedBackgroundColor(profile.MutedColor())
	item.SetLabel(option.Label())
	item.SetSelectedFunc(func(_ int, _ *cview.DropDownOption) {
		item.Apply()
	})

	for _, value := range item.values {
		var text string
		switch value {
		case nil:
			text = "<None>"
		case true:
			text = "<True>"
		case false:
			text = "<False>"
		default:
			text = fmt.Sprintf("%v", value)
		}
		item.AddOptions(cview.NewDropDownOption(text))
	}

	// Initial value
	value, err := itemAccessor.Get()
	if err != nil {
		return nil, err
	}

	i := slices.Index(item.values, value)
	if i != -1 {
		item.SetCurrentOption(i)
	} else {
		item.SetCurrentOption(0)
	}

	return item, nil
}

func (item *DialogSelectFormItem) Apply() bool {
	i, _ := item.GetCurrentOption()
	value := item.values[i]

	// Accession
	if err := item.accessor.Set(value); err != nil {
		item.errored(serrors.New("accession error").
			With("label", item.option.Label()).
			WithErrors(err),
		)

		return false
	}

	return true
}
