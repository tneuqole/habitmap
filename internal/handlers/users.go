package handlers

import (
	"database/sql"
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
		form.FieldErrors = h.parseValidationErrors(err)
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
			form.AddFieldError("Email", apperror.ErrDuplicateEmail.Message)
			return h.render(w, r, formcomponents.Signup(form))
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (h *UserHandler) GetLoginForm(w http.ResponseWriter, r *http.Request) error {
	return h.render(w, r, formcomponents.Login(forms.LoginForm{}))
}

func (h *UserHandler) PostLogin(w http.ResponseWriter, r *http.Request) error {
	var form forms.LoginForm
	if err := h.bindFormData(r, &form); err != nil {
		return err
	}

	err := validate.Struct(&form)
	if err != nil {
		form.FieldErrors = h.parseValidationErrors(err)
		return h.render(w, r, formcomponents.Login(form))
	}

	user, err := h.Queries.GetUser(r.Context(), form.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			form.AddGenericError(apperror.ErrInvalidCredentials.Message)
			return h.render(w, r, formcomponents.Login(form))
		}

		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(form.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			form.AddGenericError(apperror.ErrInvalidCredentials.Message)
			return h.render(w, r, formcomponents.Login(form))
		}
		return err
	}

	err = h.Sessions.RenewToken(r.Context())
	if err != nil {
		return err
	}

	h.Sessions.Put(r.Context(), "authenticatedUserID", user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
