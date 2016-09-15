package ob

import (
	"testing"
)

func BenchmarkParseTokens1(b *testing.B) {
	var phrase = "key"
	b.ReportAllocs()
	l := NewExprLexer()
	for i := 0; i < b.N; i++ {
		l.Parse(phrase)
	}
}
