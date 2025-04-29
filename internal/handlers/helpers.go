package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/tneuqole/habitmap/internal/apperror"
	"github.com/tneuqole/habitmap/internal/model"
	"github.com/tneuqole/habitmap/internal/util"
)

const (
	daysInWeek   = 7
	monthsInYear = 12
)

type BaseHandler struct {
	Logger  *slog.Logger
	Queries *model.Queries
}

func (h *BaseHandler) render(w http.ResponseWriter, r *http.Request, component templ.Component) error {
	return component.Render(r.Context(), w)
}

func (h *BaseHandler) getIDFromURL(r *http.Request) (int64, error) {
	param := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return -1, apperror.New(http.StatusBadRequest, "id must be an integer")
	}

	return id, nil
}

var formDecoder = form.NewDecoder()

// dest should be address of a struct
func (h *BaseHandler) bindFormData(r *http.Request, dest any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := formDecoder.Decode(dest, r.Form); err != nil {
		return apperror.New(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (h *BaseHandler) handleDBError(err error) error {
	h.Logger.Error("DATABASE_ERROR", util.ErrorSlog(err))
	if errors.Is(err, sql.ErrNoRows) {
		return apperror.New(http.StatusNotFound, "Resource does not exist")
	}

	return apperror.New(http.StatusInternalServerError, "Error reading from database")
}

// checks if the date matches "YYYY-MM" format
func validateYearMonth(fl validator.FieldLevel) bool {
	d := fl.Field().String()
	if _, err := time.Parse("2006-01", d); err == nil {
		return true
	}

	return false
}

func newValidate() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("notblank", validators.NotBlank) //nolint:errcheck,gosec
	validate.RegisterValidation("yearmonth", validateYearMonth)  //nolint:errcheck,gosec
	return validate
}

var validate = newValidate()

func (h *BaseHandler) parseValidationErrors(err error) map[string]string {
	errMsgs := make(map[string]string)

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, fieldErr := range validationErrors {
			var msg string
			switch fieldErr.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", fieldErr.Field())
			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters long", fieldErr.Field(), fieldErr.Param())
			case "max":
				msg = fmt.Sprintf("%s must be at most %s characters long", fieldErr.Field(), fieldErr.Param())
			case "alpha":
			case "alphanum":
			case "alphanumunicode":
			case "alphaunicode":
			case "ascii":
				msg = fmt.Sprintf("%s contains invalid characters", fieldErr.Field())
			case "notblank":
				msg = fmt.Sprintf("%s cannot be blank", fieldErr.Field())
			case "oneof":
				msg = fmt.Sprintf("%s must be one of [%s]", fieldErr.Field(), fieldErr.Param())
			case "yearmonth":
				msg = fmt.Sprintf("%s must be in the format YYYY or YYYY-MM", fieldErr.Field())
			default:
				h.Logger.Debug("default case", slog.String("type", fieldErr.Tag()))
				msg = fmt.Sprintf("%s is invalid", fieldErr.Field())
			}
			errMsgs[fieldErr.Field()] = msg
		}
	}

	return errMsgs
}

// generateMonth creates a 2D slice representing a month's calendar,
// with each week containing habit entries for 7 days. The week is populated based on
// the entries parameter, if no entry exists for a date then a blank Entry{} is created.
// The function pads weeks with empty days and months with empty weeks as needed.
//
// Parameters:
//
//	monthStr (string): The target month in "YYYY-MM" format.
//	entries ([]model.Entry): Habit entries to populate the calendar.
//
// Returns:
//
//	[][]model.Entry: A 2D slice with weekly habit entries.
func (h *BaseHandler) generateMonth(habitID int64, monthStr string, entries []model.Entry) [][]model.Entry {
	var month [][]model.Entry
	week := make([]model.Entry, daysInWeek)

	date, err := time.Parse("2006-01", monthStr)
	if err != nil {
		h.Logger.Error("Error parsing date", util.ErrorSlog(err))
		return month
	}

	daysInMonth := date.AddDate(0, 1, -1).Day()

	entryIdx := 0
	dayOfWeek := int(date.Weekday())
	for day := date.Day(); day <= daysInMonth; {
		for ; dayOfWeek < daysInWeek && day <= daysInMonth; dayOfWeek++ {
			if entryIdx < len(entries) && len(entries) > 0 && entries[entryIdx].EntryDate == date.Format("2006-01-02") {
				week[dayOfWeek] = entries[entryIdx]
				entryIdx++
			} else {
				entry := model.Entry{
					HabitID:   habitID,
					EntryDate: date.Format("2006-01-02"),
				}
				week[dayOfWeek] = entry
			}
			date = date.AddDate(0, 0, 1)
			day++
		}
		month = append(month, week)
		week = make([]model.Entry, daysInWeek)
		dayOfWeek = 0
	}

	for len(month) < 6 {
		week = make([]model.Entry, daysInWeek)
		month = append(month, week)
	}

	return month
}

func (h *BaseHandler) generateYearMonths(date string) ([]string, error) {
	t, err := time.Parse("2006-01", date)
	if err != nil {
		h.Logger.Error("invalid date format", util.ErrorSlog(err))
		return nil, apperror.New(http.StatusBadRequest, "invalid date format")
	}

	year := t.Year()
	months := make([]string, 0, monthsInYear)

	for month := 1; month <= 12; month++ {
		// Format each month as "YYYY-MM"
		m := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		months = append(months, m.Format("2006-01"))
	}

	return months, nil
}
