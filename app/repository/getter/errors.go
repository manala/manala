package getter

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/manala/manala/internal/serrors"
)

var detectErrorRegex = regexp.MustCompile(`(?s)^error downloading '(?P<error>.*)'$`)

func IsNotDetected(err error) bool {
	// Go-getter don't provide a convenient way to catch a getter non detection...
	return detectErrorRegex.MatchString(err.Error())
}

var (
	commandErrorCodeRegex = regexp.MustCompile(`(?s)(?P<command>.+) exited with (?P<code>\d+): (?P<details>.*)$`)
	commandErrorRegex     = regexp.MustCompile(`(?s)error running (?P<command>[^(: )]+): (?P<details>.*)$`)
	multiErrorRegex       = regexp.MustCompile(`(?s)error downloading '.*': \d+ errors occurred:\n(?P<details>.*)\n\n$`)
)

// Mimic the aws sdk error interface to avoid direct dependency on it.
type awsError interface {
	error
	Code() string
	Message() string
	OrigErr() error
}

func NewError(err error) serrors.Error {
	message := err.Error()

	if message == "subdirectory component contain path traversal out of the repository" {
		return serrors.New("subdir out of repository")
	} else
	// Aws error
	if err, ok := errors.AsType[awsError](err); ok {
		var arguments []any
		if code := err.Code(); code != "" {
			arguments = append(arguments, "code", code)
		}

		if message := err.Message(); message != "" {
			arguments = append(arguments, "message", message)
		}

		return serrors.New("aws error").
			WithArguments(arguments...).
			WithErrors(err.OrigErr()).
			WithDetails(err.Error())
	} else
	// Command error code
	if matches := commandErrorCodeRegex.FindStringSubmatch(message); matches != nil {
		code, _ := strconv.Atoi(matches[commandErrorCodeRegex.SubexpIndex("code")])
		return serrors.New("command error").
			WithArguments(
				"command", matches[commandErrorCodeRegex.SubexpIndex("command")],
				"code", code,
			).
			WithDetails(matches[commandErrorCodeRegex.SubexpIndex("details")])
	} else
	// Command error
	if matches := commandErrorRegex.FindStringSubmatch(message); matches != nil {
		return serrors.New("command error").
			WithArguments("command", matches[commandErrorRegex.SubexpIndex("command")]).
			WithDetails(matches[commandErrorRegex.SubexpIndex("details")])
	} else
	// Multi error
	if matches := multiErrorRegex.FindStringSubmatch(message); matches != nil {
		return serrors.New("unable to handle repository").
			WithDetails(matches[multiErrorRegex.SubexpIndex("details")])
	}

	return serrors.New("unable to handle repository").
		WithArguments("error", err.Error())
}
