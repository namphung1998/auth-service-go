package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleGreet(t *testing.T) {
	type request struct {
		Name string `json:"name"`
	}

	type response struct {
		Greeting string `json:"greeting"`
	}

	type result struct {
		status   int
		response response
	}

	tests := map[string]struct {
		data   request
		result result
	}{
		"empty name": {
			result: result{
				status: http.StatusBadRequest,
			},
		},
		"all is well": {
			data: request{
				Name: "Nam",
			},
			result: result{
				status: http.StatusOK,
				response: response{
					Greeting: "Hello, Nam.",
				},
			},
		},
	}

	for name, tt := range tests {
		name := name
		tt := tt

		t.Run(name, func(t *testing.T) {
			payload, _ := json.Marshal(&tt.data)
			r := httptest.NewRequest("POST", "/", bytes.NewBuffer(payload))
			w := httptest.NewRecorder()

			handler := Handler{}
			handler.HandleGreet().ServeHTTP(w, r)

			assert.Equal(t, tt.result.status, w.Result().StatusCode)

			var actualResult response
			json.Unmarshal(w.Body.Bytes(), &actualResult)
			assert.Equal(t, tt.result.response, actualResult)
		})
	}
}
