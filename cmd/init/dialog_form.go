package init

import (
	"fmt"
	"slices"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe/option"
	"github.com/manala/manala/internal/accessor"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"

	"codeberg.org/tslocum/cview"
	"github.com/gdamore/tcell/v3/color"
)

type DialogForm struct {
	*cview.Form

	errored func(error)
	applied func()
}

func NewDialogForm(title string) (*DialogForm, *cview.Flex) {
	// Form
	form := &DialogForm{Form: cview.NewForm()}
	form.SetWrapAround(true)
	form.SetPadding(1, 1, 2, 1)
	form.SetItemPadding(0)
	form.SetButtonsAlign(cview.AlignLeft)
	form.SetBackgroundColor(color.Default)
	form.SetLabelColor(DialogStyles.SecondaryColor)
	form.SetLabelColorFocused(DialogStyles.PrimaryColor)
	form.SetFieldBackgroundColor(DialogStyles.PrimaryColor)
	form.SetFieldBackgroundColorFocused(DialogStyles.SecondaryColor)
	form.SetFieldTextColor(DialogStyles.PrimaryColor)
	form.SetFieldTextColorFocused(DialogStyles.TertiaryColor)
	form.SetButtonTextColor(DialogStyles.TertiaryColor)
	form.SetButtonTextColorFocused(DialogStyles.PrimaryColor)
	form.SetButtonBackgroundColor(DialogStyles.SecondaryColor)
	form.SetButtonBackgroundColorFocused(DialogStyles.SecondaryColor)

	return form, NewDialogPanel(title, form)
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
		case *option.Select:
			item, err := NewSelectFormItem(opt, vars, form.errored)
			if err != nil {
				return serrors.New("invalid recipe option").
					WithArguments("label", opt.Label()).
					WithErrors(err)
			}
			items = append(items, item)
			form.AddFormItem(item)
		case *option.Text:
			item, err := NewDialogTextFormItem(opt, vars, form.errored)
			if err != nil {
				return serrors.New("invalid recipe option").
					WithArguments("label", opt.Label()).
					WithErrors(err)
			}
			items = append(items, item)
			form.AddFormItem(item)
		default:
			return serrors.New("unknown recipe option").
				WithArguments("label", opt.Label())
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

	option    *option.Text
	accessor  accessor.Accessor
	validator *schema.Validator
	errored   func(error)
}

func NewDialogTextFormItem(
	option *option.Text,
	vars *map[string]any,
	errored func(error),
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
	}

	// Input field
	item.SetFieldNoteTextColor(DialogStyles.SecondaryColor)
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
			WithArguments("label", item.option.Label()).
			WithErrors(err),
		)

		return false
	}
	if violations != nil {
		if errs := violations.Errors(); len(errs) > 0 {
			item.SetFieldNoteTextColor(DialogStyles.AlertColor)
			item.SetFieldNote(errs[0].Error())
		}

		return false
	}

	item.SetFieldNote(item.option.Help())
	item.SetFieldNoteTextColor(DialogStyles.SecondaryColor)

	// Accession
	if err := item.accessor.Set(value); err != nil {
		item.errored(serrors.New("accession error").
			WithArguments("label", item.option.Label()).
			WithErrors(err),
		)

		return false
	}

	return true
}

type DialogSelectFormItem struct {
	*cview.DropDown

	option   *option.Select
	values   []any
	accessor accessor.Accessor
	errored  func(error)
}

func NewSelectFormItem(
	option *option.Select,
	vars *map[string]any,
	errored func(error),
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
	item.SetDropDownTextColor(DialogStyles.TertiaryColor)
	item.SetDropDownBackgroundColor(DialogStyles.SecondaryColor)
	item.SetDropDownSelectedTextColor(DialogStyles.QuaternaryColor)
	item.SetDropDownSelectedBackgroundColor(DialogStyles.SecondaryColor)
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
			WithArguments("label", item.option.Label()).
			WithErrors(err),
		)

		return false
	}

	return true
}
