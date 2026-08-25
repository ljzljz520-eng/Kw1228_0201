package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"example.com/energycore/internal/model"
	"example.com/energycore/internal/store"
)

func (s *Server) health(response http.ResponseWriter, request *http.Request) {
	if !methodAllowed(response, request, http.MethodGet) {
		return
	}
	writeJSON(response, http.StatusOK, map[string]string{"status": "ok", "service": "energy-core-registry"})
}

func (s *Server) cores(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		filter := model.SearchFilter{Query: request.URL.Query().Get("q"), Owner: request.URL.Query().Get("owner"), Classification: request.URL.Query().Get("classification"), IncludeArchived: request.URL.Query().Get("include_archived") == "true"}
		if status := request.URL.Query().Get("status"); status != "" {
			filter.Status = model.RecordStatus(status)
		}
		items, err := s.runtime.Catalog.Search(filter)
		if err != nil {
			writeError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, items)
	case http.MethodPost:
		var input model.Record
		if err := readJSON(request, &input); err != nil {
			writeError(response, err)
			return
		}
		record, err := s.runtime.Catalog.Register(input, actor(request))
		if err != nil {
			writeError(response, err)
			return
		}
		writeJSON(response, http.StatusCreated, record)
	default:
		methodAllowed(response, request, http.MethodGet, http.MethodPost)
	}
}

func (s *Server) coreByID(response http.ResponseWriter, request *http.Request) {
	id := strings.TrimPrefix(request.URL.Path, "/cores/")
	id = strings.TrimSuffix(id, "/review")
	if id == "" {
		http.NotFound(response, request)
		return
	}
	switch request.Method {
	case http.MethodGet:
		record, err := s.runtime.Catalog.Get(id)
		if err != nil {
			writeError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, record)
	case http.MethodPut:
		var patch model.Record
		if err := readJSON(request, &patch); err != nil {
			writeError(response, err)
			return
		}
		record, err := s.runtime.Catalog.Update(id, patch, patch.Revision, actor(request))
		if err != nil {
			writeError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, record)
	case http.MethodDelete:
		if err := s.runtime.Catalog.Delete(id, actor(request)); err != nil {
			writeError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, map[string]string{"deleted": id})
	case http.MethodPost:
		s.handleReview(response, request, id)
	default:
		methodAllowed(response, request, http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPost)
	}
}

func (s *Server) handleReview(response http.ResponseWriter, request *http.Request, id string) {
	if !strings.HasSuffix(request.URL.Path, "/review") {
		http.NotFound(response, request)
		return
	}
	var input struct {
		Decision model.ReviewDecision `json:"decision"`
		Note     string               `json:"note"`
	}
	if err := readJSON(request, &input); err != nil {
		writeError(response, err)
		return
	}
	result, err := s.runtime.Review.Review(id, input.Decision, input.Note, actor(request))
	if err != nil {
		writeError(response, err)
		return
	}
	writeJSON(response, http.StatusOK, result)
}

func (s *Server) imports(response http.ResponseWriter, request *http.Request) {
	if !methodAllowed(response, request, http.MethodPost) {
		return
	}
	var input struct {
		CSV string `json:"csv"`
	}
	if err := readJSON(request, &input); err != nil {
		writeError(response, err)
		return
	}
	result, err := s.runtime.Import.ImportCSV(input.CSV, actor(request))
	if err != nil {
		writeError(response, err)
		return
	}
	writeJSON(response, http.StatusOK, result)
}

func (s *Server) summary(response http.ResponseWriter, request *http.Request) {
	if !methodAllowed(response, request, http.MethodGet) {
		return
	}
	result, err := s.runtime.Report.Summary()
	if err != nil {
		writeError(response, err)
		return
	}
	writeJSON(response, http.StatusOK, result)
}

func actor(request *http.Request) string {
	if value := request.Header.Get("X-Actor"); value != "" {
		return value
	}
	return "operator"
}

func writeJSON(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(value)
}

func writeError(response http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, store.ErrNotFound) {
		status = http.StatusNotFound
	}
	writeJSON(response, status, map[string]string{"error": err.Error()})
}
