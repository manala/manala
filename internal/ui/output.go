package ui

import (
	"github.com/stretchr/testify/mock"
	"manala/internal/ui/components"
)

type Output interface {
	Message(message *components.Message)
	Error(err error)
	List(header string, list []components.ListItem) error
	Animate(animation components.Animation, repeat int) error
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

func (mock *OutputMock) Animate(animation components.Animation, repeat int) error {
	args := mock.Called(animation, repeat)
	return args.Error(0)
}
