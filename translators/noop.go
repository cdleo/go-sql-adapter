package translators

import adapter "go-sql-adapter"

type noopTranslator struct{}

func NewNoopTranslator() adapter.SQLSyntaxTranslator {
	return &noopTranslator{}
}

func (t *noopTranslator) Translate(query string) string {
	return query
}
