package ui

import (
	"github.com/stretchr/testify/mock"
	"manala/internal/ui/components"
)

type Input interface {
	ListForm(header string, form *components.ListForm) error
	Form(header string, form *components.Form) error
}

type InputMock struct {
	mock.Mock
}

func (mock *InputMock) ListForm(header string, form *components.ListForm) error {
	args := mock.Called(header, form)
	return args.Error(0)
}

func (mock *InputMock) Form(header string, form *components.Form) error {
	args := mock.Called(header, form)
	return args.Error(0)
}
