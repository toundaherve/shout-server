package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

// Form contains the
type Form struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Location  string `json:"location"`
}

func (f *Form) validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Username, validation.Required),
		validation.Field(&f.FirstName, validation.Required),
		validation.Field(&f.LastName, validation.Required),
		validation.Field(&f.Email, validation.Required, is.Email),
		validation.Field(&f.Password, validation.Required),
		validation.Field(&f.Location, validation.Required),
	)
}

func (f *Form) format() {
	f.Username = strings.ToLower(strings.TrimSpace(f.Username))
	f.FirstName = strings.ToLower(strings.TrimSpace(f.FirstName))
	f.LastName = strings.ToLower(strings.TrimSpace(f.LastName))
	f.Email = strings.ToLower(strings.TrimSpace(f.Email))
	f.Password = strings.ToLower(strings.TrimSpace(f.Password))
	f.Location = strings.ToLower(strings.TrimSpace(f.Location))

}

// NewUserRegistrationHandler creates a new Registration Handler
func NewUserRegistrationHandler(conn *pgx.Conn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)

		var form *Form

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&form)
		if err != nil {
			resp := APIResponse{
				Code:    http.StatusBadRequest,
				Type:    "unknown",
				Message: "Please check the data.",
			}

			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(resp)
			return
		}

		form.format()

		err = form.validate()
		if err != nil {
			resp := APIResponse{
				Code:    http.StatusBadRequest,
				Type:    "unknown",
				Message: err.Error(),
			}

			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(resp)
			return
		}

		var takenEmail string
		err = conn.QueryRow(context.Background(), "SELECT email FROM users WHERE email=$1", form.Email).Scan(&takenEmail)
		if err == nil {
			resp := APIResponse{
				Code:    http.StatusBadRequest,
				Type:    "unknown",
				Message: "Email already taken.",
			}

			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(resp)
			return
		}

		id := uuid.New().String()

		_, err = conn.Exec(context.Background(), "INSERT INTO users(id, username, first_name, last_name, email, password, location) VALUES($1, $2, $3, $4, $5, $6, $7)",
			id, form.Username, form.FirstName, form.LastName, form.Email, form.Password, form.Location)
		if err != nil {
			resp := APIResponse{
				Code:    http.StatusInternalServerError,
				Type:    "unknown",
				Message: "Sorry, we cannot perform this operation now.",
			}

			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(resp)
			return
		}

		resp := APIResponse{
			Code:    http.StatusCreated,
			Type:    "unknown",
			Message: "Registration successful",
		}

		w.WriteHeader(http.StatusCreated)
		enc.Encode(resp)
		return
	}

}
