package util

import (
	"encoding/json"
	"net/http"
)

func ParseJSON(data []byte, source interface{}) error {
	err := json.Unmarshal(data, source)
	return err
}

func ToJSON(data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	return bytes, err
}

func ReadFromRequestBody(request *http.Request, result interface{}) error {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	if err != nil {
		return NewErrorException(ErrBadRequest, "Invalid request body: "+err.Error())
	}

	val := NewValidator()
	err = val.Struct(result)
	if err != nil {
		return NewErrorException(ErrValidation, "Validation error: "+err.Error())
	}

	return nil
}

func WriteToResponseBody(writer http.ResponseWriter, status int, response any) error {
	if status == 0 {
		status = http.StatusOK
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)

	if err != nil {
		return ErrInternalServer
	}

	return nil
}
