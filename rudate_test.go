package rudate

import (
	"testing"
	"time"
)

// base time: Saturday, March 15, 2025, 14:30:00 MSK
var msk = time.FixedZone("MSK", 3*60*60)
var base = time.Date(2025, 3, 15, 14, 30, 0, 0, msk)

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, msk)
}

func dateTime(y int, m time.Month, d, h, min int) time.Time {
	return time.Date(y, m, d, h, min, 0, 0, msk)
}

func TestSpecial(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"сейчас", base},
		{"сегодня", date(2025, 3, 15)},
		{"вчера", date(2025, 3, 14)},
		{"завтра", date(2025, 3, 16)},
		{"позавчера", date(2025, 3, 13)},
		{"послезавтра", date(2025, 3, 17)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestRelativeAgo(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"минуту назад", base.Add(-1 * time.Minute)},
		{"5 минут назад", base.Add(-5 * time.Minute)},
		{"час назад", base.Add(-1 * time.Hour)},
		{"2 часа назад", base.Add(-2 * time.Hour)},
		{"день назад", base.AddDate(0, 0, -1)},
		{"3 дня назад", base.AddDate(0, 0, -3)},
		{"неделю назад", base.AddDate(0, 0, -7)},
		{"2 недели назад", base.AddDate(0, 0, -14)},
		{"месяц назад", base.AddDate(0, -1, 0)},
		{"год назад", base.AddDate(-1, 0, 0)},
		{"5 лет назад", base.AddDate(-5, 0, 0)},
		{"10 секунд назад", base.Add(-10 * time.Second)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestRelativeFuture(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"через минуту", base.Add(1 * time.Minute)},
		{"через 5 минут", base.Add(5 * time.Minute)},
		{"через час", base.Add(1 * time.Hour)},
		{"через 2 часа", base.Add(2 * time.Hour)},
		{"через день", base.AddDate(0, 0, 1)},
		{"через 3 дня", base.AddDate(0, 0, 3)},
		{"через неделю", base.AddDate(0, 0, 7)},
		{"через месяц", base.AddDate(0, 1, 0)},
		{"через год", base.AddDate(1, 0, 0)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestColloquial(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"полчаса назад", base.Add(-30 * time.Minute)},
		{"через полчаса", base.Add(30 * time.Minute)},
		{"полтора часа назад", base.Add(-90 * time.Minute)},
		{"через полтора часа", base.Add(90 * time.Minute)},
		{"полгода назад", base.AddDate(0, -6, 0)},
		{"через полгода", base.AddDate(0, 6, 0)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestWrittenNumbers(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"пять минут назад", base.Add(-5 * time.Minute)},
		{"через два часа", base.Add(2 * time.Hour)},
		{"через три дня", base.AddDate(0, 0, 3)},
		{"одну минуту назад", base.Add(-1 * time.Minute)},
		{"десять секунд назад", base.Add(-10 * time.Second)},
		{"двадцать минут назад", base.Add(-20 * time.Minute)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestWeekdays(t *testing.T) {
	// Base is Saturday March 15, 2025
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"в понедельник", date(2025, 3, 17)},
		{"в среду", date(2025, 3, 19)},
		{"в пятницу", date(2025, 3, 21)},
		{"в прошлый понедельник", date(2025, 3, 10)},
		{"в прошлую среду", date(2025, 3, 12)},
		{"в прошлую пятницу", date(2025, 3, 14)},
		{"в следующий понедельник", date(2025, 3, 17)},
		{"в следующую среду", date(2025, 3, 19)},
		{"в следующее воскресенье", date(2025, 3, 16)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestDates(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"25 декабря", date(2025, 12, 25)},
		{"1 января", date(2026, 1, 1)},   // already passed → next year
		{"15 марта", date(2025, 3, 15)},   // today
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestTime(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"в 15:30", dateTime(2025, 3, 15, 15, 30)},
		{"в 10:00", dateTime(2025, 3, 15, 10, 0)},
		{"в 10 утра", dateTime(2025, 3, 15, 10, 0)},
		{"в 5 вечера", dateTime(2025, 3, 15, 17, 0)},
		{"в 3 дня", dateTime(2025, 3, 15, 15, 0)},
		{"в час дня", dateTime(2025, 3, 15, 13, 0)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestCombined(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"вчера в 10 утра", dateTime(2025, 3, 14, 10, 0)},
		{"вчера в 15:30", dateTime(2025, 3, 14, 15, 30)},
		{"завтра в 9 утра", dateTime(2025, 3, 16, 9, 0)},
		{"завтра в полдень", dateTime(2025, 3, 16, 12, 0)},
		{"сегодня утром", dateTime(2025, 3, 15, 9, 0)},
		{"сегодня вечером", dateTime(2025, 3, 15, 18, 0)},
		{"завтра утром", dateTime(2025, 3, 16, 9, 0)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestEmbeddedText(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"напомни через 5 минут", base.Add(5 * time.Minute)},
		{"напомни мне через час", base.Add(1 * time.Hour)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestQuarterHour(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"четверть часа назад", base.Add(-15 * time.Minute)},
		{"через четверть часа", base.Add(15 * time.Minute)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestDateDotFormat(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"25.12.2025", date(2025, 12, 25)},
		{"01.01.2026", date(2026, 1, 1)},
		{"15.03.2025", date(2025, 3, 15)},
		{"25.12", date(2025, 12, 25)},          // this year (future)
		{"01.01", date(2026, 1, 1)},             // already passed → next year
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestAbbreviations(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"5 мин назад", base.Add(-5 * time.Minute)},
		{"через 3 часа", base.Add(3 * time.Hour)},
		{"10 сек назад", base.Add(-10 * time.Second)},
		{"через 1 нед", base.AddDate(0, 0, 7)},
		{"2 мес назад", base.AddDate(0, -2, 0)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestWeekdayTimeCombos(t *testing.T) {
	// Base is Saturday March 15, 2025
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"в понедельник в 9 утра", dateTime(2025, 3, 17, 9, 0)},
		{"в следующую среду в 14:00", dateTime(2025, 3, 19, 14, 0)},
		{"в прошлую пятницу в 17:30", dateTime(2025, 3, 14, 17, 30)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestDateTimeCombos(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"25 декабря в 19:00", dateTime(2025, 12, 25, 19, 0)},
		{"1 января в полночь", dateTime(2026, 1, 1, 0, 0)},
		{"25.12.2025 в 18:00", dateTime(2025, 12, 25, 18, 0)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	benchmarks := []string{
		"сейчас",
		"5 минут назад",
		"через 2 часа",
		"в прошлый понедельник",
		"25 декабря в 19:00",
		"вчера в 10 утра",
		"полчаса назад",
	}

	for _, input := range benchmarks {
		b.Run(input, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = Parse(input, base)
			}
		})
	}
}
