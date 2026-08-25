package importer

import (
	"fmt"
	"strings"

	"example.com/energycore/internal/catalog"
	"example.com/energycore/internal/model"
)

type Service struct {
	catalog *catalog.Service
}

func NewService(catalogue *catalog.Service) *Service {
	return &Service{catalog: catalogue}
}

func (s *Service) ImportCSV(input, actor string) (model.ImportReport, error) {
	rows, err := ParseCSV(input)
	if err != nil {
		return model.ImportReport{}, err
	}
	report := model.ImportReport{Total: len(rows), ImportedIDs: []string{}, Rejected: []string{}}
	for index, row := range rows {
		record, registerErr := s.catalog.Register(model.Record{Code: row.Code, Name: row.Name, Description: row.Description, Owner: row.Owner, Classification: row.Classification}, actor)
		if registerErr != nil {
			report.Rejected = append(report.Rejected, fmt.Sprintf("row %d: %v", index+1, registerErr))
			continue
		}
		report.ImportedIDs = append(report.ImportedIDs, record.ID)
	}
	return report, nil
}

func (s *Service) ImportRows(rows []model.ImportRow, actor string) (model.ImportReport, error) {
	if len(rows) == 0 {
		return model.ImportReport{}, fmt.Errorf("no import rows")
	}
	parts := make([]string, 0, len(rows)+1)
	parts = append(parts, "code,name,description,owner,classification")
	for _, row := range rows {
		parts = append(parts, strings.Join([]string{row.Code, row.Name, row.Description, row.Owner, row.Classification}, ","))
	}
	return s.ImportCSV(strings.Join(parts, "\n"), actor)
}

func Preview(input string) ([]model.ImportRow, []string, error) {
	rows, err := ParseCSV(input)
	if err != nil {
		return nil, nil, err
	}
	valid := make([]model.ImportRow, 0, len(rows))
	issues := []string{}
	for index, row := range rows {
		if strings.TrimSpace(row.Code) == "" || strings.TrimSpace(row.Name) == "" || strings.TrimSpace(row.Owner) == "" {
			issues = append(issues, fmt.Sprintf("row %d is missing identity", index+1))
			continue
		}
		valid = append(valid, row)
	}
	return valid, issues, nil
}
