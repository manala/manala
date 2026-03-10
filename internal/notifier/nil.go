package notifier

type Nil struct{}

func NewNil() *Nil {
	return &Nil{}
}

func (notifier *Nil) Message(_ string) {}

func (notifier *Nil) Error(_ error) {}
