package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"wazmeow/internal/application/handlers"
)

// SetupRoutes configures all routes for the API
func SetupRoutes(router chi.Router, sessionHandler *handlers.SessionHandler) {
	// Health check endpoint
	router.Get("/health", healthCheckHandler)

	// Root endpoint
	router.Get("/", rootHandler)

	// Session management routes (direct paths as specified)
	setupSessionRoutes(router, sessionHandler)
}

// setupSessionRoutes configures session management routes
func setupSessionRoutes(router chi.Router, sessionHandler *handlers.SessionHandler) {
	router.Route("/sessions", func(r chi.Router) {
		// Session collection routes
		r.Post("/add", sessionHandler.CreateSession)
		r.Get("/list", sessionHandler.ListSessions)

		// Session-specific routes
		r.Route("/{sessionID}", func(r chi.Router) {
			r.Get("/info", sessionHandler.GetSessionInfo)
			r.Delete("/", sessionHandler.DeleteSession)
			r.Post("/connect", sessionHandler.ConnectSession)
			r.Post("/logout", sessionHandler.LogoutSession)
			r.Get("/qr", sessionHandler.GetQRCode)
			r.Post("/pairphone", sessionHandler.PairPhone)
			r.Post("/proxy/set", sessionHandler.SetProxy)
		})
	})
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"wazmeow"}`))
}

// rootHandler handles root endpoint requests
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"WazMeow API - WhatsApp Session Management","version":"1.0.0"}`))
}
