package ui

import (
	"github.com/stretchr/testify/mock"
	"manala/internal/ui/components"
)

type Output interface {
	Message(message *components.Message)
	Error(err error)
	List(header string, list []components.ListItem) error
}

type OutputMock struct {
	mock.Mock
}

func (mock *OutputMock) Message(message *components.Message) {
	mock.Called(message)
}

func (mock *OutputMock) Error(err error) {
	mock.Called(err)
}

func (mock *OutputMock) List(header string, list []components.ListItem) error {
	args := mock.Called(header, list)
	return args.Error(0)
}
