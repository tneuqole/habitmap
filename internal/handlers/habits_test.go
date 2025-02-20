package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/tneuqole/habitmap/internal/model"
	"github.com/tneuqole/habitmap/testdata"
)

func TestGetHabits(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/habits", strings.NewReader(""))
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	v, _ := NewValidate()

	testDB := testdata.NewTestDB()
	testDB.Setup()
	defer testDB.Teardown()
	q := model.New(testDB.DB)
	h := NewHabitsHandler(q, v)

	assert.NoError(t, h.GetHabits(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	// assert.Equal(t, "", rec.Body.String())
}
