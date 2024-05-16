package getter

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"manala/internal/serrors"
	"regexp"
	"strconv"
)

// 1: error
var detectErrorRegex = regexp.MustCompile(`(?s)^error downloading '(.*)'$`)

func IsNotDetected(err error) bool {
	// Go-getter don't provide a convenient way to catch a getter non detection...
	return detectErrorRegex.MatchString(err.Error())
}

// 1: command
// 2: code
// 3: details (optional)
var commandErrorCodeRegex = regexp.MustCompile(`(?s)(.+) exited with (\d+): (.*)$`)

// 1: command
// 2: details (optional)
var commandErrorRegex = regexp.MustCompile(`(?s)error running ([^(: )]+): (.*)$`)

// 1: details
var multiErrorRegex = regexp.MustCompile(`(?s)error downloading '.*': \d+ errors occurred:\n(.*)\n\n$`)

func NewError(err error) serrors.Error {
	message := err.Error()

	if message == "subdirectory component contain path traversal out of the repository" {
		return serrors.New("subdir out of repository")
	} else
	// Aws error
	if awsErr, ok := err.(awserr.Error); ok {
		return serrors.New("aws sdk error").
			WithDetails(awsErr.Error())
	} else
	// Command error code
	if matches := commandErrorCodeRegex.FindStringSubmatch(message); matches != nil {
		err := serrors.New("command error").
			WithArguments("command", matches[1]).
			WithDetails(matches[3])
		if code, _err := strconv.Atoi(matches[2]); _err == nil {
			err = err.WithArguments("code", code)
		}
		return err
	} else
	// Command error
	if matches := commandErrorRegex.FindStringSubmatch(message); matches != nil {
		return serrors.New("command error").
			WithArguments("command", matches[1]).
			WithDetails(matches[2])
	} else
	// Multi error
	if matches := multiErrorRegex.FindStringSubmatch(message); matches != nil {
		return serrors.New("unable to handle repository").
			WithDetails(matches[1])
	}

	return serrors.New("unable to handle repository").
		WithArguments("error", err.Error())
}
