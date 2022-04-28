package validator

import (
	"regexp"
)

/************/
/* Git Repo */
/************/

type gitRepoFormatChecker struct{}

var gitRepoRegex = regexp.MustCompile(`((git|ssh|http(s)?)|(git@[\w.]+))(:(//)?)([\w.@:/\-~]+)(\.git)(/)?`)

func (checker gitRepoFormatChecker) IsFormat(input interface{}) bool {
	return gitRepoRegex.MatchString(input.(string))
}

/*************/
/* File Path */
/*************/

type filePathFormatChecker struct{}

var filePathRegex = regexp.MustCompile(`^/[a-z0-9-_/]*$`)

func (checker filePathFormatChecker) IsFormat(input interface{}) bool {
	return filePathRegex.MatchString(input.(string))
}

/**********/
/* Domain */
/**********/

type domainFormatChecker struct{}

var domainRegex = regexp.MustCompile(`^([a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.)+[a-zA-Z]{2,}$`)

func (checker domainFormatChecker) IsFormat(input interface{}) bool {
	return domainRegex.MatchString(input.(string))
}
