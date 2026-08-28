package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"example.com/energycore/internal/flow017"
	"example.com/energycore/internal/model"
)

func TestHTTPRegistrationAndSearch(t *testing.T) {
	runtime, err := flow017.OpenRuntime(filepath.Join(t.TempDir(), "http.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer runtime.Close()
	server := New(runtime)
	body := `{"code":"H-1","name":"HTTP Core","description":"Deterministic endpoint fixture","owner":"Web Lab","classification":"gamma"}`
	request := httptest.NewRequest(http.MethodPost, "/cores", strings.NewReader(body))
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	var record model.Record
	if err := json.NewDecoder(response.Body).Decode(&record); err != nil {
		t.Fatal(err)
	}
	search := httptest.NewRequest(http.MethodGet, "/cores?q=HTTP", nil)
	searchResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(searchResponse, search)
	if searchResponse.Code != http.StatusOK || !strings.Contains(searchResponse.Body.String(), record.ID) {
		t.Fatalf("search status=%d body=%s", searchResponse.Code, searchResponse.Body.String())
	}
}
