package components

func NewForm(fields []FormField) *Form {
	return &Form{
		Fields: fields,
	}
}

type Form struct {
	Fields []FormField
}

func (form *Form) Submit() (bool, error) {
	ok := true
	for _, field := range form.Fields {
		_ok, err := field.Submit()
		if err != nil {
			return false, err
		}
		ok = _ok && ok
	}
	return ok, nil
}
