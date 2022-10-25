package validation

import (
	"github.com/xeipuuv/gojsonschema"
	internalReport "manala/internal/report"
	"regexp"
)

func NewError(message string, result *gojsonschema.Result, options ...ErrorOption) *Error {
	err := &Error{
		message: message,
		result:  result,
	}

	for _, option := range options {
		option(err)
	}

	return err
}

type Error struct {
	message   string
	result    *gojsonschema.Result
	reporters []Reporter
	messages  []ErrorMessage
}

func (err *Error) Error() string {
	return err.message
}

func (err *Error) Report(report *internalReport.Report) {
	// Range over result errors
	for _, result := range err.result.Errors() {
		rep := internalReport.NewReport(result.Description())

		// Custom messages
		custom := false
		for _, customMessage := range err.messages {
			if ok, message := customMessage.Match(result); ok {
				rep.Compose(
					internalReport.WithMessage(message),
				)
				custom = true
			}
		}

		if !custom {
			// Special error type treatments
			switch result.(type) {
			case *gojsonschema.InvalidTypeError:
				rep.Compose(
					internalReport.WithMessage("invalid type"),
					internalReport.WithField("expected", result.Details()["expected"]),
					internalReport.WithField("given", result.Details()["given"]),
				)
			case *gojsonschema.RequiredError:
				rep.Compose(
					internalReport.WithMessage("missing property"),
					internalReport.WithField("property", result.Details()["property"]),
				)
			case *gojsonschema.AdditionalPropertyNotAllowedError:
				rep.Compose(
					internalReport.WithMessage("additional property is not allowed"),
					internalReport.WithField("property", result.Details()[gojsonschema.STRING_PROPERTY]),
				)
			}
		}

		for _, reporter := range err.reporters {
			reporter.Report(result, rep)
		}

		report.Add(rep)
	}
}

type ErrorOption func(err *Error)

func WithMessages(messages []ErrorMessage) ErrorOption {
	return func(err *Error) {
		err.messages = append(err.messages, messages...)
	}
}

type ErrorMessage struct {
	Field      string
	FieldRegex *regexp.Regexp
	Type       string
	Property   string
	Message    string
}

func (message *ErrorMessage) Match(result gojsonschema.ResultError) (bool, string) {
	field := result.Field()

	// Try to match on field
	if message.Field != "" && message.Field != field {
		return false, ""
	}
	// Try to match on path regex
	if message.FieldRegex != nil && !message.FieldRegex.MatchString(field) {
		return false, ""
	}
	// Try to match on type
	if message.Type != "" && message.Type != result.Type() {
		return false, ""
	}
	// Try to match on property
	if message.Property != "" && message.Property != result.Details()["property"] {
		return false, ""
	}

	return true, message.Message
}
