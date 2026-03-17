package rudate

import "time"

// dictionary maps Russian word forms to their token types and values.
// All morphological forms are listed explicitly for fast lookup.
var dictionary = map[string]Token{
	// === Time units — all forms ===

	// секунда
	"секунда": {Type: TOK_UNIT_SEC}, "секунду": {Type: TOK_UNIT_SEC},
	"секунды": {Type: TOK_UNIT_SEC}, "секунд": {Type: TOK_UNIT_SEC},
	"секунде": {Type: TOK_UNIT_SEC}, "сек": {Type: TOK_UNIT_SEC},

	// минута
	"минута": {Type: TOK_UNIT_MIN}, "минуту": {Type: TOK_UNIT_MIN},
	"минуты": {Type: TOK_UNIT_MIN}, "минут": {Type: TOK_UNIT_MIN},
	"минуте": {Type: TOK_UNIT_MIN}, "мин": {Type: TOK_UNIT_MIN},

	// час
	"час": {Type: TOK_UNIT_HOUR}, "часа": {Type: TOK_UNIT_HOUR},
	"часов": {Type: TOK_UNIT_HOUR}, "часу": {Type: TOK_UNIT_HOUR},
	"часе": {Type: TOK_UNIT_HOUR}, "часик": {Type: TOK_UNIT_HOUR},

	// день
	"день": {Type: TOK_UNIT_DAY}, "дня": {Type: TOK_UNIT_DAY},
	"дней": {Type: TOK_UNIT_DAY}, "дню": {Type: TOK_UNIT_DAY},
	"дне": {Type: TOK_UNIT_DAY},

	// неделя
	"неделя": {Type: TOK_UNIT_WEEK}, "неделю": {Type: TOK_UNIT_WEEK},
	"недели": {Type: TOK_UNIT_WEEK}, "недель": {Type: TOK_UNIT_WEEK},
	"неделе": {Type: TOK_UNIT_WEEK}, "неделек": {Type: TOK_UNIT_WEEK},

	// месяц
	"месяц": {Type: TOK_UNIT_MONTH}, "месяца": {Type: TOK_UNIT_MONTH},
	"месяцев": {Type: TOK_UNIT_MONTH}, "месяце": {Type: TOK_UNIT_MONTH},
	"месяцу": {Type: TOK_UNIT_MONTH},

	// год
	"год": {Type: TOK_UNIT_YEAR}, "года": {Type: TOK_UNIT_YEAR},
	"лет": {Type: TOK_UNIT_YEAR}, "годе": {Type: TOK_UNIT_YEAR},
	"году": {Type: TOK_UNIT_YEAR}, "годик": {Type: TOK_UNIT_YEAR},

	// Abbreviations
	"нед": {Type: TOK_UNIT_WEEK},
	"мес": {Type: TOK_UNIT_MONTH},

	// === Direction ===
	"назад":  {Type: TOK_DIR_AGO},
	"через":  {Type: TOK_DIR_IN},
	"спустя": {Type: TOK_DIR_AFTER},

	// === Prepositions ===
	"в":  {Type: TOK_PREP_AT},
	"во": {Type: TOK_PREP_AT},
	"на": {Type: TOK_PREP_ON},

	// === Modifiers — all gender/case forms ===

	// прошлый
	"прошлый": {Type: TOK_MOD_LAST}, "прошлая": {Type: TOK_MOD_LAST},
	"прошлую": {Type: TOK_MOD_LAST}, "прошлое": {Type: TOK_MOD_LAST},
	"прошлом": {Type: TOK_MOD_LAST}, "прошлой": {Type: TOK_MOD_LAST},
	// предыдущий
	"предыдущий": {Type: TOK_MOD_LAST}, "предыдущая": {Type: TOK_MOD_LAST},
	"предыдущую": {Type: TOK_MOD_LAST}, "предыдущее": {Type: TOK_MOD_LAST},
	"предыдущем": {Type: TOK_MOD_LAST}, "предыдущей": {Type: TOK_MOD_LAST},

	// следующий
	"следующий": {Type: TOK_MOD_NEXT}, "следующая": {Type: TOK_MOD_NEXT},
	"следующую": {Type: TOK_MOD_NEXT}, "следующее": {Type: TOK_MOD_NEXT},
	"следующем": {Type: TOK_MOD_NEXT}, "следующей": {Type: TOK_MOD_NEXT},

	// этот
	"этот": {Type: TOK_MOD_THIS}, "эту": {Type: TOK_MOD_THIS},
	"эта": {Type: TOK_MOD_THIS}, "это": {Type: TOK_MOD_THIS},
	"этом": {Type: TOK_MOD_THIS}, "этой": {Type: TOK_MOD_THIS},

	// === Weekdays (Value = time.Weekday) ===
	"понедельник":  {Type: TOK_WEEKDAY, Value: float64(time.Monday)},
	"понедельника": {Type: TOK_WEEKDAY, Value: float64(time.Monday)},
	"пн":           {Type: TOK_WEEKDAY, Value: float64(time.Monday)},

	"вторник":  {Type: TOK_WEEKDAY, Value: float64(time.Tuesday)},
	"вторника": {Type: TOK_WEEKDAY, Value: float64(time.Tuesday)},
	"вт":       {Type: TOK_WEEKDAY, Value: float64(time.Tuesday)},

	"среда":  {Type: TOK_WEEKDAY, Value: float64(time.Wednesday)},
	"среду":  {Type: TOK_WEEKDAY, Value: float64(time.Wednesday)},
	"среды":  {Type: TOK_WEEKDAY, Value: float64(time.Wednesday)},
	"ср":     {Type: TOK_WEEKDAY, Value: float64(time.Wednesday)},

	"четверг":  {Type: TOK_WEEKDAY, Value: float64(time.Thursday)},
	"четверга": {Type: TOK_WEEKDAY, Value: float64(time.Thursday)},
	"чт":       {Type: TOK_WEEKDAY, Value: float64(time.Thursday)},

	"пятница":  {Type: TOK_WEEKDAY, Value: float64(time.Friday)},
	"пятницу":  {Type: TOK_WEEKDAY, Value: float64(time.Friday)},
	"пятницы":  {Type: TOK_WEEKDAY, Value: float64(time.Friday)},
	"пт":       {Type: TOK_WEEKDAY, Value: float64(time.Friday)},

	"суббота":  {Type: TOK_WEEKDAY, Value: float64(time.Saturday)},
	"субботу":  {Type: TOK_WEEKDAY, Value: float64(time.Saturday)},
	"субботы":  {Type: TOK_WEEKDAY, Value: float64(time.Saturday)},
	"сб":       {Type: TOK_WEEKDAY, Value: float64(time.Saturday)},

	"воскресенье":  {Type: TOK_WEEKDAY, Value: float64(time.Sunday)},
	"воскресенья":  {Type: TOK_WEEKDAY, Value: float64(time.Sunday)},
	"воскресеньем": {Type: TOK_WEEKDAY, Value: float64(time.Sunday)},
	"вс":           {Type: TOK_WEEKDAY, Value: float64(time.Sunday)},

	// === Months (Value = time.Month) ===
	"январь": {Type: TOK_MONTH, Value: 1}, "января": {Type: TOK_MONTH, Value: 1},
	"январе": {Type: TOK_MONTH, Value: 1}, "янв": {Type: TOK_MONTH, Value: 1},

	"февраль": {Type: TOK_MONTH, Value: 2}, "февраля": {Type: TOK_MONTH, Value: 2},
	"феврале": {Type: TOK_MONTH, Value: 2}, "фев": {Type: TOK_MONTH, Value: 2},

	"март": {Type: TOK_MONTH, Value: 3}, "марта": {Type: TOK_MONTH, Value: 3},
	"марте": {Type: TOK_MONTH, Value: 3}, "мар": {Type: TOK_MONTH, Value: 3},

	"апрель": {Type: TOK_MONTH, Value: 4}, "апреля": {Type: TOK_MONTH, Value: 4},
	"апреле": {Type: TOK_MONTH, Value: 4}, "апр": {Type: TOK_MONTH, Value: 4},

	"май": {Type: TOK_MONTH, Value: 5}, "мая": {Type: TOK_MONTH, Value: 5},
	"мае": {Type: TOK_MONTH, Value: 5},

	"июнь": {Type: TOK_MONTH, Value: 6}, "июня": {Type: TOK_MONTH, Value: 6},
	"июне": {Type: TOK_MONTH, Value: 6}, "июн": {Type: TOK_MONTH, Value: 6},

	"июль": {Type: TOK_MONTH, Value: 7}, "июля": {Type: TOK_MONTH, Value: 7},
	"июле": {Type: TOK_MONTH, Value: 7}, "июл": {Type: TOK_MONTH, Value: 7},

	"август": {Type: TOK_MONTH, Value: 8}, "августа": {Type: TOK_MONTH, Value: 8},
	"августе": {Type: TOK_MONTH, Value: 8}, "авг": {Type: TOK_MONTH, Value: 8},

	"сентябрь": {Type: TOK_MONTH, Value: 9}, "сентября": {Type: TOK_MONTH, Value: 9},
	"сентябре": {Type: TOK_MONTH, Value: 9}, "сен": {Type: TOK_MONTH, Value: 9},

	"октябрь": {Type: TOK_MONTH, Value: 10}, "октября": {Type: TOK_MONTH, Value: 10},
	"октябре": {Type: TOK_MONTH, Value: 10}, "окт": {Type: TOK_MONTH, Value: 10},

	"ноябрь": {Type: TOK_MONTH, Value: 11}, "ноября": {Type: TOK_MONTH, Value: 11},
	"ноябре": {Type: TOK_MONTH, Value: 11}, "ноя": {Type: TOK_MONTH, Value: 11},

	"декабрь": {Type: TOK_MONTH, Value: 12}, "декабря": {Type: TOK_MONTH, Value: 12},
	"декабре": {Type: TOK_MONTH, Value: 12}, "дек": {Type: TOK_MONTH, Value: 12},

	// === Special words ===
	"сейчас":      {Type: TOK_SPECIAL_NOW},
	"сегодня":     {Type: TOK_SPECIAL_TODAY},
	"вчера":       {Type: TOK_SPECIAL_YESTERDAY},
	"завтра":      {Type: TOK_SPECIAL_TOMORROW},
	"позавчера":   {Type: TOK_SPECIAL_DAYBEFOREYEST},
	"послезавтра": {Type: TOK_SPECIAL_DAYAFTERTOM},
	"полдень":     {Type: TOK_SPECIAL_NOON},
	"полночь":     {Type: TOK_SPECIAL_MIDNIGHT},
	"четверть":    {Type: TOK_SPECIAL_QUARTER},
	"полтора":     {Type: TOK_SPECIAL_ONEANDHALF, Value: 1.5},
	"полторы":     {Type: TOK_SPECIAL_ONEANDHALF, Value: 1.5},

	// === Time of day ===
	"утра":    {Type: TOK_DAYTIME_MORNING},
	"утром":   {Type: TOK_DAYTIME_MORNING},
	// NOTE: "дня" is NOT here — it's listed under TOK_UNIT_DAY.
	// Disambiguation happens in the parser: after NUM + HOUR → daytime; after NUM → unit.
	"днём":    {Type: TOK_DAYTIME_DAY},
	"днем":    {Type: TOK_DAYTIME_DAY},
	"вечера":  {Type: TOK_DAYTIME_EVENING},
	"вечером": {Type: TOK_DAYTIME_EVENING},
	"ночи":    {Type: TOK_DAYTIME_NIGHT},
	"ночью":   {Type: TOK_DAYTIME_NIGHT},

	// === Written numbers ===
	"ноль":         {Type: TOK_NUM, Value: 0},
	"один":         {Type: TOK_NUM, Value: 1},
	"одна":         {Type: TOK_NUM, Value: 1},
	"одну":         {Type: TOK_NUM, Value: 1},
	"одно":         {Type: TOK_NUM, Value: 1},
	"два":          {Type: TOK_NUM, Value: 2},
	"две":          {Type: TOK_NUM, Value: 2},
	"три":          {Type: TOK_NUM, Value: 3},
	"четыре":       {Type: TOK_NUM, Value: 4},
	"пять":         {Type: TOK_NUM, Value: 5},
	"шесть":        {Type: TOK_NUM, Value: 6},
	"семь":         {Type: TOK_NUM, Value: 7},
	"восемь":       {Type: TOK_NUM, Value: 8},
	"девять":       {Type: TOK_NUM, Value: 9},
	"десять":       {Type: TOK_NUM, Value: 10},
	"одиннадцать":  {Type: TOK_NUM, Value: 11},
	"двенадцать":   {Type: TOK_NUM, Value: 12},
	"тринадцать":   {Type: TOK_NUM, Value: 13},
	"четырнадцать": {Type: TOK_NUM, Value: 14},
	"пятнадцать":   {Type: TOK_NUM, Value: 15},
	"шестнадцать":  {Type: TOK_NUM, Value: 16},
	"семнадцать":   {Type: TOK_NUM, Value: 17},
	"восемнадцать": {Type: TOK_NUM, Value: 18},
	"девятнадцать": {Type: TOK_NUM, Value: 19},
	"двадцать":     {Type: TOK_NUM, Value: 20},
	"тридцать":     {Type: TOK_NUM, Value: 30},
	"сорок":        {Type: TOK_NUM, Value: 40},
	"пятьдесят":    {Type: TOK_NUM, Value: 50},
}

// compoundWords maps compound Russian words to token sequences.
// These are checked before dictionary lookup.
var compoundWords = map[string][]Token{
	"полчаса": {
		{Type: TOK_NUM, Value: 30},
		{Type: TOK_UNIT_MIN},
	},
	"полгода": {
		{Type: TOK_NUM, Value: 6},
		{Type: TOK_UNIT_MONTH},
	},
	"полдня": {
		{Type: TOK_NUM, Value: 12},
		{Type: TOK_UNIT_HOUR},
	},
}

// daytimeTokens lists token types that are time-of-day indicators.
// Used for disambiguation (e.g., "дня" can be genitive of "день" or "afternoon").
var daytimeTokens = map[TokenType]bool{
	TOK_DAYTIME_MORNING: true,
	TOK_DAYTIME_DAY:     true,
	TOK_DAYTIME_EVENING: true,
	TOK_DAYTIME_NIGHT:   true,
}
