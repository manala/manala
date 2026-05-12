package validation

type Locator interface {
	At(pointer string) (int, int)
}

type zeroLocator struct{}

func (zeroLocator) At(string) (int, int) { return 0, 0 }
