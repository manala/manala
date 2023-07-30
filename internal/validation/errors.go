package validation

import (
	"github.com/xeipuuv/gojsonschema"
	"manala/internal/errors/serrors"
	"regexp"
)

func NewError(msg string, result *gojsonschema.Result) *Error {
	return &Error{
		message:   msg,
		result:    result,
		Arguments: serrors.NewArguments(),
	}
}

type Error struct {
	message string
	result  *gojsonschema.Result
	*serrors.Arguments
	messages              []ErrorMessage
	resultErrorDecorators []ResultErrorDecorator
}

func (err *Error) Error() string {
	return err.message
}

func (err *Error) Unwrap() []error {
	var errs []error

	for _, _err := range err.result.Errors() {
		__err := NewResultError(_err, err.messages)
		for _, decorator := range err.resultErrorDecorators {
			__err = decorator.Decorate(__err)

		}
		errs = append(errs, __err)
	}

	return errs
}

func (err *Error) WithMessages(messages []ErrorMessage) *Error {
	err.messages = append(err.messages, messages...)
	return err
}

func (err *Error) WithArguments(arguments ...any) *Error {
	err.AppendArguments(arguments...)
	return err
}

func (err *Error) WithResultErrorDecorator(decorator ResultErrorDecorator) *Error {
	err.resultErrorDecorators = append(err.resultErrorDecorators, decorator)
	return err
}

/**********/
/* Result */
/**********/

type ResultErrorDecorator interface {
	Decorate(err ResultErrorInterface) ResultErrorInterface
}

type ResultErrorInterface interface {
	Path() string
	error
	serrors.ErrorArguments
	serrors.ErrorDetails
}

func NewResultError(err gojsonschema.ResultError, messages []ErrorMessage) ResultErrorInterface {
	_err := &ResultError{
		message:   err.String(),
		err:       err,
		Arguments: serrors.NewArguments(),
		Details:   serrors.NewDetails(),
	}

	// Special types
	switch err.(type) {
	case *gojsonschema.InvalidTypeError:
		_err.message = "invalid type"
		_err.AppendArguments(
			"expected", err.Details()["expected"],
			"given", err.Details()["given"],
		)
	case *gojsonschema.RequiredError:
		_err.message = "missing property"
		_err.AppendArguments(
			"property", err.Details()["property"],
		)
	case *gojsonschema.AdditionalPropertyNotAllowedError:
		_err.message = "additional property is not allowed"
		_err.AppendArguments(
			"property", err.Details()["property"],
		)
	}

	// Custom messages
	for _, message := range messages {
		if _message, ok := message.Match(err); ok {
			_err.message = _message
		}
	}

	return _err
}

type ResultError struct {
	message string
	err     gojsonschema.ResultError
	*serrors.Arguments
	*serrors.Details
}

func (err *ResultError) Path() string {
	return err.err.Field()
}

func (err *ResultError) Error() string {
	return err.message
}

/***********/
/* Message */
/***********/

type ErrorMessage struct {
	Field      string
	FieldRegex *regexp.Regexp
	Type       string
	Property   string
	Message    string
}

func (message *ErrorMessage) Match(result gojsonschema.ResultError) (string, bool) {
	field := result.Field()

	// Try to match on field
	if message.Field != "" && message.Field != field {
		return "", false
	}
	// Try to match on path regex
	if message.FieldRegex != nil && !message.FieldRegex.MatchString(field) {
		return "", false
	}
	// Try to match on type
	if message.Type != "" && message.Type != result.Type() {
		return "", false
	}
	// Try to match on property
	if message.Property != "" && message.Property != result.Details()["property"] {
		return "", false
	}

	return message.Message, true
}
