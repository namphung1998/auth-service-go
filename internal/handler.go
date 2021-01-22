package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
)

type Handler struct{}

func writeResponse(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	if data == nil {
		return
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println(err)
	}
}

func decode(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	_, err := govalidator.ValidateStruct(v)
	if err != nil {
		return err
	}
	return nil
}

// HandleGreet returns a HandlerFunc for greeting a user
func (h *Handler) HandleGreet() http.HandlerFunc {
	type request struct {
		Name string `json:"name" valid:"required"`
	}

	type response struct {
		Greeting string `json:"greeting"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decode(r, &req); err != nil {
			fmt.Println(err)
			writeResponse(w, nil, http.StatusBadRequest)
			return
		}

		writeResponse(w, response{Greeting: fmt.Sprintf("Hello, %v.", req.Name)}, http.StatusOK)
	}
}
