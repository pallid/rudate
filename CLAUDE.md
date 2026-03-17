# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

**rudate** — pure Go library for parsing Russian natural language date/time expressions. Zero external dependencies. Module path: `github.com/pallid/rudate`.

## Commands

```bash
go test ./...              # run all tests
go test -run TestFuzzy     # run specific test(s) by name
go test -bench .           # run benchmarks
go test -v                 # verbose output
go vet ./...               # static analysis
```

## Architecture

Two-stage pipeline: **Input → Lexer → Token Stream → Parser → time.Time**

| File | Role |
|------|------|
| `rudate.go` | Public API: `Parse`, `MustParse`, `Extract`, `ExtractAll`, `ParseDuration`, options |
| `tokens.go` | Token type enum (37+ variants for morphological/syntactic categories) |
| `dictionary.go` | Maps 170+ Russian word forms → tokens (all cases/genders), compound words (полчаса→30 мин) |
| `lexer.go` | Tokenizer: normalizes input, reads numbers/punctuation/words, combines written tens+ones, optional fuzzy lookup |
| `parser.go` | Recursive descent parser: special words → relative → weekday → date → time → combined |
| `fuzzy.go` | Levenshtein distance with single-row optimization, fuzzy dictionary lookup (disabled by default) |

### Parser grammar (parser.go)

Parses date and time independently, then combines. Productions in priority order:
1. **Special words**: сейчас, сегодня, вчера, завтра, позавчера, послезавтра, полдень, полночь
2. **Relative**: "через N UNIT", "N UNIT назад", "UNIT назад" (implicit 1)
3. **Weekday**: [в] [modifier] WEEKDAY — modifiers: прошлый/следующий/этот
4. **Date**: NUM MONTH [YEAR], DD.MM.YYYY, [в] [modifier] MONTH [DAY]
5. **Time**: [в] (HH:MM:SS | NUM DAYTIME | DAYTIME)

### Key design decisions

- **Morphological dictionary**: all word forms hardcoded for O(1) lookup — no NLP/stemming
- **Compound words**: pre-tokenized as multi-token sequences in dictionary
- **Extraction**: sliding window of up to 12 tokens to find parseable subsequences
- **`preferFuture` default true**: ambiguous dates resolve to future
- **Daytime context**: "дня" is context-sensitive (genitive of "день" vs afternoon indicator)
- **Fuzzy matching opt-in**: `WithFuzzy()` enables Levenshtein ≤ 2 typo correction; short words (≤2 chars) excluded

## Testing

- `rudate_test.go`: core parsing (special words, relative, colloquial, written numbers, weekdays, dates, times, combos)
- `extract_test.go`: Extract/ExtractAll/ParseDuration, custom options
- `fuzzy_test.go`: Levenshtein computation, fuzzy parse on/off

All tests use a fixed `base` time for deterministic results.
