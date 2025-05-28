package handlers

import (
	"errors"
	"net/http"

	"github.com/tneuqole/habitmap/internal/apperror"
	"github.com/tneuqole/habitmap/internal/forms"
	"github.com/tneuqole/habitmap/internal/model"
	"github.com/tneuqole/habitmap/internal/templates/formcomponents"
	"golang.org/x/crypto/bcrypt"
)

const Cost = 12

type UserHandler struct {
	*BaseHandler
}

func NewUserHandler(bh *BaseHandler) *UserHandler {
	return &UserHandler{
		BaseHandler: bh,
	}
}

func (h *UserHandler) GetSignupForm(w http.ResponseWriter, r *http.Request) error {
	return h.render(w, r, formcomponents.Signup(forms.SignupForm{}))
}

func (h *UserHandler) PostSignup(w http.ResponseWriter, r *http.Request) error {
	var form forms.SignupForm
	if err := h.bindFormData(r, &form); err != nil {
		return err
	}

	err := validate.Struct(&form)
	if err != nil {
		form.Errors = h.parseValidationErrors(err)
		return h.render(w, r, formcomponents.Signup(form))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), Cost)
	if err != nil {
		return err
	}

	params := model.CreateUserParams{
		Name:           form.Name,
		Email:          form.Email,
		HashedPassword: string(hashedPassword),
	}
	_, err = h.Queries.CreateUser(r.Context(), params)
	if err != nil {
		err = h.handleDBError(err)
		if errors.Is(err, apperror.ErrDuplicateEmail) {
			form.AddError("Email", apperror.ErrDuplicateEmail.Message)
			return h.render(w, r, formcomponents.Signup(form))
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
