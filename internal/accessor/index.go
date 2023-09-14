package accessor

import (
	"manala/internal/serrors"
	"slices"
)

type IndexAccessor interface {
	GetIndex() (int, error)
	SetIndex(index int) error
}

func NewIndex[I comparable](item *I, items []I) IndexAccessor {
	return &indexAccessor[I]{
		item:  item,
		items: items,
	}
}

type indexAccessor[I comparable] struct {
	item  *I
	items []I
}

func (accessor *indexAccessor[I]) GetIndex() (int, error) {
	if *accessor.item == *new(I) {
		return -1, nil
	}
	index := slices.Index(accessor.items, *accessor.item)
	if index == -1 {
		return 0, serrors.New("invalid item")
	}
	return index, nil
}

func (accessor *indexAccessor[I]) SetIndex(index int) error {
	if index == -1 {
		accessor.item = nil
		return nil
	}
	if index < 0 || index >= len(accessor.items) {
		return serrors.New("invalid item index")
	}
	*accessor.item = accessor.items[index]
	return nil
}
