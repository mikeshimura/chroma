package lexers

import (
	. "github.com/alecthomas/chroma" // nolint
)

// Markdown lexer.
var Markdown = Register(MustNewLexer(
	&Config{
		Name:      "markdown",
		Aliases:   []string{"md"},
		Filenames: []string{"*.md"},
		MimeTypes: []string{"text/x-markdown"},
	},
	Rules{
		"root": {
			{`^(#)([^#].+\n)`, ByGroups(GenericHeading, Text), nil},
			{`^(#{2,6})(.+\n)`, ByGroups(GenericSubheading, Text), nil},
			{`^(\s*)([*-] )(\[[ xX]\])( .+\n)`, ByGroups(Text, Keyword, Keyword, UsingSelf("inline")), nil},
			{`^(\s*)([*-])(\s)(.+\n)`, ByGroups(Text, Keyword, Text, UsingSelf("inline")), nil},
			{`^(\s*)([0-9]+\.)( .+\n)`, ByGroups(Text, Keyword, UsingSelf("inline")), nil},
			{`^(\s*>\s)(.+\n)`, ByGroups(Keyword, GenericEmph), nil},
			{"^(```\\n)([\\w\\W]*?)(^```$)", ByGroups(LiteralString, Text, LiteralString), nil},
			{"^(```)(\\w+)(\\n)([\\w\\W]*?)(^```$)", EmitterFunc(handleCodeblock), nil},
			Include("inline"),
		},
		"inline": {
			{`\\.`, Text, nil},
			{`(\s)([*_][^*_]+[*_])(\W|\n)`, ByGroups(Text, GenericEmph, Text), nil},
			{`(\s)((\*\*|__).*\3)((?=\W|\n))`, ByGroups(Text, GenericStrong, None, Text), nil},
			{`(\s)(~~[^~]+~~)((?=\W|\n))`, ByGroups(Text, GenericDeleted, Text), nil},
			{"`[^`]+`", LiteralStringBacktick, nil},
			{`[@#][\w/:]+`, NameEntity, nil},
			{`(!?\[)([^]]+)(\])(\()([^)]+)(\))`, ByGroups(Text, NameTag, Text, Text, NameAttribute, Text), nil},
			{`[^\\\s]+`, Text, nil},
			{`.`, Text, nil},
		},
	},
))

func handleCodeblock(groups []string, lexer Lexer, out func(*Token)) {
	out(&Token{String, groups[1]})
	out(&Token{String, groups[2]})
	out(&Token{Text, groups[3]})
	code := groups[4]
	lexer = Get(groups[2])
	if lexer == nil {
		out(&Token{String, code})
	} else {
		lexer.Tokenise(nil, code, out)
	}
	out(&Token{String, groups[5]})
}
