package processor

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
)

const epsilon = 1e-9

// Record - структура одной записи в JSON Lines
type Record struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Role   string  `json:"role"`
	Salary float64 `json:"salary"`
}

// Stats содержит результаты обработки файла.
type Stats struct {
	Total   int // всего проверенных json-объектов
	Matched int // количество совпадений
	Skipped int // количество объектов без совпадений
}

// Process читает JSON Lines файл и фильтрует записи по полю, записывая статистику в структуру.
// Соответственно: filterField - поле, по которому фильтруем; filterValue - значение, которое
// ищем; maxRecords - максимальное количество записей, которые мы готовы проверить
func Process(in io.Reader, out io.Writer, filterField, filterValue string, maxRecords int) (Stats, error) {
	var stats Stats

	dec := json.NewDecoder(in)
	enc := json.NewEncoder(out)

	for {
		if maxRecords > 0 && stats.Total >= maxRecords {
			break
		}

		var rec Record

		if err := dec.Decode(&rec); err != nil {
			if err == io.EOF {
				break
			}
			return stats, fmt.Errorf("error parsing input file: %w", err)
		}

		stats.Total++

		matched := filterField == "" ||
			matchesFilter(rec, filterField, filterValue)

		if matched {
			stats.Matched++

			if err := enc.Encode(rec); err != nil {
				return stats, fmt.Errorf("encode record: %w", err)
			}

			continue
		}

		stats.Skipped++
	}

	return stats, nil
}

// matchesFilter проверяет, совпадает ли значение поля filterField с filterValue
func matchesFilter(rec Record, field, expected string) bool {
	switch field {
	case "id":
		expectedID, err := strconv.Atoi(expected)
		if err != nil {
			return false
		}

		return rec.ID == expectedID
	case "name":
		return rec.Name == expected
	case "role":
		return rec.Role == expected
	case "salary":
		expectedSalary, err := strconv.ParseFloat(expected, 64)
		if err != nil {
			return false
		}
		return math.Abs(rec.Salary-expectedSalary) < epsilon
	default:
		return false
	}
}
