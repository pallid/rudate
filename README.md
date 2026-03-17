# rudate 📅

[![Go](https://img.shields.io/github/go-mod/go-version/pallid/rudate)](https://github.com/pallid/rudate)
[![Test](https://github.com/pallid/rudate/actions/workflows/ci.yml/badge.svg)](https://github.com/pallid/rudate/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Парсинг дат и времени на русском языке для Go.

Превращает фразы типа «через 5 минут», «в прошлый понедельник», «25 декабря в 19:00» в `time.Time`.

## Установка

```bash
go get github.com/pallid/rudate
```

## Использование

```go
package main

import (
    "fmt"
    "time"

    "github.com/pallid/rudate"
)

func main() {
    now := time.Now()

    // Относительное время
    t, _ := rudate.Parse("5 минут назад", now)
    fmt.Println(t) // 5 минут до now

    t, _ = rudate.Parse("через 2 часа", now)
    fmt.Println(t) // 2 часа после now

    // Даты
    t, _ = rudate.Parse("25 декабря в 19:00", now)
    fmt.Println(t) // 25 Dec at 19:00

    // Комбинации
    t, _ = rudate.Parse("в следующий понедельник в 9 утра", now)
    fmt.Println(t) // Next Monday at 09:00

    // Из произвольного текста
    t, _ = rudate.Parse("напомни мне через час", now)
    fmt.Println(t) // 1 hour from now
}
```

## Поддерживаемые выражения

### Специальные слова
`сейчас`, `сегодня`, `вчера`, `завтра`, `позавчера`, `послезавтра`, `полдень`, `полночь`

### Относительное время
- `5 минут назад`, `час назад`, `3 дня назад`, `2 недели назад`, `месяц назад`, `год назад`
- `через минуту`, `через 2 часа`, `через 3 дня`, `через неделю`, `через месяц`, `через год`

### Коллоквиальные формы
- `полчаса назад`, `через полчаса`
- `полтора часа назад`, `через полтора часа`
- `полгода назад`, `через полгода`
- `четверть часа назад`, `через четверть часа`

### Числительные прописью
`один` ... `двадцать`, `тридцать`, `сорок`, `пятьдесят` + все родовые формы (`одна`, `одну`, `два`, `две`)

### Дни недели
- `в понедельник`, `во вторник`, `в среду`, `в четверг`, `в пятницу`, `в субботу`, `в воскресенье`
- `в прошлый понедельник`, `в прошлую среду`, `в прошлое воскресенье`
- `в следующий вторник`, `в следующую пятницу`

### Даты
- `25 декабря`, `1 января`, `15 марта 2025`
- `25.12.2025`, `01.01`, `15.03`
- `в январе`, `в прошлом марте`, `в следующем июне`

### Время
- `в 15:30`, `в 10:00`, `в 9:30:15`
- `в 10 утра`, `в 5 вечера`, `в 3 дня`, `в 2 ночи`
- `в час дня`, `в час ночи`
- `в полдень`, `в полночь`

### Комбинации
- `вчера в 10 утра`, `завтра в 15:30`
- `в понедельник в 9 утра`, `в следующую среду в 14:00`
- `25 декабря в 19:00`, `25.12.2025 в 18:00`
- `сегодня утром`, `завтра вечером`

### Извлечение из текста
- `напомни через 5 минут`, `напомни мне через час`
- `встреча завтра в 10 утра`

## API

```go
// Парсинг
func Parse(text string, base time.Time, opts ...Option) (time.Time, error)
func MustParse(text string, base time.Time, opts ...Option) time.Time

// Извлечение из текста
func Extract(text string, base time.Time, opts ...Option) (time.Time, int, int, error)
func ExtractAll(text string, base time.Time, opts ...Option) []Match

// Опции
func WithPreferFuture(prefer bool) Option     // default: true
func WithDefaultMorning(hour int) Option       // default: 9
func WithDefaultEvening(hour int) Option       // default: 18
```

## Производительность

```
BenchmarkParse/сейчас                310 ns/op    224 B/op    4 allocs/op
BenchmarkParse/5_минут_назад         609 ns/op    384 B/op    5 allocs/op
BenchmarkParse/через_2_часа          683 ns/op    384 B/op    5 allocs/op
BenchmarkParse/в_прошлый_понедельник 1055 ns/op   384 B/op    5 allocs/op
BenchmarkParse/25_декабря_в_19:00    1060 ns/op   704 B/op    6 allocs/op
BenchmarkParse/вчера_в_10_утра       1036 ns/op   704 B/op    6 allocs/op
BenchmarkParse/полчаса_назад          735 ns/op   384 B/op    5 allocs/op
```

**~300-1000 нс на парс** — в 500,000 раз быстрее вызова LLM.

## Архитектура

```
Входная строка
    ↓
Лексер (dictionary.go) → морфологический словарь, 170+ словоформ → []Token
    ↓
Парсер (parser.go) → рекурсивный спуск → time.Time
```

- **Zero dependencies** — чистый Go, без внешних пакетов
- Морфология инкапсулирована в лексере (все падежи, рода, числа)
- Синтаксис — в парсере (рекурсивный спуск)

## Лицензия

MIT
