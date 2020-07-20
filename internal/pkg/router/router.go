package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/cors"

	"github.com/dink10/enlabs/internal/pkg/logger"
	"github.com/dink10/enlabs/internal/pkg/server"
	"github.com/dink10/enlabs/third_party/swagger"
)

// NewDefaultRouter returns router with CORS and request logging middlewares,
// health-check and swagger documentation end-points.
func NewDefaultRouter(enableLogging bool) Router {
	mux := chi.NewRouter()

	corsMiddleware := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Gismart-Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler

	mux.Use(corsMiddleware)
	mux.Use(render.SetContentType(render.ContentTypeJSON))
	mux.Use(middleware.RequestID)
	mux.Use(middleware.StripSlashes)
	mux.Use(middleware.RealIP)

	if enableLogging {
		mux.Use(middleware.RequestLogger(logger.NewRequestLogger()))
	}

	mux.Get("/", server.HealthCheck())
	mux.Get("/health", server.HealthCheck())
	mux.Get("/swagger/*", swagger.Documentation())

	return Router{mux: mux}
}

// Router wraps chi.Router.
type Router struct {
	mux chi.Router
}

// Routes is a type used while creating a subrouter.
type Routes map[string]http.Handler

// AddSubRouter adds subrouter to Router's mux.
//
// Usage example:
//     r := router.NewDefaultRouter()
//     r.AddSubRouter("/v1", router.Routes{
//         "/users": userProvider.Router(),
//         "/books": bookProvider.Router(),
//     })
//
func (r *Router) AddSubRouter(pattern string, routes Routes) {
	r.mux.Route(pattern, func(router chi.Router) {
		for pattern, handler := range routes {
			router.Mount(pattern, handler)
		}
	})
}

// Handler return Router's mux.
func (r *Router) Handler() http.Handler {
	return r.mux
}
