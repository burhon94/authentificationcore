package app

import (
	"github.com/burhon94/alifMux/pkg/mux"
	"github.com/burhon94/authentificationcore/pkg/core/add"
	"github.com/burhon94/authentificationcore/pkg/core/token"
	"github.com/burhon94/authentificationcore/pkg/core/user"
	"github.com/burhon94/authentificationcore/pkg/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

type Server struct {
	router   *mux.ExactMux
	pool     *pgxpool.Pool
	secret   jwt.Secret
	tokenSvc *token.Service
	userSvc  *user.Service
	addUser  *add.Service
}

func NewServer(router *mux.ExactMux, pool *pgxpool.Pool, secret jwt.Secret, tokenSvc *token.Service, userSvc *user.Service, addUser *add.Service) *Server {
	return &Server{router: router, pool: pool, secret: secret, tokenSvc: tokenSvc, userSvc: userSvc, addUser: addUser}
}

func (s *Server) Start() {
	s.InitRoutes()
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}
