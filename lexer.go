package rudate

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lexer tokenizes Russian natural language date/time strings.
type Lexer struct {
	input        string
	pos          int
	tokens       []Token
	fuzzyEnabled bool // enable fuzzy matching for typos
}

// NewLexer creates a new Lexer for the given input string.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: strings.ToLower(strings.TrimSpace(input)),
	}
}

// Tokenize processes the input and returns a slice of tokens.
func (l *Lexer) Tokenize() []Token {
	l.tokens = nil

	for l.pos < len(l.input) {
		// Skip whitespace
		if l.skipWhitespace() {
			continue
		}

		r, size := utf8.DecodeRuneInString(l.input[l.pos:])

		switch {
		case r == ':':
			l.tokens = append(l.tokens, Token{Type: TOK_COLON, Raw: ":", Pos: l.pos})
			l.pos += size
		case r == '.':
			l.tokens = append(l.tokens, Token{Type: TOK_DOT, Raw: ".", Pos: l.pos})
			l.pos += size
		case r == '-':
			// Check if it's a negative number or ordinal suffix like "1-го"
			l.tokens = append(l.tokens, Token{Type: TOK_DASH, Raw: "-", Pos: l.pos})
			l.pos += size
		case unicode.IsDigit(r):
			l.readNumber()
		case unicode.IsLetter(r):
			l.readWord()
		default:
			// Skip unknown characters
			l.pos += size
		}
	}

	l.tokens = append(l.tokens, Token{Type: TOK_EOF, Pos: l.pos})
	return l.tokens
}

// skipWhitespace advances past whitespace, returns true if any was skipped.
func (l *Lexer) skipWhitespace() bool {
	start := l.pos
	for l.pos < len(l.input) {
		r, size := utf8.DecodeRuneInString(l.input[l.pos:])
		if !unicode.IsSpace(r) {
			break
		}
		l.pos += size
	}
	return l.pos > start
}

// readNumber reads a numeric token (integer).
func (l *Lexer) readNumber() {
	start := l.pos
	for l.pos < len(l.input) {
		r, size := utf8.DecodeRuneInString(l.input[l.pos:])
		if !unicode.IsDigit(r) {
			break
		}
		l.pos += size
	}
	raw := l.input[start:l.pos]
	val, _ := strconv.ParseFloat(raw, 64)
	l.tokens = append(l.tokens, Token{Type: TOK_NUM, Value: val, Raw: raw, Pos: start})

	// Skip ordinal suffixes like "-го", "-е", "-й"
	if l.pos < len(l.input) && l.input[l.pos] == '-' {
		savedPos := l.pos
		l.pos++ // skip '-'
		// Try to read suffix
		suffStart := l.pos
		for l.pos < len(l.input) {
			r, size := utf8.DecodeRuneInString(l.input[l.pos:])
			if !unicode.IsLetter(r) {
				break
			}
			l.pos += size
		}
		suffix := l.input[suffStart:l.pos]
		// Known ordinal suffixes
		switch suffix {
		case "го", "е", "й", "я", "му", "м":
			// Valid ordinal suffix, consumed
		default:
			// Not an ordinal suffix, revert
			l.pos = savedPos
		}
	}
}

// readWord reads a word token and looks it up in dictionaries.
func (l *Lexer) readWord() {
	start := l.pos
	for l.pos < len(l.input) {
		r, size := utf8.DecodeRuneInString(l.input[l.pos:])
		if !unicode.IsLetter(r) {
			break
		}
		l.pos += size
	}

	word := l.input[start:l.pos]

	// Check compound words first
	if tokens, ok := compoundWords[word]; ok {
		for _, t := range tokens {
			t.Raw = word
			t.Pos = start
			l.tokens = append(l.tokens, t)
		}
		return
	}

	// Look up in dictionary
	if tok, ok := dictionary[word]; ok {
		tok.Raw = word
		tok.Pos = start
		l.tokens = append(l.tokens, tok)
		return
	}

	// Additive number words: check if this is a tens+ones compound
	// e.g., after "двадцать" we might see "один" → 21
	if len(l.tokens) > 0 {
		last := l.tokens[len(l.tokens)-1]
		if last.Type == TOK_NUM && isRoundTens(last.Value) {
			if tok, ok := dictionary[word]; ok && tok.Type == TOK_NUM && tok.Value < 10 {
				// Combine: двадцать + один = 21
				l.tokens[len(l.tokens)-1].Value += tok.Value
				l.tokens[len(l.tokens)-1].Raw += " " + word
				return
			}
		}
	}

	// Fuzzy matching: try to find closest dictionary word
	if l.fuzzyEnabled {
		if tok, ok := fuzzyLookup(word, 2); ok {
			tok.Pos = start
			l.tokens = append(l.tokens, tok)
			return
		}
	}

	// Unknown word — skip it (allows extracting dates from arbitrary text)
	l.tokens = append(l.tokens, Token{Type: TOK_WORD, Raw: word, Pos: start})
}

// isRoundTens checks if a number is a round tens value (20, 30, 40, 50).
func isRoundTens(v float64) bool {
	return v == 20 || v == 30 || v == 40 || v == 50
}
