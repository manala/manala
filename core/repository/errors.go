package repository

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"manala/core"
	"manala/internal/errors/serrors"
	"regexp"
	"strconv"
)

var downloadingErrorRegex = regexp.MustCompile(`(?m)^error downloading '.*'$`)

func NewGetterResult(url string) *GetterResult {
	return &GetterResult{
		url: url,
	}
}

type GetterResult struct {
	url         string
	detectError error
	getErrors   []error
}

func (result *GetterResult) SetDetectError(err error, protocol string) {
	result.detectError = NewGetterError(err, protocol)
}

func (result *GetterResult) AddGetError(err error, protocol string) {
	result.getErrors = append(result.getErrors, NewGetterError(err, protocol))
}

func (result *GetterResult) Error(err error) error {
	// No error :)
	if err == nil {
		return nil
	}

	// Only one detection error can occur, because getter client breaks the getter detection loop
	// at the first encountered one.
	if result.detectError != nil {
		return serrors.Wrap("unable to detect repository url protocol", result.detectError).
			WithArguments("url", result.url)
	}

	// Multiple errors can occur during the getters "get" phase
	if len(result.getErrors) > 0 {
		return serrors.Wraps("unable to load repository url", result.getErrors...).
			WithArguments("url", result.url)
	}

	// An error matching "downloadingErrorRegex" means no getters has detected url,
	// and no errors were returned during these detections.
	if downloadingErrorRegex.MatchString(err.Error()) {
		return &core.UnsupportedRepositoryError{Url: result.url}
	}

	// Unknown error
	return err
}

// 1: command
// 2: code
// 3: details (optional)
var commandErrorCodeRegex = regexp.MustCompile(`(?s)(.+) exited with (\d+): (.*)$`)

// 1: command
// 2: details (optional)
var commandErrorRegex = regexp.MustCompile(`(?s)error running ([^(: )]+): (.*)$`)

func NewGetterError(err error, protocol string) *GetterError {
	_err := &GetterError{
		message:   err.Error(),
		Arguments: serrors.NewArguments(),
		Details:   serrors.NewDetails(),
	}

	// Protocol
	_err.AppendArguments("protocol", protocol)

	// Aws error
	if awsErr, ok := err.(awserr.Error); ok {
		_err.message = "aws sdk error"
		_err.SetDetails(awsErr.Error())
	} else
	// Command error code
	if matches := commandErrorCodeRegex.FindStringSubmatch(_err.message); matches != nil {
		_err.message = "command error"
		_err.AppendArguments("command", matches[1])
		if code, __err := strconv.Atoi(matches[2]); __err == nil {
			_err.AppendArguments("code", code)
		}
		_err.SetDetails(matches[3])
	} else
	// Command error
	if matches := commandErrorRegex.FindStringSubmatch(_err.message); matches != nil {
		_err.message = "command error"
		_err.AppendArguments("command", matches[1])
		_err.SetDetails(matches[2])
	}

	return _err
}

type GetterError struct {
	message string
	*serrors.Arguments
	*serrors.Details
}

func (err *GetterError) Error() string {
	return err.message
}
