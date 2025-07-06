package handlers

import (
	"net/http"

	"github.com/tneuqole/habitmap/internal/ctxutil"
	"github.com/tneuqole/habitmap/internal/logutil"
	"github.com/tneuqole/habitmap/internal/templates/pages"
)

func (h *BaseHandler) renderErrorPage(w http.ResponseWriter, r *http.Request, statusCode int) {
	logger := ctxutil.GetLogger(r.Context())
	sessionData := h.Session.Data(r.Context())

	w.WriteHeader(statusCode)

	var err error
	switch statusCode {
	case http.StatusNotFound:
		err = h.render(w, r, pages.Error404(sessionData))
	default:
		err = h.render(w, r, pages.Error500(sessionData))
	}

	if err != nil {
		logger.Error("could not write error response", logutil.ErrorSlog(err))
	}
}

func (h *BaseHandler) Error404(w http.ResponseWriter, r *http.Request) {
	h.renderErrorPage(w, r, http.StatusNotFound)
}
