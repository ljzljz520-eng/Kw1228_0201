package importer

import "testing"

func TestImportPreview(t *testing.T) {
	rows, issues, err := Preview("code,name,description,owner,classification\nA,Alpha,Stable unit,Lab,alpha\nB,,Missing,Lab,beta")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || len(issues) != 1 {
		t.Fatalf("rows=%v issues=%v", rows, issues)
	}
}
