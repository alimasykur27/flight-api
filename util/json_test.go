package util_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"flight-api/util"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---- Unit Test Parse JSON ----
// ---- helpers ----
// assertCanonicalJSONMatch menyamakan hasil dengan cara:
// 1) marshal 'got' (source) ke JSON
// 2) unmarshal keduanya ke struktur generic (map/array/bool/number string) pakai UseNumber
// 3) samakan dengan DeepEqual setelah normalisasi number -> float64
func assertCanonicalJSONMatch(t *testing.T, got interface{}, want map[string]interface{}) {
	t.Helper()

	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("failed to marshal got: %v", err)
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("failed to marshal want: %v", err)
	}

	gotVal, err := decodeCanonical(gotJSON)
	if err != nil {
		t.Fatalf("failed to decode got: %v", err)
	}
	wantVal, err := decodeCanonical(wantJSON)
	if err != nil {
		t.Fatalf("failed to decode want: %v", err)
	}

	if !reflect.DeepEqual(gotVal, wantVal) {
		t.Fatalf("mismatch.\n got: %s\nwant: %s", bytes.TrimSpace(gotJSON), bytes.TrimSpace(wantJSON))
	}
}

// decodeCanonical mengubah JSON jadi struktur interface{} dengan number sebagai float64
func decodeCanonical(b []byte) (interface{}, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()

	var v interface{}
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return normalizeJSONNumber(v), nil
}

// normalizeJSONNumber: json.Number -> float64; rekursif untuk map/array
func normalizeJSONNumber(v interface{}) interface{} {
	switch x := v.(type) {
	case json.Number:
		// coba parse float (cukup untuk test)
		f, _ := x.Float64()
		return f
	case map[string]interface{}:
		m := make(map[string]interface{}, len(x))
		for k, val := range x {
			m[k] = normalizeJSONNumber(val)
		}
		return m
	case []interface{}:
		arr := make([]interface{}, len(x))
		for i, val := range x {
			arr[i] = normalizeJSONNumber(val)
		}
		return arr
	default:
		return v
	}
}

func TestParseJSON(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		source      interface{}            // pointer ke struct target decode
		expected    map[string]interface{} // ekspektasi payload JSON (pakai kunci JSON)
		wantErr     bool
		errContains string
	}

	tests := []testCase{
		{
			name: "valid input 1 (int field)",
			data: []byte(`{"name":"John","age":30}`),
			source: &struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{},
			expected: map[string]interface{}{"name": "John", "age": float64(30)},
		},
		{
			name: "valid input 2 (int price)",
			data: []byte(`{"product":"laptop","price":200000}`),
			source: &struct {
				Product string `json:"product"`
				Price   int    `json:"price"`
			}{},
			expected: map[string]interface{}{"product": "laptop", "price": float64(200000)},
		},
		{
			name: "valid input 3 (string float + bool)",
			data: []byte(`{"product":"motor","berat":"25.50","is_available":true}`),
			source: &struct {
				Product     string `json:"product"`
				Berat       string `json:"berat"` // data asli "25.50" (string), decode ke string
				IsAvailable bool   `json:"is_available"`
			}{},
			expected: map[string]interface{}{"product": "motor", "berat": "25.50", "is_available": true},
		},
		{
			name: "type compatible: json number -> float32",
			data: []byte(`{"name":"John","age":30}`),
			source: &struct {
				Name string  `json:"name"`
				Age  float32 `json:"age"`
			}{},
			// setelah marshal balik ke JSON, angka akan jadi float64 saat dinormalisasi
			expected: map[string]interface{}{"name": "John", "age": float64(30)},
		},
		{
			name: "invalid input: json number -> string (should error)",
			data: []byte(`{"product":"laptop","price":200000}`),
			source: &struct {
				Product string `json:"product"`
				Price   string `json:"price"`
			}{},
			wantErr:     true,
			errContains: "cannot unmarshal number into Go struct field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Pastikan source pointer ke struct
			rt := reflect.TypeOf(tt.source)
			if rt == nil || rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Struct {
				t.Fatalf("source must be a non-nil pointer to struct, got %T", tt.source)
			}

			err := util.ParseJSON(tt.data, tt.source)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("error should contain %q, got %v", tt.errContains, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Bandingkan hasil decode (source) vs expected map dalam bentuk canonical JSON
			assertCanonicalJSONMatch(t, tt.source, tt.expected)
		})
	}
}

// ---- Unit Test ReadFromRequestBody ----
func TestReadFromRequestBody(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	t.Run("decodes valid JSON", func(t *testing.T) {
		body := []byte(`{"name":"John","age":42}`)
		req := httptest.NewRequest(http.MethodPost, "/x", bytes.NewReader(body))

		var got payload
		util.ReadFromRequestBody(req, &got)

		if got.Name != "John" || got.Age != 42 {
			t.Fatalf("unexpected decode result: %+v", got)
		}
	})

	t.Run("panics on invalid JSON", func(t *testing.T) {
		body := []byte(`{"name":"John","age":}`)
		req := httptest.NewRequest(http.MethodPost, "/x", bytes.NewReader(body))

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic, got none")
			}
		}()

		var got payload
		util.ReadFromRequestBody(req, &got) // should panic via PanicIfError
	})
}

// ---- Unit Test WriteToResponseBody ----
type errorWriter struct{}

func (e *errorWriter) Header() http.Header {
	return http.Header{}
}
func (e *errorWriter) WriteHeader(statusCode int) {}
func (e *errorWriter) Write(b []byte) (int, error) {
	return 0, errors.New("write error")
}

func TestWriteToResponseBody(t *testing.T) {
	t.Run("valid - status 0 defaults to 200", func(t *testing.T) {
		rr := httptest.NewRecorder()
		resp := map[string]string{"msg": "ok"}

		util.WriteToResponseBody(rr, 0, resp)

		if rr.Code != http.StatusOK {
			assert.IsType(t, http.StatusOK, rr.Code)
			assert.Equal(t, http.StatusOK, rr.Code)
			t.Fatalf("expected status 200, got %d", rr.Code)
		}
		if rr.Header().Get("Content-Type") != "application/json" {
			assert.IsType(t, "application/json", rr.Header().Get("Content-Type"))
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
			t.Fatalf("expected application/json, got %s", rr.Header().Get("Content-Type"))
		}
		var got map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
			assert.IsType(t, map[string]string{}, got)
			assert.Equal(t, map[string]string{"msg": "ok"}, got)
			t.Fatalf("response not valid JSON: %v", err)
		}
		if got["msg"] != "ok" {
			assert.IsType(t, "ok", got["msg"])
			assert.Equal(t, "ok", got["msg"])
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("valid - custom status 201", func(t *testing.T) {
		rr := httptest.NewRecorder()
		resp := map[string]string{"msg": "created"}

		util.WriteToResponseBody(rr, http.StatusCreated, resp)

		if rr.Code != http.StatusCreated {
			assert.IsType(t, http.StatusCreated, rr.Code)
			assert.Equal(t, http.StatusCreated, rr.Code)
			t.Fatalf("expected status 201, got %d", rr.Code)
		}
	})

	t.Run("invalid - unencodable type (channel)", func(t *testing.T) {
		rr := httptest.NewRecorder()
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic, got none")
			}
		}()

		resp := map[string]interface{}{"bad": make(chan int)}
		util.WriteToResponseBody(rr, 200, resp) // should panic
	})

	t.Run("invalid - writer error", func(t *testing.T) {
		w := &errorWriter{}
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic, got none")
			}
		}()

		resp := map[string]string{"msg": "fail"}
		util.WriteToResponseBody(w, 200, resp) // should panic
	})
}
