package translators

import (
	"regexp"
	"strings"

	adapter "github.com/cdleo/go-sql-adapter"
)

type postgresTranslator struct {
	paramRegExp     *regexp.Regexp
	sourceSQLSyntax adapter.SQLSyntax
}

func NewPostgresTranslator(sourceSQLSyntax adapter.SQLSyntax) adapter.SQLSyntaxTranslator {
	return &postgresTranslator{
		regexp.MustCompile(":[1-9]"),
		sourceSQLSyntax,
	}
}

func (s *postgresTranslator) Translate(query string) string {

	if s.sourceSQLSyntax == adapter.SQLSyntax_Oracle {
		return s.paramRegExp.ReplaceAllStringFunc(query, func(m string) string {
			return strings.Replace(m, ":", "$", 1)
		})
	} else {
		return query
	}

}
