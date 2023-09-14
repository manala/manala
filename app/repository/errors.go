package repository

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"manala/app"
	"manala/internal/serrors"
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
		return serrors.New("unable to detect repository url protocol").
			WithArguments("url", result.url).
			WithErrors(result.detectError)
	}

	// Multiple errors can occur during the getters "get" phase
	if len(result.getErrors) > 0 {
		return serrors.New("unable to load repository url").
			WithArguments("url", result.url).
			WithErrors(result.getErrors...)
	}

	// An error matching "downloadingErrorRegex" means no getters has detected url,
	// and no errors were returned during these detections.
	if downloadingErrorRegex.MatchString(err.Error()) {
		return &app.UnsupportedRepositoryError{Url: result.url}
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

func NewGetterError(err error, protocol string) serrors.Error {
	message := err.Error()
	arguments := []any{"protocol", protocol}
	details := ""

	// Aws error
	if awsErr, ok := err.(awserr.Error); ok {
		message = "aws sdk error"
		details = awsErr.Error()
	} else
	// Command error code
	if matches := commandErrorCodeRegex.FindStringSubmatch(message); matches != nil {
		message = "command error"
		arguments = append(arguments, "command", matches[1])
		if code, _err := strconv.Atoi(matches[2]); _err == nil {
			arguments = append(arguments, "code", code)
		}
		details = matches[3]
	} else
	// Command error
	if matches := commandErrorRegex.FindStringSubmatch(message); matches != nil {
		message = "command error"
		arguments = append(arguments, "command", matches[1])
		details = matches[2]
	}

	return serrors.New(message).
		WithArguments(arguments...).
		WithDetails(details)
}
