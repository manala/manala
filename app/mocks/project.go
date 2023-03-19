package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"manala/app/interfaces"
	internalReport "manala/internal/report"
	internalTemplate "manala/internal/template"
)

func MockProject() *ProjectMock {
	return &ProjectMock{}
}

type ProjectMock struct {
	mock.Mock
}

func (rec *ProjectMock) Dir() string {
	args := rec.Called()
	return args.String(0)
}

func (rec *ProjectMock) Recipe() interfaces.Recipe {
	args := rec.Called()
	return args.Get(0).(interfaces.Recipe)
}

func (rec *ProjectMock) Vars() map[string]interface{} {
	args := rec.Called()
	return args.Get(0).(map[string]interface{})
}

func (rec *ProjectMock) Template() *internalTemplate.Template {
	args := rec.Called()
	return args.Get(0).(*internalTemplate.Template)
}

/************/
/* Manifest */
/************/

func MockProjectManifest() *ProjectManifestMock {
	return &ProjectManifestMock{}
}

type ProjectManifestMock struct {
	mock.Mock
}

func (man *ProjectManifestMock) Recipe() string {
	args := man.Called()
	return args.String(0)
}

func (man *ProjectManifestMock) Repository() string {
	args := man.Called()
	return args.String(0)
}

func (man *ProjectManifestMock) Vars() map[string]interface{} {
	args := man.Called()
	return args.Get(0).(map[string]interface{})
}

func (man *ProjectManifestMock) ReadFrom(reader io.Reader) error {
	args := man.Called(reader)
	return args.Error(0)
}

func (man *ProjectManifestMock) Report(result gojsonschema.ResultError, report *internalReport.Report) {
	man.Called(result, report)
}
