package validation

import (
	"github.com/xeipuuv/gojsonschema"
	internalReport "manala/internal/report"
)

type Reporter interface {
	Report(result gojsonschema.ResultError, report *internalReport.Report)
}
