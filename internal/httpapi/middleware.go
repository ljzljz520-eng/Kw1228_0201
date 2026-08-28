package httpapi

import "net/http"

func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("X-Request-Sequence", "1")
		next.ServeHTTP(response, request)
	})
}

func (s *Server) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		next.ServeHTTP(response, request)
	})
}

func methodAllowed(response http.ResponseWriter, request *http.Request, methods ...string) bool {
	for _, method := range methods {
		if request.Method == method {
			return true
		}
	}
	response.Header().Set("Allow", methods[0])
	http.Error(response, "method not allowed", http.StatusMethodNotAllowed)
	return false
}
