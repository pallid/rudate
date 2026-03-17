package main

import (
	"fmt"
	"time"

	"github.com/pallid/rudate"
)

func main() {
	now := time.Now()

	phrases := []string{
		// Специальные слова
		"сейчас",
		"вчера",
		"послезавтра",

		// Относительные
		"5 минут назад",
		"через 2 часа",
		"полчаса назад",
		"через полтора часа",
		"четверть часа назад",
		"полгода назад",

		// Числительные прописью
		"пять минут назад",
		"через двадцать минут",

		// Дни недели
		"в понедельник",
		"в прошлую среду",
		"в следующую пятницу",

		// Даты
		"25 декабря",
		"1 января 2026",
		"25.12.2025",

		// Время
		"в 15:30",
		"в 10 утра",
		"в 5 вечера",
		"в час дня",

		// Комбинации
		"вчера в 10 утра",
		"завтра в 15:30",
		"в понедельник в 9 утра",
		"25 декабря в 19:00",
		"сегодня вечером",

		// Из текста
		"напомни мне через 5 минут",
	}

	fmt.Printf("Текущее время: %s\n\n", now.Format("2006-01-02 15:04:05"))

	for _, phrase := range phrases {
		t, err := rudate.Parse(phrase, now)
		if err != nil {
			fmt.Printf("  %-40s → ОШИБКА: %v\n", phrase, err)
			continue
		}
		fmt.Printf("  %-40s → %s\n", phrase, t.Format("2006-01-02 15:04:05"))
	}

	// Пример Extract
	fmt.Println("\n--- Extract ---")
	text := "напомни мне через 5 минут позвонить маме"
	t, start, end, err := rudate.Extract(text, now)
	if err == nil {
		fmt.Printf("  Текст: %q\n", text)
		fmt.Printf("  Найдено: %q → %s\n", text[start:end], t.Format("2006-01-02 15:04:05"))
	}

	// Пример ParseDuration
	fmt.Println("\n--- ParseDuration ---")
	d, err := rudate.ParseDuration("полтора часа", now)
	if err == nil {
		fmt.Printf("  полтора часа = %v\n", d)
	}
}
