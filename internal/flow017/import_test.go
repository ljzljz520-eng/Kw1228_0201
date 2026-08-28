package flow017

import "testing"

func TestWorkflowImportReport(t *testing.T) {
	runtime := testRuntime(t)
	input := "code,name,description,owner,classification\nSC-101,Orion Core,Stable plasma transfer module,North Lab,beta\nSC-102,,Missing name,North Lab,beta"
	result, summary, err := runtime.ImportReport(input, "importer")
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 2 || len(result.ImportedIDs) != 1 || len(result.Rejected) != 1 {
		t.Fatalf("report = %+v", result)
	}
	if summary.Total != 1 || summary.Imported != 1 || summary.Rejected != 1 {
		t.Fatalf("summary = %+v", summary)
	}
}
