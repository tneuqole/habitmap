package handlers

import (
	"net/http"

	"github.com/tneuqole/habitmap/internal/templates/pages"
)

func (h *BaseHandler) GetHome(w http.ResponseWriter, r *http.Request) error {
	return h.render(w, r, pages.Home(h.Session.Data(r.Context())))
}
