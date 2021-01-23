package httpservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/namphung1998/auth-service-go/internal"
	"github.com/namphung1998/auth-service-go/internal/mock"
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
		data           request
		result         result
		userServiceDep func(*gomock.Controller) internal.UserService
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
			t.Parallel()
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

func TestHandleCreateUser(t *testing.T) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type result struct {
		status int
	}

	goodRequest := request{
		Email:    "nam@test.com",
		Password: "password",
	}

	tests := map[string]struct {
		data           request
		result         result
		userServiceDep func(*gomock.Controller) internal.UserService
	}{
		"empty request": {
			result: result{
				status: http.StatusBadRequest,
			},
		},
		"malformed email": {
			data: request{
				Email:    "bad",
				Password: "password",
			},
			result: result{
				status: http.StatusBadRequest,
			},
		},
		"service returns an error": {
			data: goodRequest,
			result: result{
				status: http.StatusInternalServerError,
			},
			userServiceDep: func(c *gomock.Controller) internal.UserService {
				m := mock.NewMockUserService(c)
				m.EXPECT().Create(gomock.Eq(goodRequest.Email), gomock.Eq(goodRequest.Password)).Return(errors.New("random")).Times(1)
				return m
			},
		},
		"all is well": {
			data: goodRequest,
			result: result{
				status: http.StatusCreated,
			},
			userServiceDep: func(c *gomock.Controller) internal.UserService {
				m := mock.NewMockUserService(c)
				m.EXPECT().Create(gomock.Eq(goodRequest.Email), gomock.Eq(goodRequest.Password)).Return(nil).Times(1)
				return m
			},
		},
	}

	for name, tt := range tests {
		name := name
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			payload, _ := json.Marshal(&tt.data)
			r := httptest.NewRequest("POST", "/", bytes.NewBuffer(payload))
			w := httptest.NewRecorder()

			handler := Handler{}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			if tt.userServiceDep != nil {
				handler.userService = tt.userServiceDep(ctrl)
			}

			handler.HandleCreateUser().ServeHTTP(w, r)

			assert.Equal(t, tt.result.status, w.Result().StatusCode)
		})
	}
}
