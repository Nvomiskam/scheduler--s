// Пакет rules предоставляет логику вычисления дат выполнения для повторяющихся задач.
package rules

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// NextDate вычисляет следующую дату выполнения задачи.
//
// Параметры:
//   - now: текущее время, относительно которого вычисляется следующая дата
//   - dateStr: начальная дата в формате DateFormat
//   - repeat:   правило повторения в формате:
//   - "y"       - ежегодно
//   - "d N"     - каждые N дней (1 ≤ N ≤ 400)
//   - "w 1,3,5" - указанные дни недели (1-пн, 7-вс)
//   - "m d1,d2" - указанные дни месяца:
//   - 1-31    - конкретные числа
//   - -1      - последний день месяца
//   - -2      - предпоследний день
//   - Опционально: "m 1,15 1,6" (1 и 15 число января и июня)
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("не указано правило")
	}

	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return "", errors.New("некорректное исходное время")
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "y":
		return handleYearlyRepeat(now, date)
	case "d":
		return handleDailyRepeat(now, date, parts)
	case "w":
		return handleWeeklyRepeat(now, date, parts)
	case "m":
		return handleMonthlyRepeat(now, date, parts)
	default:
		return "", errors.New("неизвестное правило")
	}
}

// handleMonthlyRepeat обрабатывает повторение по дням месяца
func handleMonthlyRepeat(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("не указаны дни месяца")
	}

	monthDays := make(map[int]bool)
	for _, dayStr := range strings.Split(parts[1], ",") {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < -2 || day > 31 || day == 0 {
			return "", errors.New("неверный формат дней месяца")
		}
		monthDays[day] = true
	}

	months := make(map[int]bool)
	if len(parts) > 2 {
		for _, monthStr := range strings.Split(parts[2], ",") {
			month, err := strconv.Atoi(monthStr)
			if err != nil || month < 1 || month > 12 {
				return "", errors.New("неверный формат месяцев")
			}
			months[month] = true
		}
	} else {
		for i := 1; i <= 12; i++ {
			months[i] = true
		}
	}

	next := date
	for {
		next = next.AddDate(0, 0, 1)
		currentMonth := int(next.Month())
		currentDay := next.Day()
		lastDay := lastDayOfMonth(next)

		if !months[currentMonth] {
			continue
		}

		dayMatches := false

		for d := range monthDays {
			if d > 0 && currentDay == d {
				dayMatches = true
				break
			}
		}

		if monthDays[-1] && currentDay == lastDay {
			dayMatches = true
		}
		if monthDays[-2] && currentDay == lastDay-1 {
			dayMatches = true
		}

		if dayMatches && (next.After(now) || next.Equal(now)) {
			return next.Format(DateFormat), nil
		}
	}
}

// handleWeeklyRepeat обрабатывает повторение по дням недели
func handleWeeklyRepeat(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("не указаны дни недели")
	}

	weekdays := make(map[int]bool)
	for _, dayStr := range strings.Split(parts[1], ",") {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day > 7 {
			return "", errors.New("неверный формат дней недели")
		}
		weekdays[day] = true
	}

	start := now
	if date.After(now) {
		start = date
	}
	next := start.AddDate(0, 0, 1)

	for i := 0; i < 400; i++ {
		weekday := int(next.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		if weekdays[weekday] && !next.Before(now) {
			return next.Format(DateFormat), nil
		}

		next = next.AddDate(0, 0, 1)
	}

	return "", errors.New("не удалось найти подходящую дату")
}

// handleYearlyRepeat обрабатывает ежегодное повторение
func handleYearlyRepeat(now, date time.Time) (string, error) {
	next := date

	for {
		next = next.AddDate(1, 0, 0)

		// Корректировка для 29 февраля в невисокосный год
		if next.Month() == time.February && next.Day() == 29 && !isLeapYear(next.Year()) {
			next = time.Date(next.Year(), time.March, 1, 0, 0, 0, 0, next.Location())
		}

		if next.After(now) {
			break
		}
	}

	return next.Format(DateFormat), nil
}

// handleDailyRepeat обрабатывает повторение с интервалом
func handleDailyRepeat(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("не указан интервал в днях")
	}

	days, err := strconv.Atoi(parts[1])
	if err != nil || days <= 0 {
		return "", errors.New("неверный формат правила")
	}

	if days > 400 {
		return "", errors.New("превышен максимально допустимый интервал")
	}

	next := date
	for {
		next = next.AddDate(0, 0, days)
		if next.After(now) {
			return next.Format(DateFormat), nil
		}
	}
}

// isLeapYear проверяет - является ли год високосным
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// isLeapYear возвращает последний день месяца
func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
