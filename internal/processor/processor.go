package processor

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Record - структура одной записи в JSON Lines
type Record struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Role   string  `json:"role"`
	Salary float32 `json:"salary"`
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
func Process(inputPath, filterField, filterValue string, maxRecords int) (Stats, error) {
	var stats Stats

	file, err := os.Open(inputPath)
	if err != nil {
		return stats, fmt.Errorf("error opening input file: %w", err)
	}

	defer file.Close()

	dec := json.NewDecoder(file)

	for {
		if maxRecords > 0 && stats.Total >= maxRecords {
			break
		}

		var rec Record

		if err = dec.Decode(&rec); err != nil {
			if err == io.EOF {
				break
			}
			return stats, fmt.Errorf("error parsing input file: %w", err)
		}

		stats.Total++

		if filterField == "" || filterValue == "" {
			stats.Matched++
			printRecord(rec)
			continue
		}

		if matchesFilter(rec, filterField, filterValue) {
			stats.Matched++
			printRecord(rec)
		} else {
			stats.Skipped++
		}
	}

	return stats, nil
}

// matchesFilter проверяет, совпадает ли значение поля filterField с filterValue
func matchesFilter(rec Record, field, expected string) bool {
	switch field {
	case "id":
		return fmt.Sprint(rec.ID) == expected
	case "name":
		return rec.Name == expected
	case "role":
		return rec.Role == expected
	case "salary":
		expectedFloat, err := strconv.ParseFloat(expected, 32)
		if err != nil {
			return false
		}
		return rec.Salary == float32(expectedFloat)
	default:
		return false
	}
}

// printRecord вспомогательная функция для дебага
func printRecord(r Record) {
	data, err := json.Marshal(r)
	if err != nil {
		log.Printf("marshal record: %v", err)
		return
	}
	fmt.Println(string(data))
}
