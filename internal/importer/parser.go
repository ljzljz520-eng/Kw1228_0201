package importer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	"example.com/energycore/internal/model"
)

func ParseCSV(input string) ([]model.ImportRow, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.FieldsPerRecord = -1
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read csv header: %w", err)
	}
	indexes, err := mapHeaders(headers)
	if err != nil {
		return nil, err
	}
	rows := []model.ImportRow{}
	for {
		fields, readErr := reader.Read()
		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			return nil, fmt.Errorf("read csv row: %w", readErr)
		}
		if len(fields) == 0 || strings.TrimSpace(strings.Join(fields, "")) == "" {
			continue
		}
		row := model.ImportRow{Code: field(fields, indexes, "code"), Name: field(fields, indexes, "name"), Description: field(fields, indexes, "description"), Owner: field(fields, indexes, "owner"), Classification: field(fields, indexes, "classification")}
		rows = append(rows, row)
	}
	return rows, nil
}

func mapHeaders(headers []string) (map[string]int, error) {
	indexes := map[string]int{}
	for index, header := range headers {
		key := strings.ToLower(strings.TrimSpace(header))
		indexes[key] = index
	}
	for _, required := range []string{"code", "name", "owner"} {
		if _, ok := indexes[required]; !ok {
			return nil, fmt.Errorf("missing csv column %s", required)
		}
	}
	return indexes, nil
}

func field(fields []string, indexes map[string]int, name string) string {
	index, ok := indexes[name]
	if !ok || index >= len(fields) {
		return ""
	}
	return strings.TrimSpace(fields[index])
}
