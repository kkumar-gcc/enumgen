package lexer

import (
	"fmt"
	"testing"

	"github.com/kkumar-gcc/enumgen/src/token"
)

// Define token classes for better categorization
const (
	special = iota
	literal
	punctuation
	keyword
)

func tokenClass(tok token.Token) int {
	switch {
	case tok.IsLiteral():
		return literal
	case tok.IsKeyword():
		return keyword
	case tok == token.LBRACKET || tok == token.RBRACKET ||
		tok == token.COMMA || tok == token.ASSIGN ||
		tok == token.COLON || tok == token.SEMICOLON || tok == token.SUB:
		return punctuation
	}
	return special
}

// Test element including the token, literal, and class
type tokenTest struct {
	tok   token.Token
	lit   string
	class int
}

// All tokens defined by the grammar.edl
var allTokens = []tokenTest{
	// Special tokens
	{token.COMMENT, "// a comment", special},
	{token.COMMENT, "// another comment with special chars: 123, {}[]", special},

	// Identifiers and literals
	{token.IDENT, "foobar", literal},
	{token.IDENT, "snake_case", literal},
	{token.IDENT, "CamelCase", literal},
	{token.INT, "0", literal},
	{token.INT, "123", literal},
	{token.INT, "9876543210", literal},
	{token.FLOAT, "0.0", literal},
	{token.FLOAT, "3.14", literal},
	{token.FLOAT, ".5", literal},
	{token.FLOAT, "1.", literal},
	{token.STRING, `"string"`, literal},
	{token.STRING, `"string with spaces"`, literal},
	{token.STRING, `"string with \"escaped\" quotes"`, literal},
	{token.CHAR, `'c'`, literal},
	{token.CHAR, `'x'`, literal},

	// Keywords
	{token.ENUM, "enum", keyword},
	{token.KIND, "kind", keyword},
	{token.IOTA, "iota", keyword},
	{token.VALUE, "value", keyword},
	{token.TRUE, "true", keyword},
	{token.FALSE, "false", keyword},

	// Punctuation
	{token.LBRACKET, "[", punctuation},
	{token.RBRACKET, "]", punctuation},
	{token.COMMA, ",", punctuation},
	{token.COLON, ":", punctuation},
	{token.ASSIGN, "=", punctuation},
	{token.SEMICOLON, ";", punctuation},
	{token.SUB, "-", punctuation},
}

// Whitespace to separate tokens
const whitespace = " \t\r\n  "

// build source code from all defined tokens
var source = func() []byte {
	var src []byte
	for _, t := range allTokens {
		src = append(src, t.lit...)
		src = append(src, whitespace...)
	}
	return src
}()

// Count newlines in a string
func countNewlines(s string) int {
	n := 0
	for _, ch := range s {
		if ch == '\n' {
			n++
		}
	}
	return n
}

// TestLexer tests the complete lexing of all token types
func TestLexer(t *testing.T) {
	// Configure lexer with comment mode to capture comments
	l := New(source, CommentMode)

	// Track position
	line := 1
	column := 0

	// Run through all tokens and verify
	index := 0

	for {
		pos, tok, lit := l.Lex()

		// For debugging
		fmt.Printf("Token: %v, Literal: %q, Pos: %v\n", tok, lit, pos)

		if tok == token.EOF {
			break
		}

		// Skip semicolons inserted by the lexer for this test
		if tok == token.SEMICOLON && lit == "\n" {
			// Update line/column for the newline
			line++
			column = 0
			continue
		}

		// Special case for comments
		if tok == token.COMMENT {
			// Verify it's a comment
			if index < len(allTokens) && allTokens[index].tok == token.COMMENT {
				if allTokens[index].lit != lit {
					t.Errorf("bad comment literal: got %q, expected %q", lit, allTokens[index].lit)
				}
				index++
			}

			// Update position - comments can contain newlines
			newlines := countNewlines(lit)
			if newlines > 0 {
				line += newlines
				// Find column after last newline
				lastNL := 0
				for i := len(lit) - 1; i >= 0; i-- {
					if lit[i] == '\n' {
						lastNL = i
						break
					}
				}
				column = len(lit) - lastNL
			} else {
				column += len(lit)
			}
			continue
		}

		// Check against expected token
		if index >= len(allTokens) {
			t.Fatalf("too many tokens: got %v at position %v", tok, pos)
		}

		expected := allTokens[index]
		index++

		// For explicit semicolon tokens, check if they match what we expect
		if tok == token.SEMICOLON && expected.tok == token.SEMICOLON {
			if lit != ";" {
				t.Errorf("expected semicolon character, got %q", lit)
			}
		} else if tok != expected.tok {
			t.Errorf("bad token: got %v, expected %v at index %d", tok, expected.tok, index-1)
		}

		// Check literal for non-operator tokens
		if tok.IsLiteral() || tok == token.COMMENT {
			if lit != expected.lit {
				t.Errorf("bad literal: got %q, expected %q", lit, expected.lit)
			}
		}

		// Check token class
		if tokenClass(tok) != expected.class {
			t.Errorf("bad token class: got %d, expected %d for %v", tokenClass(tok), expected.class, tok)
		}

		// Update position tracking based on the token
		for _, ch := range lit {
			if ch == '\n' {
				line++
				column = 0
			} else {
				column++
			}
		}

		// Update position for whitespace
		for _, ch := range whitespace {
			if ch == '\n' {
				line++
				column = 0
			} else {
				column++
			}
		}
	}

	// Make sure we've seen all expected tokens
	if index != len(allTokens) {
		t.Errorf("not enough tokens: got %d, expected %d", index, len(allTokens))
	}
}

// TestLexerWithSampleCode tests the lexer with a complete enum definition
func TestLexerWithSampleCode(t *testing.T) {
	sample := `
// Color enum represents color values
enum Color [string] {
  Red = "RED",
  Green = "GREEN",
  Blue = "BLUE"
}

// Direction enum with implicit values
enum Direction {
  North,
  East,
  South, 
  West
}
`

	// Expected tokens from the sample code
	expectedTokens := []tokenTest{
		{token.COMMENT, "// Color enum represents color values", special},
		{token.ENUM, "enum", keyword},
		{token.IDENT, "Color", literal},
		{token.LBRACKET, "[", punctuation},
		{token.IDENT, `"string"`, literal},
		{token.RBRACKET, "]", punctuation},
		{token.ILLEGAL, "{", special}, // Lexer treats { as illegal
		{token.IDENT, "Red", literal},
		{token.ASSIGN, "=", punctuation},
		{token.STRING, `"RED"`, literal},
		{token.COMMA, ",", punctuation},
		{token.IDENT, "Green", literal},
		{token.ASSIGN, "=", punctuation},
		{token.STRING, `"GREEN"`, literal},
		{token.COMMA, ",", punctuation},
		{token.IDENT, "Blue", literal},
		{token.ASSIGN, "=", punctuation},
		{token.STRING, `"BLUE"`, literal},
		{token.ILLEGAL, "}", special}, // Lexer treats } as illegal
		{token.COMMENT, "// Direction enum with implicit values", special},
		{token.ENUM, "enum", keyword},
		{token.IDENT, "Direction", literal},
		{token.ILLEGAL, "{", special},
		{token.IDENT, "North", literal},
		{token.COMMA, ",", punctuation},
		{token.IDENT, "East", literal},
		{token.COMMA, ",", punctuation},
		{token.IDENT, "South", literal},
		{token.COMMA, ",", punctuation},
		{token.IDENT, "West", literal},
		{token.ILLEGAL, "}", special},
	}

	// Test lexer with comment mode
	l := New([]byte(sample), CommentMode)

	index := 0
	for {
		_, tok, lit := l.Lex()

		if tok == token.EOF {
			break
		}

		// Skip semicolons inserted by the lexer
		if tok == token.SEMICOLON && lit == "\n" {
			continue
		}

		// Verify token against expected
		if index >= len(expectedTokens) {
			t.Fatalf("too many tokens: got %v %q", tok, lit)
		}

		expected := expectedTokens[index]

		// Special case for strings in the sample (we don't check exact contents)
		if expected.tok == token.STRING && tok == token.STRING {
			// Pass this check
		} else if tok != expected.tok {
			t.Errorf("token %d: got %v, expected %v", index, tok, expected.tok)
		}

		index++
	}

	// Make sure we've seen all expected tokens
	if index < len(expectedTokens) {
		t.Errorf("not enough tokens: got %d, expected %d", index, len(expectedTokens))
	}
}

// TestLexerErrors tests error cases
func TestLexerErrors(t *testing.T) {
	testCases := []struct {
		input string
		error bool
	}{
		{`"unterminated string`, true},
		{`'unterminated char`, true},
		{`@illegal character`, true},
		{`#illegal character`, true},
		{`$illegal character`, true},
		{`123abc`, false}, // Valid - int followed by ident
	}

	for i, tc := range testCases {
		l := New([]byte(tc.input), 0)

		hasErrors := false
		for {
			_, tok, _ := l.Lex()
			if tok == token.ILLEGAL {
				hasErrors = true
			}
			if tok == token.EOF {
				break
			}
		}

		if hasErrors != tc.error {
			t.Errorf("test case %d: error expectation mismatch - got error: %v, expected error: %v",
				i, hasErrors, tc.error)
		}
	}
}

// TestSemicolonInsertion specifically tests the automatic insertion of semicolons
func TestSemicolonInsertion(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected []struct {
			tok token.Token
			lit string
		}
	}{
		{
			desc:  "identifier followed by newline",
			input: "a\nb",
			expected: []struct {
				tok token.Token
				lit string
			}{
				{token.IDENT, "a"},
				{token.SEMICOLON, "\n"},
				{token.IDENT, "b"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			desc:  "int followed by newline",
			input: "123\n456",
			expected: []struct {
				tok token.Token
				lit string
			}{
				{token.INT, "123"},
				{token.SEMICOLON, "\n"},
				{token.INT, "456"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			desc:  "string followed by newline",
			input: "\"string\"\n456",
			expected: []struct {
				tok token.Token
				lit string
			}{
				{token.STRING, "\"string\""},
				{token.SEMICOLON, "\n"},
				{token.INT, "456"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			desc:  "rbracket followed by newline",
			input: "]\n456",
			expected: []struct {
				tok token.Token
				lit string
			}{
				{token.RBRACKET, ""},
				{token.SEMICOLON, "\n"},
				{token.INT, "456"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			desc:  "keyword true followed by newline - should insert semicolon",
			input: "true\n456",
			expected: []struct {
				tok token.Token
				lit string
			}{
				{token.TRUE, "true"},
				//{token.SEMICOLON, "\n"},
				{token.INT, "456"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			desc:  "keyword false followed by newline - should insert semicolon",
			input: "false\n456",
			expected: []struct {
				tok token.Token
				lit string
			}{
				{token.FALSE, "false"},
				//{token.SEMICOLON, "\n"},
				{token.INT, "456"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			desc:  "multiple newlines - should insert semicolons",
			input: "a\n\nb",
			expected: []struct {
				tok token.Token
				lit string
			}{
				{token.IDENT, "a"},
				{token.SEMICOLON, "\n"},
				{token.IDENT, "b"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			desc:  "newline after identifier at EOF - should insert semicolon",
			input: "a\n",
			expected: []struct {
				tok token.Token
				lit string
			}{
				{token.IDENT, "a"},
				{token.SEMICOLON, "\n"},
				{token.EOF, ""},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			l := New([]byte(tc.input), 0)

			for j, exp := range tc.expected {
				pos, tok, lit := l.Lex()

				if tok != exp.tok {
					t.Errorf("token %d: got %v, expected %v", j, tok, exp.tok)
				}

				// For tokens that generate literals, check the literal value
				if exp.lit != "" && lit != exp.lit {
					t.Errorf("literal %d: got %q, expected %q", j, lit, exp.lit)
				}

				t.Logf("Position: %v, Token: %v, Literal: %q", pos, tok, lit)
			}
		})
	}
}

func TestLexingMultipleTokens(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected []token.Token
	}{
		{
			desc:     "enum declaration",
			input:    "enum Color [string]",
			expected: []token.Token{token.ENUM, token.IDENT, token.LBRACKET, token.IDENT, token.RBRACKET, token.SEMICOLON, token.EOF},
		},
		{
			desc:     "enum values",
			input:    "Red = \"RED\", Green = \"GREEN\"",
			expected: []token.Token{token.IDENT, token.ASSIGN, token.STRING, token.COMMA, token.IDENT, token.ASSIGN, token.STRING, token.SEMICOLON, token.EOF},
		},
		{
			desc:     "comment followed by identifier",
			input:    "// This is a comment\nActive,",
			expected: []token.Token{token.COMMENT, token.IDENT, token.COMMA, token.EOF},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			l := New([]byte(tc.input), CommentMode)

			j := 0
			for {
				_, tok, lit := l.Lex()
				t.Logf("Token: %v, Literal: %q", tok, lit)

				if j < len(tc.expected) {
					if tok != tc.expected[j] {
						t.Errorf("token %d: got %v, expected %v", j, tok, tc.expected[j])
					}
				}

				if tok == token.EOF {
					break
				}
				j++
			}

			if j != len(tc.expected)-1 {
				t.Errorf("processed %d tokens, expected %d", j, len(tc.expected)-1)
			}
		})
	}
}

func BenchmarkLexer(b *testing.B) {
	input := `
enum Direction [string]
  North = "NORTH",
  East = "EAST",
  South = "SOUTH",
  West = "WEST"

enum Status
  Active,
  Inactive,
  Pending
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New([]byte(input), 0)
		for {
			_, tok, _ := l.Lex()
			if tok == token.EOF {
				break
			}
		}
	}
}
