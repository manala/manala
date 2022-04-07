package yaml

import (
	"github.com/goccy/go-yaml"
	internalErrors "manala/internal/errors"
	"regexp"
	"strconv"
)

var errorRegex = regexp.MustCompile(`\x1b\[91m\[(?P<line>\d+):(?P<column>\d+)] (?P<message>.*)\x1b\[0m\n`)

func Error(file string, err error) *internalErrors.InternalError {
	_err := internalErrors.New("yaml processing error").
		WithField("file", file)

	trace := yaml.FormatError(err, true, true)
	if trace != err.Error() {
		if match := errorRegex.FindStringSubmatch(trace); match != nil {
			if line, err := strconv.Atoi(match[1]); err == nil {
				_ = _err.WithField("line", line)
			}
			if column, err := strconv.Atoi(match[2]); err == nil {
				_ = _err.WithField("column", column)
			}
			_ = _err.
				WithField("message", match[3]).
				WithTrace(errorRegex.ReplaceAllLiteralString(trace, ""))
		} else {
			_ = _err.WithTrace(trace)
		}
	} else {
		_ = _err.WithError(err)
	}

	return _err
}

func CommentTagError(path string, err error) *internalErrors.InternalError {
	return internalErrors.New("yaml comment tag error").
		WithField("path", path).
		WithError(err)
}
