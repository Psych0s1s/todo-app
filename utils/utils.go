package utils

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату задачи согласно правилам повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	const layout = "20060102"
	start, err := time.Parse(layout, date)
	if err != nil {
		return "", errors.New("время не может быть преобразовано в корректную дату")
	}

	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	switch repeat[0] {
	case 'y':
		return handleYearly(now, start)
	case 'd':
		return handleDaily(now, start, repeat)
	case 'w':
		return handleWeekly(now, start, repeat)
	case 'm':
		return handleMonthly(now, start, repeat)
	default:
		return "", fmt.Errorf("указан неверный формат: %s", repeat)
	}
}

func handleYearly(now time.Time, start time.Time) (string, error) {

	next := start.AddDate(1, 0, 0)
	for !next.After(now) {
		next = next.AddDate(1, 0, 0)
	}

	return next.Format("20060102"), nil
}

func handleDaily(now, start time.Time, repeat string) (string, error) {
	parts := strings.Split(repeat, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("указан неверный формат: %s", repeat)
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("указан неверный формат: %s", repeat)
	}
	if days < 1 || days > 400 {
		return "", fmt.Errorf("d %d — превышен максимально допустимый интервал", days)
	}
	next := start.AddDate(0, 0, days)
	for !next.After(now) {
		next = next.AddDate(0, 0, days)
	}
	return next.Format("20060102"), nil
}

func handleWeekly(now, start time.Time, repeat string) (string, error) {
	// Убираем 'w' и возможный пробел после него, затем разбиваем по запятым
	repeat = strings.TrimSpace(repeat[1:])
	parts := strings.Split(repeat, ",")
	if len(parts) == 0 {
		return "", fmt.Errorf("указан неверный формат: %s", repeat)
	}

	daysOfWeek := []time.Weekday{}
	for _, part := range parts {
		day, err := strconv.Atoi(strings.TrimSpace(part)) // Преобразуем строку в число
		if err != nil || day < 1 || day > 7 {
			return "", fmt.Errorf("указан неверный формат: %s", repeat)
		}
		// Преобразуем число в соответствующий день недели (1 - понедельник, 7 - воскресенье)
		daysOfWeek = append(daysOfWeek, time.Weekday(day%7)) // %7 чтобы 7 соответствовало воскресенью
	}

	// Сортируем дни недели по возрастанию
	sort.Slice(daysOfWeek, func(i, j int) bool {
		return daysOfWeek[i] < daysOfWeek[j]
	})

	next := findNextWeekday(start, daysOfWeek)
	// Находим следующую подходящую дату, которая больше текущей даты (now)
	for !next.After(now) {
		next = findNextWeekday(next.AddDate(0, 0, 1), daysOfWeek)
	}
	// Возвращаем следующую дату в формате YYYYMMDD
	return next.Format("20060102"), nil
}

func findNextWeekday(start time.Time, daysOfWeek []time.Weekday) time.Time {
	// Перебираем дни недели в списке
	for _, day := range daysOfWeek {
		if start.Weekday() <= day {
			// Если текущий день недели меньше или равен указанному дню, возвращаем эту дату
			return start.AddDate(0, 0, int(day-start.Weekday()))
		}
	}
	// Если все дни в списке меньше текущего дня недели, добавляем 7 дней к самому первому дню в списке
	return start.AddDate(0, 0, int(7-start.Weekday()+daysOfWeek[0]))
}

func handleMonthly(now, start time.Time, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat[1:])
	parts := strings.Split(repeat, " ")
	if len(parts) == 0 || len(parts) > 2 {
		return "", fmt.Errorf("указан неверный формат: %s", repeat)
	}

	// Обрабатываем дни месяца
	daysPart := strings.Split(parts[0], ",")
	daysMap := make(map[int]bool)
	for _, day := range daysPart {
		dayInt, err := strconv.Atoi(day)
		if err != nil || dayInt < -2 || dayInt == 0 || dayInt > 31 {
			return "", fmt.Errorf("указан неверный формат дня месяца: %s", day)
		}
		daysMap[dayInt] = true
	}

	// Обрабатываем месяцы
	monthsMap := make(map[int]bool)
	if len(parts) == 2 {
		for _, m := range strings.Split(parts[1], ",") {
			month, err := strconv.Atoi(m)
			if err != nil || month < 1 || month > 12 {
				return "", fmt.Errorf("указан неверный формат месяца: %s", parts[1])
			}
			monthsMap[month] = true
		}
	} else {
		for i := 1; i <= 12; i++ {
			monthsMap[i] = true
		}
	}

	// Проверяем каждый день, начиная с даты start
	for next := start; ; next = next.AddDate(0, 0, 1) {
		day := next.Day()
		month := int(next.Month())

		// Проверяем положительные и отрицательные дни месяца
		if daysMap[day] || daysMap[day-daysInMonth(next.Month(), next.Year())-1] {
			// Проверяем месяцы
			if monthsMap[month] {
				if next.After(now) {
					return next.Format("20060102"), nil
				}
			}
		}

		// Если мы проходим больше года, это защита от бесконечных циклов
		if next.Year() > now.Year()+1 {
			break
		}
	}

	return "", errors.New("не удалось найти следующую подходящую дату")
}

// определяем число дней в месяце
func daysInMonth(month time.Month, year int) int {
	switch month {
	case time.February:
		if isLeapYear(year) {
			return 29
		}
		return 28
	case time.April, time.June, time.September, time.November:
		return 30
	default:
		return 31
	}
}

// определяем, является ли год високосным
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
