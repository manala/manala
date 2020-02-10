package validator

import (
	"fmt"
	"github.com/mingrammer/commonregex"
	"github.com/xeipuuv/gojsonschema"
	"manala/models"
	"regexp"
)

func ValidateProject(prj models.ProjectInterface) error {
	result, err := gojsonschema.Validate(
		gojsonschema.NewGoLoader(prj.Recipe().Schema()),
		gojsonschema.NewGoLoader(prj.Vars()),
	)
	if err != nil {
		return err
	}
	if !result.Valid() {
		err := "project config errors:"
		for _, e := range result.Errors() {
			err += "\n- " + e.String()
		}
		return fmt.Errorf(err)
	}

	return nil
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

var DomainRegex = regexp.MustCompile(`^([a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.)+[a-zA-Z]{2,}$`)

func (f DomainFormatChecker) IsFormat(input interface{}) bool {
	return DomainRegex.MatchString(input.(string))
}

func init() {
	gojsonschema.FormatCheckers.Add("git-repo", GitRepoFormatChecker{})
	gojsonschema.FormatCheckers.Add("file-path", FilePathFormatChecker{})
	gojsonschema.FormatCheckers.Add("domain", DomainFormatChecker{})
}
