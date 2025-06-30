package server

import (
	"fmt"
	"net/http"

	"github.com/savin000/audit-log/internal/server/handlers"
	"github.com/savin000/audit-log/internal/server/routes"
)

func New(port uint32, h *handlers.Handler) *http.Server {
	mux := routes.RegisterRoutes(h)
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}
