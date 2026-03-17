package rudate

import (
	"fmt"
	"time"
)

// Parser implements a recursive descent parser for Russian date/time expressions.
type Parser struct {
	tokens      []Token
	pos         int
	base        time.Time
	result      time.Time
	opts        options
	lastUsedPos int // track last token position that contributed to result
}

type options struct {
	preferFuture   bool
	defaultMorning int // default hour for "утром"
	defaultDay     int // default hour for "днём"
	defaultEvening int // default hour for "вечером"
	defaultNight   int // default hour for "ночью"
	fuzzy          bool // enable fuzzy matching for typos
}

func defaultOptions() options {
	return options{
		preferFuture:   true,
		defaultMorning: 9,
		defaultDay:     14,
		defaultEvening: 18,
		defaultNight:   23,
	}
}

// newParser creates a new Parser.
func newParser(tokens []Token, base time.Time, opts options) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		base:   base,
		result: base,
		opts:   opts,
	}
}

// current returns the current token.
func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TOK_EOF}
	}
	return p.tokens[p.pos]
}

// peek returns the token at offset from current position.
func (p *Parser) peek(offset int) Token {
	idx := p.pos + offset
	if idx >= len(p.tokens) || idx < 0 {
		return Token{Type: TOK_EOF}
	}
	return p.tokens[idx]
}

// advance moves to the next token.
func (p *Parser) advance() Token {
	t := p.current()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return t
}

// expect consumes the current token if it matches the given type.
func (p *Parser) expect(tt TokenType) (Token, bool) {
	if p.current().Type == tt {
		return p.advance(), true
	}
	return Token{}, false
}

// skipWords skips any TOK_WORD tokens (noise words like "мне", "меня", etc.)
func (p *Parser) skipWords() {
	for p.current().Type == TOK_WORD {
		p.advance()
	}
}

// Parse attempts to parse the token stream into a time.Time.
func (p *Parser) Parse() (time.Time, error) {
	// Try to parse date part
	dateSet := false
	timeSet := false

	for p.current().Type != TOK_EOF {
		p.skipWords()
		if p.current().Type == TOK_EOF {
			break
		}

		parsed := false

		// Try date productions (only if date not yet set)
		if !parsed && !dateSet {
			if t, ok := p.trySpecial(); ok {
				p.result = t
				dateSet = true
				parsed = true
				p.lastUsedPos = p.pos
			}
		}

		if !parsed && !dateSet {
			if t, ok := p.tryRelative(); ok {
				p.result = t
				dateSet = true
				parsed = true
				p.lastUsedPos = p.pos
			}
		}

		if !parsed && !dateSet {
			if t, ok := p.tryWeekday(); ok {
				p.result = t
				dateSet = true
				parsed = true
				p.lastUsedPos = p.pos
			}
		}

		if !parsed && !dateSet {
			if t, ok := p.tryDateDotFormat(); ok {
				p.result = t
				dateSet = true
				parsed = true
				p.lastUsedPos = p.pos
			}
		}

		if !parsed && !dateSet {
			if t, ok := p.tryDate(); ok {
				p.result = t
				dateSet = true
				parsed = true
				p.lastUsedPos = p.pos
			}
		}

		// Try time (always, even after date is set — for combos like "вчера в 10 утра")
		if !parsed && !timeSet {
			if t, ok := p.tryTime(); ok {
				// Apply time to current result (keep date part)
				year, month, day := p.result.Date()
				hour, min, sec := t.Clock()
				p.result = time.Date(year, month, day, hour, min, sec, 0, p.base.Location())
				timeSet = true
				parsed = true
				p.lastUsedPos = p.pos
			}
		}

		if !parsed {
			// Skip unrecognized token
			p.advance()
		}
	}

	if !dateSet && !timeSet {
		return time.Time{}, fmt.Errorf("rudate: не удалось распознать дату/время")
	}

	return p.result, nil
}

// trySpecial tries to parse special words: сейчас, сегодня, вчера, завтра, позавчера, послезавтра
func (p *Parser) trySpecial() (time.Time, bool) {
	cur := p.current()
	switch cur.Type {
	case TOK_SPECIAL_NOW:
		p.advance()
		return p.base, true
	case TOK_SPECIAL_TODAY:
		p.advance()
		return truncateDay(p.base), true
	case TOK_SPECIAL_YESTERDAY:
		p.advance()
		return truncateDay(p.base.AddDate(0, 0, -1)), true
	case TOK_SPECIAL_TOMORROW:
		p.advance()
		return truncateDay(p.base.AddDate(0, 0, 1)), true
	case TOK_SPECIAL_DAYBEFOREYEST:
		p.advance()
		return truncateDay(p.base.AddDate(0, 0, -2)), true
	case TOK_SPECIAL_DAYAFTERTOM:
		p.advance()
		return truncateDay(p.base.AddDate(0, 0, 2)), true
	case TOK_SPECIAL_NOON:
		p.advance()
		d := truncateDay(p.base)
		return d.Add(12 * time.Hour), true
	case TOK_SPECIAL_MIDNIGHT:
		p.advance()
		// Next midnight
		d := truncateDay(p.base.AddDate(0, 0, 1))
		return d, true
	}
	return time.Time{}, false
}

// tryRelative tries to parse relative expressions:
//   - "через N единиц" (future)
//   - "N единиц назад" (past)
//   - "полчаса назад", "полтора часа назад"
//   - "четверть часа назад"
func (p *Parser) tryRelative() (time.Time, bool) {
	saved := p.pos

	// Pattern 1: "через" NUM UNIT
	if p.current().Type == TOK_DIR_IN {
		p.advance()
		p.skipWords()

		amount, unit, ok := p.readAmountUnit()
		if ok {
			return p.applyRelative(p.base, amount, unit, 1), true
		}
		p.pos = saved
		return time.Time{}, false
	}

	// Pattern 2: NUM UNIT "назад"
	// Pattern 3: UNIT "назад" (implicit 1)
	// Pattern 4: "полтора" UNIT "назад"
	// Pattern 5: "четверть" UNIT "назад"
	amount, unit, ok := p.readAmountUnit()
	if ok {
		p.skipWords()
		if p.current().Type == TOK_DIR_AGO {
			p.advance()
			return p.applyRelative(p.base, amount, unit, -1), true
		}
		if p.current().Type == TOK_DIR_AFTER {
			p.advance()
			return p.applyRelative(p.base, amount, unit, 1), true
		}
		// No direction word — could be implicit. Revert.
		p.pos = saved
		return time.Time{}, false
	}

	p.pos = saved
	return time.Time{}, false
}

// readAmountUnit reads a number (or special like полтора, четверть) followed by a time unit.
// Returns (amount, unitToken, ok).
func (p *Parser) readAmountUnit() (float64, TokenType, bool) {
	saved := p.pos

	var amount float64 = 1

	cur := p.current()
	switch {
	case cur.Type == TOK_NUM:
		amount = cur.Value
		p.advance()
	case cur.Type == TOK_SPECIAL_ONEANDHALF:
		amount = 1.5
		p.advance()
	case cur.Type == TOK_SPECIAL_QUARTER:
		amount = 0.25
		p.advance()
	case cur.isUnit():
		// Implicit 1, e.g., "минуту назад", "час назад"
		amount = 1
		// Don't advance — the unit check below will handle it
	default:
		p.pos = saved
		return 0, 0, false
	}

	p.skipWords()

	// Now expect a time unit
	if p.current().isUnit() {
		unit := p.current().Type
		p.advance()
		return amount, unit, true
	}

	// Special case: amount was read but no unit follows — might be a standalone number
	p.pos = saved
	return 0, 0, false
}

// applyRelative applies a relative offset to base time.
func (p *Parser) applyRelative(base time.Time, amount float64, unit TokenType, direction int) time.Time {
	switch unit {
	case TOK_UNIT_SEC:
		return base.Add(time.Duration(direction) * time.Duration(amount*float64(time.Second)))
	case TOK_UNIT_MIN:
		return base.Add(time.Duration(direction) * time.Duration(amount*float64(time.Minute)))
	case TOK_UNIT_HOUR:
		return base.Add(time.Duration(direction) * time.Duration(amount*float64(time.Hour)))
	case TOK_UNIT_DAY:
		days := int(amount)
		if amount != float64(days) {
			// Fractional days
			return base.Add(time.Duration(direction) * time.Duration(amount*24*float64(time.Hour)))
		}
		return base.AddDate(0, 0, direction*days)
	case TOK_UNIT_WEEK:
		weeks := int(amount)
		if amount != float64(weeks) {
			return base.Add(time.Duration(direction) * time.Duration(amount*7*24*float64(time.Hour)))
		}
		return base.AddDate(0, 0, direction*weeks*7)
	case TOK_UNIT_MONTH:
		months := int(amount)
		if amount != float64(months) {
			// Approximate fractional months as 30 days
			return base.Add(time.Duration(direction) * time.Duration(amount*30*24*float64(time.Hour)))
		}
		return base.AddDate(0, direction*months, 0)
	case TOK_UNIT_YEAR:
		years := int(amount)
		if amount != float64(years) {
			return base.Add(time.Duration(direction) * time.Duration(amount*365*24*float64(time.Hour)))
		}
		return base.AddDate(direction*years, 0, 0)
	}
	return base
}

// tryWeekday tries to parse weekday expressions:
//   - "понедельник", "в понедельник"
//   - "в прошлый понедельник"
//   - "в следующую среду"
func (p *Parser) tryWeekday() (time.Time, bool) {
	saved := p.pos

	// Optional "в"/"во"
	if p.current().Type == TOK_PREP_AT {
		p.advance()
	}

	// Optional modifier
	mod := TOK_EOF
	if p.current().Type == TOK_MOD_LAST || p.current().Type == TOK_MOD_NEXT || p.current().Type == TOK_MOD_THIS {
		mod = p.current().Type
		p.advance()
	}

	// Weekday
	if p.current().Type == TOK_WEEKDAY {
		weekday := time.Weekday(int(p.current().Value))
		p.advance()

		switch mod {
		case TOK_MOD_LAST:
			return truncateDay(prevWeekday(p.base, weekday)), true
		case TOK_MOD_NEXT:
			return truncateDay(nextWeekday(p.base, weekday)), true
		default:
			// Default: prefer future (next occurrence)
			if p.opts.preferFuture {
				return truncateDay(nextWeekday(p.base, weekday)), true
			}
			return truncateDay(prevWeekday(p.base, weekday)), true
		}
	}

	p.pos = saved
	return time.Time{}, false
}

// tryDate tries to parse date expressions:
//   - "25 декабря", "1 января 2025"
//   - "в январе", "в прошлом марте", "в следующем июне"
func (p *Parser) tryDate() (time.Time, bool) {
	saved := p.pos

	// Pattern 1: NUM MONTH (YEAR)?
	if p.current().Type == TOK_NUM {
		day := int(p.current().Value)
		if day >= 1 && day <= 31 {
			p.advance()
			p.skipWords()
			if p.current().Type == TOK_MONTH {
				month := time.Month(int(p.current().Value))
				p.advance()

				year := p.base.Year()
				// Check for year
				if p.current().Type == TOK_NUM && p.current().Value > 31 {
					year = int(p.current().Value)
					p.advance()
				} else if p.opts.preferFuture {
					// If date already passed this year, use next year
					candidate := time.Date(year, month, day, 0, 0, 0, 0, p.base.Location())
					if candidate.Before(truncateDay(p.base)) {
						year++
					}
				}

				return time.Date(year, month, day, 0, 0, 0, 0, p.base.Location()), true
			}
		}
		p.pos = saved
	}

	// Pattern 2: "в" (modifier)? MONTH
	if p.current().Type == TOK_PREP_AT {
		p.advance()

		mod := TOK_EOF
		if p.current().Type == TOK_MOD_LAST || p.current().Type == TOK_MOD_NEXT {
			mod = p.current().Type
			p.advance()
		}

		if p.current().Type == TOK_MONTH {
			month := time.Month(int(p.current().Value))
			p.advance()

			year := p.base.Year()
			_, curMonth, curDay := p.base.Date()
			hour, min, sec := p.base.Clock()

			switch mod {
			case TOK_MOD_LAST:
				if month >= curMonth {
					year--
				}
			case TOK_MOD_NEXT:
				if month <= curMonth {
					year++
				}
			default:
				// Nearest future month
				if p.opts.preferFuture {
					if month < curMonth || (month == curMonth && curDay > 1) {
						year++
					}
				}
			}

			// Check for day number after month
			if p.current().Type == TOK_NUM {
				d := int(p.current().Value)
				if d >= 1 && d <= 31 {
					p.advance()
					return time.Date(year, month, d, 0, 0, 0, 0, p.base.Location()), true
				}
			}

			return time.Date(year, month, 1, hour, min, sec, 0, p.base.Location()), true
		}

		p.pos = saved
		return time.Time{}, false
	}

	// Pattern 3: bare MONTH (without "в")
	if p.current().Type == TOK_MONTH {
		month := time.Month(int(p.current().Value))
		p.advance()

		year := p.base.Year()
		if p.opts.preferFuture {
			_, curMonth, _ := p.base.Date()
			if month <= curMonth {
				year++
			}
		}

		hour, min, sec := p.base.Clock()
		return time.Date(year, month, 1, hour, min, sec, 0, p.base.Location()), true
	}

	p.pos = saved
	return time.Time{}, false
}

// tryTime tries to parse time expressions:
//   - "в 15:30", "15:30"
//   - "в 10 утра", "в 5 вечера"
//   - "в час дня"
//   - "в полдень", "в полночь"
func (p *Parser) tryTime() (time.Time, bool) {
	saved := p.pos

	// Optional "в"
	if p.current().Type == TOK_PREP_AT {
		p.advance()
	}

	// "полдень" / "полночь"
	if p.current().Type == TOK_SPECIAL_NOON {
		p.advance()
		d := truncateDay(p.result)
		return d.Add(12 * time.Hour), true
	}
	if p.current().Type == TOK_SPECIAL_MIDNIGHT {
		p.advance()
		return truncateDay(p.result.AddDate(0, 0, 1)), true
	}

	// "час" as shorthand for 1
	if p.current().Type == TOK_UNIT_HOUR {
		p.advance()
		// "час дня" = 13:00, "час ночи" = 01:00
		if dt, ok := p.tryDaytime(); ok {
			hour := applyDaytime(1, dt)
			d := truncateDay(p.result)
			return d.Add(time.Duration(hour) * time.Hour), true
		}
		p.pos = saved
		return time.Time{}, false
	}

	// NUM followed by ":" or daytime
	if p.current().Type == TOK_NUM {
		hour := int(p.current().Value)
		if hour < 0 || hour > 23 {
			p.pos = saved
			return time.Time{}, false
		}
		p.advance()

		min := 0
		sec := 0

		// Check for HH:MM
		if p.current().Type == TOK_COLON {
			p.advance()
			if p.current().Type == TOK_NUM {
				min = int(p.current().Value)
				p.advance()

				// Check for HH:MM:SS
				if p.current().Type == TOK_COLON {
					p.advance()
					if p.current().Type == TOK_NUM {
						sec = int(p.current().Value)
						p.advance()
					}
				}
			}
		}

		// Check for "часов" (optional word)
		if p.current().Type == TOK_UNIT_HOUR {
			p.advance()
		}

		// Check for daytime modifier
		if dt, ok := p.tryDaytime(); ok {
			hour = applyDaytime(hour, dt)
		}

		if min >= 0 && min <= 59 && sec >= 0 && sec <= 59 {
			d := truncateDay(p.result)
			return d.Add(time.Duration(hour)*time.Hour +
				time.Duration(min)*time.Minute +
				time.Duration(sec)*time.Second), true
		}
	}

	// Daytime without specific hour: "утром", "вечером"
	if dt, ok := p.tryDaytime(); ok {
		d := truncateDay(p.result)
		var hour int
		switch dt {
		case TOK_DAYTIME_MORNING:
			hour = p.opts.defaultMorning
		case TOK_DAYTIME_DAY:
			hour = p.opts.defaultDay
		case TOK_DAYTIME_EVENING:
			hour = p.opts.defaultEvening
		case TOK_DAYTIME_NIGHT:
			hour = p.opts.defaultNight
		}
		return d.Add(time.Duration(hour) * time.Hour), true
	}

	p.pos = saved
	return time.Time{}, false
}

// tryDaytime tries to parse a daytime token (утра, дня, вечера, ночи).
// Special handling for "дня" which is ambiguous (genitive of "день" vs time of day).
// In time context (after a number + hour), "дня" means afternoon.
func (p *Parser) tryDaytime() (TokenType, bool) {
	if daytimeTokens[p.current().Type] {
		tt := p.current().Type
		p.advance()
		return tt, true
	}
	// "дня" is mapped to TOK_UNIT_DAY in dictionary, but in time context it means afternoon
	if p.current().Type == TOK_UNIT_DAY && p.current().Raw == "дня" {
		p.advance()
		return TOK_DAYTIME_DAY, true
	}
	return 0, false
}

// applyDaytime adjusts hour based on daytime indicator.
func applyDaytime(hour int, dt TokenType) int {
	switch dt {
	case TOK_DAYTIME_MORNING:
		// 1-11 утра → as is
		return hour
	case TOK_DAYTIME_DAY:
		// 1-4 дня → +12, 12 дня → 12
		if hour >= 1 && hour <= 4 {
			return hour + 12
		}
		return hour
	case TOK_DAYTIME_EVENING:
		// 5-11 вечера → +12
		if hour >= 1 && hour <= 11 {
			return hour + 12
		}
		return hour
	case TOK_DAYTIME_NIGHT:
		// 1-3 ночи → as is (01:00-03:00)
		// 12 ночи → 00:00
		if hour == 12 {
			return 0
		}
		return hour
	}
	return hour
}

// tryDateDotFormat tries to parse dates in DD.MM.YYYY or DD.MM format.
func (p *Parser) tryDateDotFormat() (time.Time, bool) {
	saved := p.pos

	// Need: NUM DOT NUM (DOT NUM)?
	if p.current().Type != TOK_NUM {
		return time.Time{}, false
	}
	day := int(p.current().Value)
	if day < 1 || day > 31 {
		return time.Time{}, false
	}
	p.advance()

	if p.current().Type != TOK_DOT {
		p.pos = saved
		return time.Time{}, false
	}
	p.advance()

	if p.current().Type != TOK_NUM {
		p.pos = saved
		return time.Time{}, false
	}
	month := int(p.current().Value)
	if month < 1 || month > 12 {
		p.pos = saved
		return time.Time{}, false
	}
	p.advance()

	year := p.base.Year()
	if p.current().Type == TOK_DOT {
		p.advance()
		if p.current().Type == TOK_NUM {
			y := int(p.current().Value)
			if y > 100 { // full year like 2025
				year = y
			} else { // two-digit year like 25
				year = 2000 + y
			}
			p.advance()
		}
	} else if p.opts.preferFuture {
		candidate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, p.base.Location())
		if candidate.Before(truncateDay(p.base)) {
			year++
		}
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, p.base.Location()), true
}

// === Helper functions ===

// truncateDay returns a date truncated to the start of the day.
func truncateDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// prevWeekday returns the previous occurrence of a weekday.
func prevWeekday(t time.Time, day time.Weekday) time.Time {
	d := int(t.Weekday()) - int(day)
	if d <= 0 {
		d += 7
	}
	return t.AddDate(0, 0, -d)
}

// nextWeekday returns the next occurrence of a weekday.
func nextWeekday(t time.Time, day time.Weekday) time.Time {
	d := int(day) - int(t.Weekday())
	if d <= 0 {
		d += 7
	}
	return t.AddDate(0, 0, d)
}
