package components

import "manala/internal/accessor"

func NewListForm(list []ListItem, accessor accessor.IndexAccessor) (*ListForm, error) {
	// Initial index
	index, err := accessor.GetIndex()
	if err != nil {
		return nil, err
	}

	return &ListForm{
		List:     list,
		index:    index,
		accessor: accessor,
	}, nil
}

type ListForm struct {
	List     []ListItem
	index    int
	accessor accessor.IndexAccessor
}

func (form *ListForm) GetIndex() int {
	return form.index
}

func (form *ListForm) SetIndex(index int) {
	form.index = index
}

func (form *ListForm) Submit() error {
	return form.accessor.SetIndex(form.index)
}
