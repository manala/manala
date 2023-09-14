package schema

import (
	"github.com/xeipuuv/gojsonschema"
	"regexp"
)

func init() {
	gojsonschema.FormatCheckers.
		Add("git-repo", GitRepoFormatChecker{}).
		Add("file-path", FilePathFormatChecker{}).
		Add("domain", DomainFormatChecker{})
}

/************/
/* Git Repo */
/************/

type GitRepoFormatChecker struct{}

func (checker GitRepoFormatChecker) IsFormat(input any) bool {
	return gitRepoRegex.MatchString(input.(string))
}

var gitRepoRegex = regexp.MustCompile(`((git|ssh|http(s)?)|(git@[\w.]+))(:(//)?)([\w.@:/\-~]+)(\.git)(/)?`)

/*************/
/* File Path */
/*************/

type FilePathFormatChecker struct{}

func (checker FilePathFormatChecker) IsFormat(input any) bool {
	return filePathRegex.MatchString(input.(string))
}

var filePathRegex = regexp.MustCompile(`^/[a-z0-9-_/]*$`)

/**********/
/* Domain */
/**********/

type DomainFormatChecker struct{}

func (checker DomainFormatChecker) IsFormat(input any) bool {
	return domainRegex.MatchString(input.(string))
}

var domainRegex = regexp.MustCompile(`^([a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.)+[a-zA-Z]{2,}$`)
