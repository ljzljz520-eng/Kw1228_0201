package httpapi

import (
	"net/http"
	"strings"
)

func (s *Server) adminRoutes() {
	s.mux.HandleFunc("/review-queue", s.reviewQueue)
	s.mux.HandleFunc("/cores/", s.coreActions)
}

func (s *Server) reviewQueue(response http.ResponseWriter, request *http.Request) {
	if !methodAllowed(response, request, http.MethodGet) {
		return
	}
	queue, err := s.runtime.ReviewQueue()
	if err != nil {
		writeError(response, err)
		return
	}
	writeJSON(response, http.StatusOK, queue)
}

func (s *Server) coreActions(response http.ResponseWriter, request *http.Request) {
	path := strings.TrimPrefix(request.URL.Path, "/cores/")
	if !strings.Contains(path, "/actions/") {
		s.coreByID(response, request)
		return
	}
	parts := strings.SplitN(path, "/actions/", 2)
	id, action := parts[0], parts[1]
	if request.Method != http.MethodPost {
		methodAllowed(response, request, http.MethodPost)
		return
	}
	switch action {
	case "submit":
		record, err := s.runtime.Catalog.Submit(id, actor(request))
		if err != nil {
			writeError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, record)
	case "publish":
		record, err := s.runtime.Publish(id, actor(request))
		if err != nil {
			writeError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, record)
	case "archive":
		record, err := s.runtime.Catalog.Archive(id, actor(request))
		if err != nil {
			writeError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, record)
	case "reset-review":
		if err := s.runtime.ResetReviewState(id); err != nil {
			writeError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, map[string]string{"reset": id})
	default:
		http.NotFound(response, request)
	}
}
