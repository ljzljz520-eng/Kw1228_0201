package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

func readJSON(request *http.Request, target any) error {
	if request.Body == nil {
		return errors.New("request body is required")
	}
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	return nil
}
