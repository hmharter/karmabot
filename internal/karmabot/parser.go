package karmabot

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"go.uber.org/zap"
)

func ParseCallouts(logger *zap.SugaredLogger, s string) []string {
	s = strings.ReplaceAll(s, "\u00a0", " ")
	regex := regexp.MustCompile("<?@[0-9A-Za-z]+>? ?(\\++|\\-+)")
	callouts := regex.FindAllString(s, -1)
	logger.Debugf("callouts: %v", callouts)

	return callouts
}

func ParseUserID(logger *zap.SugaredLogger, s string) string {
	regex := regexp.MustCompile("@[0-9A-Za-z]+")
	user := regex.FindString(s)
	logger.Debugf("user: %v", user)
	_, i := utf8.DecodeRuneInString(user)

	return user[i:]
}

func ParseKarma(logger *zap.SugaredLogger, s string) string {
	regex := regexp.MustCompile("(\\++|\\-+)")
	karma := regex.FindString(s)
	logger.Debugf("karma: %v", karma)
	return karma
}

func ParseSlashFlags(s string) []Flag {
	regex := regexp.MustCompile("--(?P<key>[a-z]+)=\\\"(?P<val>.+)\\\"")
	flags := regex.FindAllStringSubmatch(s, -1)

	fs := []Flag{}

	// This is totally overkill because there is only one flag right now, but you never know.
	for _, flag := range flags {
		result := make(map[string]string)
		for i, name := range regex.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = flag[i]
			}
		}
		f := &Flag{
			Key:   result["key"],
			Value: result["val"],
		}
		fs = append(fs, *f)
	}

	return fs
}
