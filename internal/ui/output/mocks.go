package output

import (
	"github.com/stretchr/testify/mock"
	"manala/internal/ui/components"
)

type Mock struct {
	mock.Mock
}

func (mock *Mock) Message(message *components.Message) {
	mock.Called(message)
}

func (mock *Mock) Error(err error) {
	mock.Called(err)
}

func (mock *Mock) Table(table *components.Table) {
	mock.Called(table)
}
