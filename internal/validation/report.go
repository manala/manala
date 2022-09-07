package validation

import (
	"github.com/xeipuuv/gojsonschema"
	internalReport "manala/internal/report"
)

type Reporter interface {
	Report(result gojsonschema.ResultError, report *internalReport.Report)
}

func WithReporter(reporter Reporter) ErrorOption {
	return func(err *Error) {
		err.reporters = append(err.reporters, reporter)
	}
}
