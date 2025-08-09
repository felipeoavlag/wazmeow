package http

import (
	"net/http"
	"time"

	"wazmeow/internal/app"
	"wazmeow/internal/http/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// NewRouter cria um novo roteador HTTP usando use cases
func NewRouter(container *app.DependencyContainer) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "wazmeow-api",
		})
	})

	// Criar handler de sessões com use cases
	sessionHandler := handlers.NewSessionHandler(
		container.CreateSessionUC,
		container.ConnectSessionUC,
		container.ListSessionsUC,
		container.GetQRCodeUC,
		container.DeleteSessionUC,
		container.LogoutSessionUC,
		container.PairPhoneUC,
		container.GetSessionInfoUC,
		container.SetProxyUC,
	)

	// Rotas de sessões
	r.Route("/sessions", func(r chi.Router) {
		r.Get("/", sessionHandler.ListSessions)
		r.Post("/add", sessionHandler.CreateSession)

		r.Route("/{sessionId}", func(r chi.Router) {
			r.Get("/", sessionHandler.GetSessionInfo)
			r.Delete("/", sessionHandler.DeleteSession)
			r.Post("/connect", sessionHandler.ConnectSession)
			r.Post("/logout", sessionHandler.LogoutSession)
			r.Get("/qr", sessionHandler.GetQRCode)
			r.Post("/pair", sessionHandler.PairPhone)
			r.Post("/proxy", sessionHandler.SetProxy)
		})
	})

	return r
}
