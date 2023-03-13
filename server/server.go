package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
	"todo/database/handler"
	"todo/middleware"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRoutes() *Server {
	routes := chi.NewRouter()
	routes.Route("/api", func(api chi.Router) {
		api.Post("/register", handler.RegisterUser)
		api.Post("/login", handler.LoginUser)
		api.Route("/todo", func(todo chi.Router) {
			todo.Use(middleware.AuthMiddleware)
			todo.Post("/task", handler.CreateTask)
			todo.Get("/all-task", handler.GetAllTask)
			todo.Put("/{id}", handler.UpdateUser)
			todo.Delete("/{id}", handler.DeleteTask)
		})
	})
	return &Server{
		Router: routes,
	}
}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return svc.server.ListenAndServe()
}
