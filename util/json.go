package util

import (
	"encoding/json"
	"net/http"
)

func ParseJSON(data []byte, source interface{}) error {
	err := json.Unmarshal(data, source)
	return err
}

func ReadFromRequestBody(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	PanicIfError(err)
}

func WriteToResponseBody(writer http.ResponseWriter, status int, response any) {
	if status == 0 {
		status = http.StatusOK
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	PanicIfError(err)
}
