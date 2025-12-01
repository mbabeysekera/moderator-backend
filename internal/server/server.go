package server

import (
	"fmt"
	"log"

	"coolbreez.lk/moderator/config"
	"coolbreez.lk/moderator/internal/routes"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Port      string
	AppServer *gin.Engine
}

func New() (*Server, error) {
	engine := gin.Default()
	err := routes.RegisterRoutes(engine)
	if err != nil {
		log.Fatalf("Server initialization failed: %v", err)
		return nil, err
	}
	appPort, err := config.GetServerPort()
	if err != nil {
		log.Fatal("Application server port has not been set")
		return nil, err
	}

	return &Server{
		Port:      appPort,
		AppServer: engine,
	}, nil
}

func (s *Server) Start() error {
	err := s.AppServer.Run(fmt.Sprintf(":%s", s.Port))
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
