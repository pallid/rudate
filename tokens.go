package rudate

// TokenType represents the type of a lexical token.
type TokenType int

const (
	TOK_EOF TokenType = iota
	TOK_NUM           // число (целое или дробное)

	// Time units
	TOK_UNIT_SEC   // секунд*
	TOK_UNIT_MIN   // минут*
	TOK_UNIT_HOUR  // час*
	TOK_UNIT_DAY   // день/дн*
	TOK_UNIT_WEEK  // недел*
	TOK_UNIT_MONTH // месяц*
	TOK_UNIT_YEAR  // год/лет

	// Direction
	TOK_DIR_AGO   // назад
	TOK_DIR_IN    // через
	TOK_DIR_AFTER // спустя

	// Prepositions
	TOK_PREP_AT // в/во
	TOK_PREP_ON // на

	// Modifiers
	TOK_MOD_LAST // прошлый/прошлую/прошлое/прошлом/предыдущий
	TOK_MOD_NEXT // следующий/следующую/следующее/следующем
	TOK_MOD_THIS // этот/эту/это/этом

	// Weekday (Value = 0-6 matching time.Weekday: 0=Sunday)
	TOK_WEEKDAY

	// Month (Value = 1-12 matching time.Month)
	TOK_MONTH

	// Special words
	TOK_SPECIAL_NOW           // сейчас
	TOK_SPECIAL_TODAY         // сегодня
	TOK_SPECIAL_YESTERDAY     // вчера
	TOK_SPECIAL_TOMORROW      // завтра
	TOK_SPECIAL_DAYBEFOREYEST // позавчера
	TOK_SPECIAL_DAYAFTERTOM   // послезавтра
	TOK_SPECIAL_NOON          // полдень
	TOK_SPECIAL_MIDNIGHT      // полночь
	TOK_SPECIAL_HALF          // пол- (полчаса, полгода)
	TOK_SPECIAL_QUARTER       // четверть
	TOK_SPECIAL_ONEANDHALF    // полтора/полторы

	// Time of day
	TOK_DAYTIME_MORNING // утра/утром
	TOK_DAYTIME_DAY     // дня/днём
	TOK_DAYTIME_EVENING // вечера/вечером
	TOK_DAYTIME_NIGHT   // ночи/ночью

	// Punctuation
	TOK_COLON // :
	TOK_DOT   // .
	TOK_DASH  // -

	// Unknown
	TOK_WORD // unrecognized word
)

// Token represents a single lexical token.
type Token struct {
	Type  TokenType
	Value float64 // for NUM, WEEKDAY (0-6), MONTH (1-12)
	Raw   string  // original text
	Pos   int     // byte offset in source string
}

// isUnit returns true if the token is a time unit.
func (t Token) isUnit() bool {
	return t.Type >= TOK_UNIT_SEC && t.Type <= TOK_UNIT_YEAR
}
