package notifier

func NewNil() *Nil {
	return &Nil{}
}

type Nil struct{}

func (notifier *Nil) Message(_ string) {}

func (notifier *Nil) Error(_ error) {}
