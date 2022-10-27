package repository

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"manala/core"
	internalReport "manala/internal/report"
	"regexp"
	"strconv"
)

func NewGetterResult() *GetterResult {
	return &GetterResult{}
}

type GetterResult struct {
	detectError *GetterError
	getErrors   []*GetterError
}

func (result *GetterResult) SetDetectError(err error, protocol string) {
	result.detectError = NewGetterError(err, protocol)
}

func (result *GetterResult) AddGetError(err error, protocol string) {
	result.getErrors = append(result.getErrors, NewGetterError(err, protocol))
}

func (result *GetterResult) HandleError(err error) *internalReport.Error {
	// No error :)
	if err == nil {
		return nil
	}

	// Only one detection error can occur, because getter client breaks the getter detection loop
	// at the first encountered one.
	if result.detectError != nil {
		return internalReport.NewError(result.detectError).
			WithMessage("unable to detect repository url protocol")
	}

	// Multiple errors can occur during the getters "get" phase
	if len(result.getErrors) == 1 {
		return internalReport.NewError(result.getErrors[0]).
			WithMessage("unable to load repository url")
	}
	if len(result.getErrors) > 1 {
		_err := internalReport.NewError(fmt.Errorf("unable to load repository url"))

		for _, __err := range result.getErrors {
			_err = _err.Add(__err)
		}

		return _err
	}

	// An error matching "downloadingErrorRegex" means no getters has detected url,
	// and no errors were returned during these detections.
	if downloadingErrorRegex.MatchString(err.Error()) {
		return internalReport.NewError(
			core.NewUnsupportedRepositoryError("unsupported repository url"),
		)
	}

	// Unknown error
	return internalReport.NewError(err)
}

func NewGetterError(err error, protocol string) *GetterError {
	newError := &GetterError{
		error:    err,
		protocol: protocol,
		fields:   map[string]interface{}{},
	}

	message := err.Error()

	// Aws error
	if awsErr, ok := err.(awserr.Error); ok {
		newError.message = "aws sdk error"
		newError.trace = fmt.Sprintln(awsErr.Error())
	} else
	// Command error code
	if matches := commandErrorCodeRegex.FindStringSubmatch(message); matches != nil {
		newError.message = "command error"
		newError.fields["command"] = matches[1]
		if code, _err := strconv.Atoi(matches[2]); _err == nil {
			newError.fields["code"] = code
		}
		newError.trace = matches[3]
	} else
	// Command error
	if matches := commandErrorRegex.FindStringSubmatch(message); matches != nil {
		newError.message = "command error"
		newError.fields["command"] = matches[1]
		newError.trace = matches[2]
	}

	return newError
}

type GetterError struct {
	error
	protocol string
	message  string
	fields   map[string]interface{}
	trace    string
}

func (err *GetterError) Unwrap() error {
	return err.error
}

func (err *GetterError) Error() string {
	if err.message != "" {
		return err.message
	}

	return err.error.Error()
}

func (err *GetterError) Report(report *internalReport.Report) {
	report.Compose(
		internalReport.WithField("protocol", err.protocol),
	)

	// Fields
	for key, value := range err.fields {
		report.Compose(
			internalReport.WithField(key, value),
		)
	}

	// Trace
	if err.trace != "" {
		report.Compose(
			internalReport.WithTrace(err.trace),
		)
	}
}

var downloadingErrorRegex = regexp.MustCompile(`(?m)^error downloading '.*'$`)

// 1: command
// 2: code
// 3: trace (optional)
var commandErrorCodeRegex = regexp.MustCompile(`(?s)(.+) exited with (\d+): (.*)$`)

// 1: command
// 2: trace (optional)
var commandErrorRegex = regexp.MustCompile(`(?s)error running ([^(: )]+): (.*)$`)
