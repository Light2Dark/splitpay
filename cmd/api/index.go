package main

import (
	"net/http"

	"github.com/Light2Dark/splitpay/internal/templates"
)

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.Index().Render(r.Context(), w)
}
