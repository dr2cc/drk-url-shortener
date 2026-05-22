package handler

// type Handler struct {
// 	service *service.Service
// }

// func New(service *service.Service, cfg config.Config, log *slog.Logger) *chi.Mux {
// 	// 3️⃣handler
// 	r := chi.NewRouter()
// 	// Мы передаем настроенный logger внутрь middleware slog-chi
// 	r.Use(slogchi.New(log))
// 	r.Post("/", shortenText(service, cfg, log))
// 	r.Get("/{id}", redirect(service, log))
// 	return r
// }
