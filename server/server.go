package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo/database/handler"
	"todo/middleware"
)

type Server struct {
	*gin.Engine
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRoutes() *Server {
	routes := gin.Default()
	userAuth := routes.Group("/api")
	{
		userAuth.POST("/register", handler.RegisterUser)
		userAuth.POST("/login", handler.LoginUser)

		userTask := userAuth.Group("/todo")
		{
			userTask.Use(middleware.AuthMiddleware())
			userTask.POST("/task", handler.CreateTask)
			userTask.GET("/all-task", handler.GetAllTask)
			userTask.PUT("/:id", handler.UpdateUser)
			userTask.DELETE("/:id", handler.DeleteTask)
		}
	}
	return &Server{
		Engine: routes,
	}
}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Engine,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return svc.server.ListenAndServe()
}
