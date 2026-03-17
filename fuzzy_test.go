package rudate

import (
	"testing"
	"time"
)

func TestLevenshtein(t *testing.T) {
	cases := []struct {
		a, b     string
		expected int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", 1},
		{"через", "через", 0},
		{"чирез", "через", 1},
		{"минут", "минтут", 1},
		{"минут", "минут", 0},
		{"завтра", "завтар", 2},
		{"вчера", "вчира", 1},
		{"понедельник", "панедельник", 1},
		{"секунд", "сикунд", 1},
		{"abc", "xyz", 3},
	}

	for _, c := range cases {
		t.Run(c.a+"_"+c.b, func(t *testing.T) {
			got := levenshtein(c.a, c.b)
			if got != c.expected {
				t.Errorf("levenshtein(%q, %q) = %d, want %d", c.a, c.b, got, c.expected)
			}
		})
	}
}

func TestFuzzyParse(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		// Distance 1
		{"чирез 5 минут", base.Add(5 * time.Minute)},
		{"через 5 минтут", base.Add(5 * time.Minute)},
		{"вчира", date(2025, 3, 14)},
		{"завтар", date(2025, 3, 16)},
		{"сийчас", base},
		{"час нащад", base.Add(-1 * time.Hour)},

		// Distance 2
		{"чиерз 5 минут", base.Add(5 * time.Minute)},
		{"позовчера", date(2025, 3, 13)},

		// Combined with correct words
		{"чирез час", base.Add(1 * time.Hour)},
		{"минтуу назад", base.Add(-1 * time.Minute)},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input, base, WithFuzzy())
			if err != nil {
				t.Fatalf("Parse(%q, WithFuzzy()) error: %v", c.input, err)
			}
			if !got.Equal(c.expected) {
				t.Errorf("Parse(%q, WithFuzzy()) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}

func TestFuzzyDisabledByDefault(t *testing.T) {
	// Without fuzzy, a fully garbled input should fail
	_, err := Parse("чирез минтут нащад", base)
	if err == nil {
		t.Error("Parse with all typos should fail without WithFuzzy()")
	}
}

func BenchmarkFuzzyParse(b *testing.B) {
	b.Run("fuzzy_off", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = Parse("через 5 минут", base)
		}
	})
	b.Run("fuzzy_on_no_typo", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = Parse("через 5 минут", base, WithFuzzy())
		}
	})
	b.Run("fuzzy_on_with_typo", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = Parse("чирез 5 минтут", base, WithFuzzy())
		}
	})
}
