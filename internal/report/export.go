package report

import (
	"encoding/json"
	"fmt"

	"example.com/energycore/internal/model"
)

func JSONSummary(summary Summary) ([]byte, error) {
	return json.MarshalIndent(summary, "", "  ")
}

func CSVRows(records []model.Record) string {
	output := "id,code,name,status,owner,classification\n"
	for _, record := range records {
		output += fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", record.ID, record.Code, record.Name, FormatStatus(record.Status), record.Owner, record.Classification)
	}
	return output
}
