package importer

import (
	"encoding/csv"
	"fmt"
	"strings"

	"example.com/energycore/internal/model"
)

type RowIssue struct {
	Row     int    `json:"row"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateRows(rows []model.ImportRow) []RowIssue {
	issues := []RowIssue{}
	seen := map[string]int{}
	for index, row := range rows {
		rowNumber := index + 2
		if strings.TrimSpace(row.Code) == "" {
			issues = append(issues, RowIssue{Row: rowNumber, Field: "code", Message: "code is required"})
		}
		if strings.TrimSpace(row.Name) == "" {
			issues = append(issues, RowIssue{Row: rowNumber, Field: "name", Message: "name is required"})
		}
		if strings.TrimSpace(row.Owner) == "" {
			issues = append(issues, RowIssue{Row: rowNumber, Field: "owner", Message: "owner is required"})
		}
		code := strings.ToUpper(strings.TrimSpace(row.Code))
		if prior, ok := seen[code]; ok && code != "" {
			issues = append(issues, RowIssue{Row: rowNumber, Field: "code", Message: fmt.Sprintf("duplicate of row %d", prior)})
		}
		if code != "" {
			seen[code] = rowNumber
		}
	}
	return issues
}

func EncodeCSV(rows []model.ImportRow) (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.Write([]string{"code", "name", "description", "owner", "classification"}); err != nil {
		return "", err
	}
	for _, row := range rows {
		if err := writer.Write([]string{row.Code, row.Name, row.Description, row.Owner, row.Classification}); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func CanonicalRows(rows []model.ImportRow) []model.ImportRow {
	canonical := make([]model.ImportRow, 0, len(rows))
	for _, row := range rows {
		canonical = append(canonical, model.ImportRow{Code: strings.ToUpper(model.NormalizeText(row.Code)), Name: model.NormalizeText(row.Name), Description: model.NormalizeText(row.Description), Owner: model.NormalizeText(row.Owner), Classification: strings.ToLower(model.NormalizeText(row.Classification))})
	}
	return canonical
}

func (s *Service) ValidateAndImport(rows []model.ImportRow, actor string) (model.ImportReport, []RowIssue, error) {
	canonical := CanonicalRows(rows)
	issues := ValidateRows(canonical)
	if len(issues) > 0 {
		return model.ImportReport{Total: len(canonical), ImportedIDs: []string{}, Rejected: []string{}}, issues, nil
	}
	report, err := s.ImportRows(canonical, actor)
	return report, issues, err
}
