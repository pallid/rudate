// Package rudate provides natural language date/time parsing for Russian.
//
// It parses human-readable Russian phrases like "через 5 минут",
// "в прошлый понедельник", "25 декабря в 19:00" into time.Time values.
//
// The parser uses a two-stage approach:
//  1. Lexer with a morphological dictionary normalizes Russian word forms into tokens
//  2. Recursive descent parser processes the token stream into time values
//
// Usage:
//
//	t, err := rudate.Parse("через 5 минут", time.Now())
//	t, err := rudate.Parse("вчера в 10 утра", time.Now())
//	t, err := rudate.Parse("25 декабря в 19:00", time.Now())
package rudate

import (
	"fmt"
	"time"
)

// Option configures the parser behavior.
type Option func(*options)

// WithPreferFuture sets whether ambiguous dates should resolve to the future.
// Default: true.
func WithPreferFuture(prefer bool) Option {
	return func(o *options) {
		o.preferFuture = prefer
	}
}

// WithDefaultMorning sets the default hour for "утром" (morning). Default: 9.
func WithDefaultMorning(hour int) Option {
	return func(o *options) {
		o.defaultMorning = hour
	}
}

// WithDefaultDay sets the default hour for "днём" (afternoon). Default: 14.
func WithDefaultDay(hour int) Option {
	return func(o *options) {
		o.defaultDay = hour
	}
}

// WithDefaultEvening sets the default hour for "вечером" (evening). Default: 18.
func WithDefaultEvening(hour int) Option {
	return func(o *options) {
		o.defaultEvening = hour
	}
}

// WithDefaultNight sets the default hour for "ночью" (night). Default: 23.
func WithDefaultNight(hour int) Option {
	return func(o *options) {
		o.defaultNight = hour
	}
}

// WithFuzzy enables fuzzy matching for typos (Levenshtein distance ≤ 2).
// This allows parsing inputs like "чирез 5 минтут" as "через 5 минут".
// Disabled by default for maximum performance.
func WithFuzzy() Option {
	return func(o *options) {
		o.fuzzy = true
	}
}

// Parse parses a Russian natural language date/time string relative to base.
func Parse(text string, base time.Time, opts ...Option) (time.Time, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}

	lexer := NewLexer(text)
	lexer.fuzzyEnabled = o.fuzzy
	tokens := lexer.Tokenize()

	parser := newParser(tokens, base, o)
	return parser.Parse()
}

// MustParse is like Parse but panics on error.
func MustParse(text string, base time.Time, opts ...Option) time.Time {
	t, err := Parse(text, base, opts...)
	if err != nil {
		panic(err)
	}
	return t
}

// Extract finds the first date/time mention in arbitrary text.
// Returns the parsed time, start and end byte offsets, and an error.
func Extract(text string, base time.Time, opts ...Option) (time.Time, int, int, error) {
	matches := ExtractAll(text, base, opts...)
	if len(matches) == 0 {
		return time.Time{}, 0, 0, fmt.Errorf("rudate: дата/время не найдена в тексте")
	}
	m := matches[0]
	return m.Time, m.Start, m.End, nil
}

// Match describes a date/time found in text.
type Match struct {
	Time  time.Time
	Start int    // byte offset of match start
	End   int    // byte offset of match end
	Text  string // original matched substring
}

// ExtractAll finds all date/time mentions in arbitrary text.
// It uses a sliding window over tokens to find parseable date/time subsequences.
func ExtractAll(text string, base time.Time, opts ...Option) []Match {
	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}

	lexer := NewLexer(text)
	tokens := lexer.Tokenize()

	var matches []Match
	i := 0
	for i < len(tokens) {
		if tokens[i].Type == TOK_EOF {
			break
		}
		// Skip plain words that can't start a date expression
		if tokens[i].Type == TOK_WORD {
			i++
			continue
		}

		// Try parsing from position i with increasing window sizes
		bestLen := 0
		var bestTime time.Time
		var bestStart, bestEnd int

		for windowEnd := i + 1; windowEnd <= len(tokens) && windowEnd <= i+12; windowEnd++ {
			// Build a sub-slice ending with EOF
			sub := make([]Token, windowEnd-i+1)
			copy(sub, tokens[i:windowEnd])
			sub[len(sub)-1] = Token{Type: TOK_EOF}

			p := newParser(sub, base, o)
			t, err := p.Parse()
			if err == nil {
				// Use lastUsedPos to determine actual consumed range
				consumed := p.lastUsedPos
				if consumed < 1 {
					consumed = 1
				}
				if consumed > bestLen {
					bestLen = consumed
					bestTime = t
					bestStart = tokens[i].Pos
					// End at the last actually used token
					lastIdx := i + consumed - 1
					if lastIdx >= len(tokens) {
						lastIdx = len(tokens) - 1
					}
					lastTok := tokens[lastIdx]
					bestEnd = lastTok.Pos + len(lastTok.Raw)
				}
			}
		}

		if bestLen > 0 {
			matches = append(matches, Match{
				Time:  bestTime,
				Start: bestStart,
				End:   bestEnd,
				Text:  text[bestStart:bestEnd],
			})
			i += bestLen
		} else {
			i++
		}
	}

	return matches
}

// ParseDuration parses a Russian duration expression and returns time.Duration.
// Examples: "5 минут", "полтора часа", "2 часа 30 минут"
func ParseDuration(text string, base time.Time, opts ...Option) (time.Duration, error) {
	t, err := Parse(text+" назад", base, opts...)
	if err != nil {
		return 0, fmt.Errorf("rudate: не удалось распознать длительность: %w", err)
	}
	return base.Sub(t), nil
}
