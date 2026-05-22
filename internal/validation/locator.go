package validation

type Locator interface {
	ValueAt(pointer string) (int, int)
	PropertyAt(pointer string) (int, int)
}

type zeroLocator struct{}

func (zeroLocator) ValueAt(location string) (int, int)    { return 0, 0 }
func (zeroLocator) PropertyAt(location string) (int, int) { return 0, 0 }
