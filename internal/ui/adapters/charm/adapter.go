package charm

import (
	"io"
)

func New(err io.Writer) *Adapter {
	return &Adapter{
		err: err,
	}
}

type Adapter struct {
	err io.Writer
}
