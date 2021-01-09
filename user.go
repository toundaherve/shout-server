package main

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
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

func (f Form) validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Username, validation.Required),
		validation.Field(&f.FirstName, validation.Required),
		validation.Field(&f.LastName, validation.Required),
		validation.Field(&f.Email, validation.Required, is.Email),
		validation.Field(&f.Password, validation.Required),
		validation.Field(&f.Location, validation.Required),
	)
}

// UserRegistrationHandler handles the creation of a new user
func UserRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)

	var form Form

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

	err = form.validate()
	if err != nil {
		resp := APIResponse{
			Code:    http.StatusBadRequest,
			Type:    "unknown",
			Message: err,
		}

		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(resp)
		return
	}

}
