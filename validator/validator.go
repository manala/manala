package validator

import (
	"github.com/mingrammer/commonregex"
	"github.com/xeipuuv/gojsonschema"
	"manala/models"
	"regexp"
)

func ValidateValue(value interface{}, schema map[string]interface{}) error {
	result, err := gojsonschema.Validate(
		gojsonschema.NewGoLoader(schema),
		gojsonschema.NewGoLoader(value),
	)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return &ValueValidationError{
			Errors: result.Errors(),
		}
	}

	return nil
}

type ValueValidationError struct {
	Errors []gojsonschema.ResultError
}

func (err *ValueValidationError) Error() string {
	str := ""
	for _, e := range err.Errors {
		str += "\n- " + e.Description()
	}
	return str
}

func ValidateProject(prj models.ProjectInterface) error {
	result, err := gojsonschema.Validate(
		gojsonschema.NewGoLoader(prj.Recipe().Schema()),
		gojsonschema.NewGoLoader(prj.Vars()),
	)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return &ProjectValidationError{
			Errors: result.Errors(),
		}
	}

	return nil
}

type ProjectValidationError struct {
	Errors []gojsonschema.ResultError
}

func (err *ProjectValidationError) Error() string {
	str := "project config errors:"
	for _, e := range err.Errors {
		str += "\n- " + e.String()
	}
	return str
}

/**************************/
/* Custom Format Checkers */
/**************************/

type GitRepoFormatChecker struct{}

func (f GitRepoFormatChecker) IsFormat(input interface{}) bool {
	return commonregex.GitRepoRegex.MatchString(input.(string))
}

type FilePathFormatChecker struct{}

var FilePathRegex = regexp.MustCompile(`^/[a-z0-9-_/]*$`)

func (f FilePathFormatChecker) IsFormat(input interface{}) bool {
	return FilePathRegex.MatchString(input.(string))
}

type DomainFormatChecker struct{}

var DomainRegex = regexp.MustCompile(`^([a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.)+[a-zA-Z]{2,}$`)

func (f DomainFormatChecker) IsFormat(input interface{}) bool {
	return DomainRegex.MatchString(input.(string))
}

func init() {
	gojsonschema.FormatCheckers.Add("git-repo", GitRepoFormatChecker{})
	gojsonschema.FormatCheckers.Add("file-path", FilePathFormatChecker{})
	gojsonschema.FormatCheckers.Add("domain", DomainFormatChecker{})
}
