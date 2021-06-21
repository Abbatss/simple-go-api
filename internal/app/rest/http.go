package rest

import (
	"encoding/json"
	"net/http"
)

func PlainText(response http.ResponseWriter, status int, body string) {
	response.Header().Set("Content-Type", "text/plain")
	response.WriteHeader(status)
	_, _ = response.Write([]byte(body))
}

func JSON(response http.ResponseWriter, status int, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		PlainText(response, http.StatusInternalServerError, "Wow, things are really broken")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_, _ = response.Write(body)
}

func BindJSON(request *http.Request, data interface{}) error {
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}
