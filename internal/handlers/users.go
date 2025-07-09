package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/tneuqole/habitmap/internal/apperror"
	"github.com/tneuqole/habitmap/internal/forms"
	"github.com/tneuqole/habitmap/internal/model"
	"github.com/tneuqole/habitmap/internal/templates/formcomponents"
	"github.com/tneuqole/habitmap/internal/templates/pages"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	*BaseHandler
}

func NewUserHandler(bh *BaseHandler) *UserHandler {
	return &UserHandler{
		BaseHandler: bh,
	}
}

func (h *UserHandler) GetSignupForm(w http.ResponseWriter, r *http.Request) error {
	return h.render(w, r, formcomponents.Signup(h.Session.Data(r.Context()), forms.SignupForm{}))
}

func (h *UserHandler) PostSignup(w http.ResponseWriter, r *http.Request) error {
	sessionData := h.Session.Data(r.Context())

	var form forms.SignupForm
	if err := h.bindFormData(r, &form); err != nil {
		return err
	}

	err := validate.Struct(&form)
	if err != nil {
		form.FieldErrors = h.parseValidationErrors(r.Context(), err)
		return h.render(w, r, formcomponents.Signup(sessionData, form))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	params := model.CreateUserParams{
		Name:           form.Name,
		Email:          form.Email,
		HashedPassword: string(hashedPassword),
	}
	userID, err := h.Queries.CreateUser(r.Context(), params)
	if err != nil {
		err = h.handleDBError(err)
		if errors.Is(err, apperror.ErrDuplicateEmail) {
			form.AddFieldError("Email", apperror.ErrDuplicateEmail.Message)
			return h.render(w, r, formcomponents.Signup(sessionData, form))
		}
	}

	h.Session.SetUserID(r.Context(), userID)

	http.Redirect(w, r, "/habits", http.StatusSeeOther)
	return nil
}

func (h *UserHandler) GetLoginForm(w http.ResponseWriter, r *http.Request) error {
	return h.render(w, r, formcomponents.Login(h.Session.Data(r.Context()), forms.LoginForm{}))
}

func (h *UserHandler) PostLogin(w http.ResponseWriter, r *http.Request) error {
	sessionData := h.Session.Data(r.Context())

	var form forms.LoginForm
	if err := h.bindFormData(r, &form); err != nil {
		return err
	}

	err := validate.Struct(&form)
	if err != nil {
		form.FieldErrors = h.parseValidationErrors(r.Context(), err)
		return h.render(w, r, formcomponents.Login(sessionData, form))
	}

	user, err := h.Queries.GetUser(r.Context(), form.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			form.AddGenericError(apperror.ErrInvalidCredentials.Message)
			return h.render(w, r, formcomponents.Login(sessionData, form))
		}

		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(form.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			form.AddGenericError(apperror.ErrInvalidCredentials.Message)
			return h.render(w, r, formcomponents.Login(sessionData, form))
		}
		return err
	}

	err = h.Session.RenewToken(r.Context())
	if err != nil {
		return err
	}

	h.Session.SetUserID(r.Context(), user.ID)

	http.Redirect(w, r, "/habits", http.StatusSeeOther)
	return nil
}

func (h *UserHandler) PostLogout(w http.ResponseWriter, r *http.Request) error {
	h.Session.RemoveUserID(r.Context())
	h.Session.SetFlash(r.Context(), "You've been logged out.")

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (h *UserHandler) GetAccount(w http.ResponseWriter, r *http.Request) error {
	sessionData := h.Session.Data(r.Context())
	if !sessionData.IsAuthenticated {
		return h.render(w, r, pages.Error404(sessionData))
	}

	user, err := h.Queries.GetUserByID(r.Context(), *sessionData.UserID)
	if err != nil {
		return h.handleDBError(err)
	}

	return h.render(w, r, pages.Account(sessionData, user))
}
