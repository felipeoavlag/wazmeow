package router

import (
	"net/http"
	"time"

	"wazmeow/internal/container"
	"wazmeow/internal/http/handlers/chat"
	"wazmeow/internal/http/handlers/contact"
	"wazmeow/internal/http/handlers/group"
	"wazmeow/internal/http/handlers/message"
	"wazmeow/internal/http/handlers/middleware"
	"wazmeow/internal/http/handlers/newsletter"
	"wazmeow/internal/http/handlers/session"
	"wazmeow/internal/http/handlers/webhook"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// healthCheck handler para verificação de saúde da API
// @Summary Verificação de saúde da API
// @Description Retorna o status de saúde da API WazMeow
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "API funcionando corretamente"
// @Router /health [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "wazmeow-api",
	})
}

// NewRouter cria um novo roteador HTTP usando handlers refatorados
func NewRouter(container *container.Container) http.Handler {
	r := chi.NewRouter()

	// Middlewares globais
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
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
	r.Get("/health", healthCheck)

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Criar handlers refatorados
	sessionHandler := session.NewHandler(
		container.GetCreateSessionUseCase(),
		container.GetConnectSessionUseCase(),
		container.GetListSessionsUseCase(),
		container.GetGetQRCodeUseCase(),
		container.GetDeleteSessionUseCase(),
		container.GetLogoutSessionUseCase(),
		container.GetPairPhoneUseCase(),
		container.GetGetSessionInfoUseCase(),
		container.GetSetProxyUseCase(),
	)

	webhookHandler := webhook.NewHandler(
		container.GetSetWebhookUseCase(),
		container.GetGetWebhookUseCase(),
	)

	messageHandler := message.NewHandler(
		container.GetSendTextMessageUseCase(),
		container.GetSendMediaMessageUseCase(),
		container.GetSendImageMessageUseCase(),
		container.GetSendAudioMessageUseCase(),
		container.GetSendDocumentMessageUseCase(),
		container.GetSendVideoMessageUseCase(),
		container.GetSendStickerMessageUseCase(),
		container.GetSendLocationMessageUseCase(),
		container.GetSendContactMessageUseCase(),
		container.GetSendButtonsMessageUseCase(),
		container.GetSendListMessageUseCase(),
		container.GetSendPollMessageUseCase(),
		container.GetSendEditMessageUseCase(),
		container.GetDeleteMessageUseCase(),
		container.GetReactMessageUseCase(),
	)

	contactHandler := contact.NewHandler(
		container.GetGetUserInfoUseCase(),
		container.GetCheckUserUseCase(),
		container.GetGetAvatarUseCase(),
		container.GetGetContactsUseCase(),
	)

	chatHandler := chat.NewHandler(
		container.GetSendPresenceUseCase(),
		container.GetChatPresenceUseCase(),
		container.GetMarkReadUseCase(),
		container.GetDownloadImageUseCase(),
		container.GetDownloadVideoUseCase(),
		container.GetDownloadAudioUseCase(),
		container.GetDownloadDocumentUseCase(),
	)

	groupHandler := group.NewHandler(
		container.GetCreateGroupUseCase(),
		container.GetSetGroupPhotoUseCase(),
		container.GetUpdateGroupParticipantsUseCase(),
		container.GetLeaveGroupUseCase(),
		container.GetJoinGroupUseCase(),
		container.GetGetGroupInfoUseCase(),
		container.GetListGroupsUseCase(),
		container.GetGetGroupInviteLinkUseCase(),
		container.GetRevokeGroupInviteLinkUseCase(),
		container.GetSetGroupNameUseCase(),
		container.GetSetGroupTopicUseCase(),
	)

	newsletterHandler := newsletter.NewHandler(
		container.GetListNewsletterUseCase(),
	)

	// Rotas de sessões (sem middleware de validação de sessionID)
	r.Route("/sessions", func(r chi.Router) {
		r.Use(middleware.CombinedValidation) // Middleware básico
		r.Get("/", sessionHandler.ListSessions)
		r.Post("/add", sessionHandler.CreateSession)

		r.Route("/{sessionID}", func(r chi.Router) {
			r.Use(middleware.SessionValidator) // Middleware de validação de sessionID
			r.Get("/", sessionHandler.GetSessionInfo)
			r.Delete("/", sessionHandler.DeleteSession)
			r.Post("/connect", sessionHandler.ConnectSession)
			r.Post("/logout", sessionHandler.LogoutSession)
			r.Get("/qr", sessionHandler.GetQRCode)
			r.Post("/pair", sessionHandler.PairPhone)
			r.Post("/proxy", sessionHandler.SetProxy)

			// Rotas de webhook dentro de sessions
			r.Route("/webhook", func(r chi.Router) {
				r.Get("/", webhookHandler.FindWebhook) // GET /sessions/{sessionID}/webhook
				r.Post("/set", webhookHandler.SetWebhook)
				r.Get("/find", webhookHandler.FindWebhook) // Manter compatibilidade
			})
		})
	})

	// Rotas de mensagens
	r.Route("/message", func(r chi.Router) {
		r.Use(middleware.CombinedValidation)
		r.Route("/{sessionID}", func(r chi.Router) {
			r.Use(middleware.SessionValidator)
			r.Route("/send", func(r chi.Router) {
				// Mensagens básicas
				r.Post("/text", messageHandler.SendTextMessage)
				r.Post("/media", messageHandler.SendMediaMessage)

				// Mensagens específicas
				r.Post("/image", messageHandler.SendImageMessage)
				r.Post("/audio", messageHandler.SendAudioMessage)
				r.Post("/document", messageHandler.SendDocumentMessage)
				r.Post("/video", messageHandler.SendVideoMessage)
				r.Post("/sticker", messageHandler.SendStickerMessage)
				r.Post("/location", messageHandler.SendLocationMessage)
				r.Post("/contact", messageHandler.SendContactMessage)
				r.Post("/buttons", messageHandler.SendButtonsMessage)
				r.Post("/list", messageHandler.SendListMessage)
				r.Post("/poll", messageHandler.SendPollMessage)
				r.Post("/edit", messageHandler.SendEditMessage)
			})

			// Operações de mensagem
			r.Post("/delete", messageHandler.DeleteMessage)
			r.Post("/react", messageHandler.ReactMessage)
		})
	})

	// Rotas de webhook globais
	r.Route("/webhook", func(r chi.Router) {
		r.Use(middleware.CombinedValidation)
		// Rotas globais de webhook
		r.Get("/events", webhookHandler.GetSupportedEvents)

		r.Route("/{sessionID}", func(r chi.Router) {
			r.Use(middleware.SessionValidator)
			r.Get("/", webhookHandler.FindWebhook)          // GET /webhook/{sessionID}
			r.Post("/", webhookHandler.ReceiveWebhookEvent) // POST /webhook/{sessionID} - Receber eventos do webhook
			r.Post("/set", webhookHandler.SetWebhook)
			r.Get("/find", webhookHandler.FindWebhook) // Manter compatibilidade
		})
	})

	// Rotas de contato
	r.Route("/contact", func(r chi.Router) {
		r.Use(middleware.CombinedValidation)
		r.Route("/{sessionID}", func(r chi.Router) {
			r.Use(middleware.SessionValidator)
			r.Post("/info", contactHandler.GetUserInfo)
			r.Post("/check", contactHandler.CheckUser)
			r.Post("/avatar", contactHandler.GetAvatar)
			r.Get("/list", contactHandler.GetContacts)
		})
	})

	// Rotas de chat
	r.Route("/chat", func(r chi.Router) {
		r.Use(middleware.CombinedValidation)
		r.Route("/{sessionID}", func(r chi.Router) {
			r.Use(middleware.SessionValidator)
			r.Post("/presence", chatHandler.SendPresence)
			r.Post("/chatpresence", chatHandler.ChatPresence)
			r.Post("/markread", chatHandler.MarkRead)
			r.Post("/download/image", chatHandler.DownloadImage)
			r.Post("/download/video", chatHandler.DownloadVideo)
			r.Post("/download/audio", chatHandler.DownloadAudio)
			r.Post("/download/document", chatHandler.DownloadDocument)
		})
	})

	// Rotas de grupo
	r.Route("/group", func(r chi.Router) {
		r.Use(middleware.CombinedValidation)
		r.Route("/{sessionID}", func(r chi.Router) {
			r.Use(middleware.SessionValidator)
			r.Post("/create", groupHandler.CreateGroup)
			r.Post("/photo", groupHandler.SetGroupPhoto)
			r.Post("/participants", groupHandler.UpdateGroupParticipants)
			r.Post("/leave", groupHandler.LeaveGroup)
			r.Post("/join", groupHandler.JoinGroup)
			r.Post("/info", groupHandler.GetGroupInfo)
			r.Get("/list", groupHandler.ListGroups)
			r.Post("/invitelink", groupHandler.GetGroupInviteLink)
			r.Delete("/invitelink", groupHandler.RevokeGroupInviteLink)
			r.Post("/name", groupHandler.SetGroupName)
			r.Post("/topic", groupHandler.SetGroupTopic)
		})
	})

	// Rotas de newsletter
	r.Route("/newsletter", func(r chi.Router) {
		r.Use(middleware.CombinedValidation)
		r.Route("/{sessionID}", func(r chi.Router) {
			r.Use(middleware.SessionValidator)
			r.Get("/list", newsletterHandler.ListNewsletter)
		})
	})

	return r
}
