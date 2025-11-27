package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name            string
		status          int
		data            interface{}
		wantStatus      int
		wantBody        string
		wantErr         bool
		wantContentType string
	}{
		{
			name:            "encode simple map",
			status:          http.StatusOK,
			data:            map[string]string{"message": "hello"},
			wantStatus:      http.StatusOK,
			wantBody:        `{"message":"hello"}`,
			wantErr:         false,
			wantContentType: "application/json",
		},
		{
			name:            "encode struct",
			status:          http.StatusCreated,
			data:            testStruct{Name: "test", Value: 42},
			wantStatus:      http.StatusCreated,
			wantBody:        `{"name":"test","value":42}`,
			wantErr:         false,
			wantContentType: "application/json",
		},
		{
			name:            "encode empty map",
			status:          http.StatusOK,
			data:            map[string]string{},
			wantStatus:      http.StatusOK,
			wantBody:        `{}`,
			wantErr:         false,
			wantContentType: "application/json",
		},
		{
			name:            "encode nil map",
			status:          http.StatusOK,
			data:            map[string]string(nil),
			wantStatus:      http.StatusOK,
			wantBody:        `null`,
			wantErr:         false,
			wantContentType: "application/json",
		},
		{
			name:            "encode array",
			status:          http.StatusOK,
			data:            []string{"a", "b", "c"},
			wantStatus:      http.StatusOK,
			wantBody:        `["a","b","c"]`,
			wantErr:         false,
			wantContentType: "application/json",
		},
		{
			name:            "encode with different status code",
			status:          http.StatusBadRequest,
			data:            map[string]string{"error": "bad request"},
			wantStatus:      http.StatusBadRequest,
			wantBody:        `{"error":"bad request"}`,
			wantErr:         false,
			wantContentType: "application/json",
		},
		{
			name:            "encode nested structure",
			status:          http.StatusOK,
			data:            map[string]interface{}{"user": map[string]string{"name": "john"}},
			wantStatus:      http.StatusOK,
			wantBody:        `{"user":{"name":"john"}}`,
			wantErr:         false,
			wantContentType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test", nil)

			err := encode(w, r, tt.status, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if w.Code != tt.wantStatus {
				t.Errorf("encode() status = %v, want %v", w.Code, tt.wantStatus)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != tt.wantContentType {
				t.Errorf("encode() Content-Type = %v, want %v", contentType, tt.wantContentType)
			}

			// Compare JSON by unmarshaling and remarshaling to normalize whitespace
			var gotJSON, wantJSON interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &gotJSON); err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.wantBody), &wantJSON); err != nil {
				t.Fatalf("failed to unmarshal expected body: %v", err)
			}

			gotBytes, _ := json.Marshal(gotJSON)
			wantBytes, _ := json.Marshal(wantJSON)

			if string(gotBytes) != string(wantBytes) {
				t.Errorf("encode() body = %v, want %v", string(gotBytes), string(wantBytes))
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name    string
		body    string
		want    interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name:    "decode simple map",
			body:    `{"message":"hello"}`,
			want:    map[string]interface{}{"message": "hello"},
			wantErr: false,
		},
		{
			name:    "decode struct",
			body:    `{"name":"test","value":42}`,
			want:    testStruct{Name: "test", Value: 42},
			wantErr: false,
		},
		{
			name:    "decode empty object",
			body:    `{}`,
			want:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "decode array",
			body:    `["a","b","c"]`,
			want:    []interface{}{"a", "b", "c"},
			wantErr: false,
		},
		{
			name:    "decode invalid JSON",
			body:    `{invalid json}`,
			want:    map[string]interface{}(nil),
			wantErr: true,
			errMsg:  "decode json",
		},
		{
			name:    "decode empty body",
			body:    ``,
			want:    map[string]interface{}(nil),
			wantErr: true,
			errMsg:  "decode json",
		},
		{
			name:    "decode malformed JSON",
			body:    `{"name":"test"`,
			want:    map[string]interface{}(nil),
			wantErr: true,
			errMsg:  "decode json",
		},
		{
			name:    "decode nested structure",
			body:    `{"user":{"name":"john","age":30}}`,
			want:    map[string]interface{}{"user": map[string]interface{}{"name": "john", "age": float64(30)}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.body))

			// Use type assertion to call decode with the correct type
			switch tt.want.(type) {
			case map[string]interface{}:
				got, err := decode[map[string]interface{}](r)
				if (err != nil) != tt.wantErr {
					t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.wantErr && err != nil {
					if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
						t.Errorf("decode() error = %v, want error containing %v", err, tt.errMsg)
					}
					return
				}
				// Compare as JSON strings
				gotJSON, _ := json.Marshal(got)
				wantJSON, _ := json.Marshal(tt.want)
				if string(gotJSON) != string(wantJSON) {
					t.Errorf("decode() = %v, want %v", string(gotJSON), string(wantJSON))
				}

			case testStruct:
				got, err := decode[testStruct](r)
				if (err != nil) != tt.wantErr {
					t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr && got != tt.want.(testStruct) {
					t.Errorf("decode() = %v, want %v", got, tt.want)
				}

			case []interface{}:
				got, err := decode[[]interface{}](r)
				if (err != nil) != tt.wantErr {
					t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					gotJSON, _ := json.Marshal(got)
					wantJSON, _ := json.Marshal(tt.want)
					if string(gotJSON) != string(wantJSON) {
						t.Errorf("decode() = %v, want %v", string(gotJSON), string(wantJSON))
					}
				}
			}
		})
	}
}

func TestEncode_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		wantErr bool
	}{
		{
			name:    "encode channel (unencodable type)",
			data:    make(chan int),
			wantErr: true,
		},
		{
			name:    "encode function (unencodable type)",
			data:    func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test", nil)

			err := encode(w, r, http.StatusOK, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), "encode json") {
					t.Errorf("encode() error = %v, want error containing 'encode json'", err)
				}
			}
		})
	}
}

func TestDecode_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		body    io.Reader
		wantErr bool
	}{
		{
			name:    "decode with nil body",
			body:    nil,
			wantErr: true,
		},
		{
			name:    "decode with closed reader",
			body:    io.NopCloser(bytes.NewReader([]byte(`{"test":"value"}`))),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/test", tt.body)

			_, err := decode[map[string]interface{}](r)

			if (err != nil) != tt.wantErr {
				t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncodeDecode_RoundTrip(t *testing.T) {
	type testData struct {
		Name   string   `json:"name"`
		Age    int      `json:"age"`
		Active bool     `json:"active"`
		Tags   []string `json:"tags"`
	}

	tests := []struct {
		name string
		data testData
	}{
		{
			name: "round trip with complete struct",
			data: testData{
				Name:   "John Doe",
				Age:    30,
				Active: true,
				Tags:   []string{"developer", "golang"},
			},
		},
		{
			name: "round trip with empty tags",
			data: testData{
				Name:   "Jane Doe",
				Age:    25,
				Active: false,
				Tags:   []string{},
			},
		},
		{
			name: "round trip with nil tags",
			data: testData{
				Name:   "Bob Smith",
				Age:    40,
				Active: true,
				Tags:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			if err := encode(w, r, http.StatusOK, tt.data); err != nil {
				t.Fatalf("encode() failed: %v", err)
			}

			// Decode
			r2 := httptest.NewRequest(http.MethodPost, "/test", w.Body)
			got, err := decode[testData](r2)
			if err != nil {
				t.Fatalf("decode() failed: %v", err)
			}

			// Compare
			if got.Name != tt.data.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.data.Name)
			}
			if got.Age != tt.data.Age {
				t.Errorf("Age = %v, want %v", got.Age, tt.data.Age)
			}
			if got.Active != tt.data.Active {
				t.Errorf("Active = %v, want %v", got.Active, tt.data.Active)
			}

			// Compare tags (handling nil vs empty slice)
			gotJSON, _ := json.Marshal(got.Tags)
			wantJSON, _ := json.Marshal(tt.data.Tags)
			if string(gotJSON) != string(wantJSON) {
				t.Errorf("Tags = %v, want %v", got.Tags, tt.data.Tags)
			}
		})
	}
}
