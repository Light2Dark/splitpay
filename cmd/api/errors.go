package main

import (
	"fmt"
	"net/http"

	"github.com/Light2Dark/splitpay/internal/templates"
)

func (app application) logError(w http.ResponseWriter, r *http.Request, errorMessage string, err error) {
	app.logger.Error(errorMessage, "error", err)
	templates.Error(fmt.Sprintf("%s", errorMessage,)).Render(r.Context(), w)
}
