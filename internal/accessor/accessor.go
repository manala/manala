package accessor

type Accessor interface {
	Get() (any, error)
	Set(value any) error
}

func New(value *any) Accessor {
	return accessor{
		value: value,
	}
}

type accessor struct {
	value *any
}

func (accessor accessor) Get() (any, error) {
	return *accessor.value, nil
}

func (accessor accessor) Set(value any) error {
	*accessor.value = value
	return nil
}
