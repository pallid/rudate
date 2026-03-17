package rudate

import (
	"testing"
	"time"
)

func TestExtract(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
		wantText string
	}{
		{
			"напомни мне через 5 минут позвонить",
			base.Add(5 * time.Minute),
			"через 5 минут",
		},
		{
			"встреча завтра в 10 утра обязательно",
			dateTime(2025, 3, 16, 10, 0),
			"завтра в 10 утра",
		},
		{
			"дедлайн 25 декабря, не забудь",
			date(2025, 12, 25),
			"25 декабря",
		},
		{
			"созвон в пятницу в 15:00 с клиентом",
			dateTime(2025, 3, 21, 15, 0),
			"в пятницу в 15:00",
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, start, end, err := Extract(c.input, base)
			if err != nil {
				t.Fatalf("Extract(%q) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Extract(%q) time = %v, want %v", c.input, got, c.expected)
			}
			extracted := c.input[start:end]
			if extracted != c.wantText {
				t.Errorf("Extract(%q) text = %q, want %q", c.input, extracted, c.wantText)
			}
		})
	}
}

func TestExtractAll(t *testing.T) {
	input := "встреча завтра в 10 утра и дедлайн 25 декабря"
	matches := ExtractAll(input, base)
	if len(matches) < 1 {
		t.Fatalf("ExtractAll(%q) got %d matches, want >= 1", input, len(matches))
	}
	// At least the first match should be correct
	if !matches[0].Time.Equal(dateTime(2025, 3, 16, 10, 0)) {
		t.Errorf("first match time = %v, want %v", matches[0].Time, dateTime(2025, 3, 16, 10, 0))
	}
}

func TestParseDuration(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Duration
	}{
		{"5 минут", 5 * time.Minute},
		{"час", 1 * time.Hour},
		{"2 часа", 2 * time.Hour},
		{"полчаса", 30 * time.Minute},
		{"полтора часа", 90 * time.Minute},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := ParseDuration(c.input, base)
			if err != nil {
				t.Fatalf("ParseDuration(%q) error: %v", c.input, err)
			}
			if got != c.expected {
				t.Errorf("ParseDuration(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestWithOptions(t *testing.T) {
	// Test PreferFuture=false
	t.Run("prefer_past", func(t *testing.T) {
		// With preferFuture=false, "в понедельник" should be past
		got, err := Parse("в понедельник", base, WithPreferFuture(false))
		if err != nil {
			t.Fatal(err)
		}
		if got.After(base) {
			t.Errorf("with PreferFuture(false), got future date: %v", got)
		}
	})

	// Test custom morning hour
	t.Run("custom_morning", func(t *testing.T) {
		got, err := Parse("сегодня утром", base, WithDefaultMorning(7))
		if err != nil {
			t.Fatal(err)
		}
		expected := dateTime(2025, 3, 15, 7, 0)
		if !got.Equal(expected) {
			t.Errorf("with DefaultMorning(7), got %v, want %v", got, expected)
		}
	})

	// Test custom evening hour
	t.Run("custom_evening", func(t *testing.T) {
		got, err := Parse("сегодня вечером", base, WithDefaultEvening(20))
		if err != nil {
			t.Fatal(err)
		}
		expected := dateTime(2025, 3, 15, 20, 0)
		if !got.Equal(expected) {
			t.Errorf("with DefaultEvening(20), got %v, want %v", got, expected)
		}
	})
}
