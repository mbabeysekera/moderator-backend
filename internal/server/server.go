package server

import (
	"fmt"
	"log"

	"coolbreez.lk/moderator/config"
	"coolbreez.lk/moderator/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Port      string
	AppServer *gin.Engine
}

func New(pool *pgxpool.Pool) (*Server, error) {
	engine := gin.Default()
	basePath, err := config.GetBasePath()
	if err != nil {
		return nil, err
	}

	routerGroup := engine.Group(basePath)
	routes.RegisterHealthCheckRoutes(routerGroup)
	routes.RegisterEventRoutes(routerGroup, pool)

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
