package httpservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/namphung1998/auth-service-go/internal"
)

type Handler struct {
	userService internal.UserService
}

func NewHandler(userService internal.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

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

// HandleCreateUser returns a HandlerFunc that handles registering a new user
func (h *Handler) HandleCreateUser() http.HandlerFunc {
	type request struct {
		Email    string `json:"email" valid:"email,required"`
		Password string `json:"password" valid:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decode(r, &req); err != nil {
			fmt.Println(err)
			writeResponse(w, nil, http.StatusBadRequest)
			return
		}

		if err := h.userService.Create(req.Email, req.Password); err != nil {
			fmt.Println(err)
			switch err.(type) {
			case *internal.InvalidRequestError:
				writeResponse(w, nil, http.StatusBadRequest)
			case *internal.EmailInUseError:
				writeResponse(w, nil, http.StatusConflict)
			default:
				writeResponse(w, nil, http.StatusInternalServerError)
			}
			return
		}

		writeResponse(w, nil, http.StatusCreated)
	}
}

// HandleLogin returns a HandlerFunc that handles logging in a user
func (h *Handler) HandleLogin() http.HandlerFunc {
	type request struct {
		Email    string `json:"email" valid:"email,required"`
		Password string `json:"password" valid:"required"`
	}

	type response struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decode(r, &req); err != nil {
			fmt.Println(err)
			writeResponse(w, nil, http.StatusBadRequest)
			return
		}

		token, err := h.userService.Login(req.Email, req.Password)
		if err != nil {
			fmt.Println(err)
			switch err.(type) {
			case *internal.UserNotFoundError:
				writeResponse(w, nil, http.StatusNotFound)
			case *internal.IncorrectPasswordError:
				writeResponse(w, nil, http.StatusUnauthorized)
			default:
				writeResponse(w, nil, http.StatusInternalServerError)
			}
			return
		}

		writeResponse(w, response{Token: token}, http.StatusOK)
	}
}
