package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/burhon94/alifMux/pkg/mux"
	"github.com/burhon94/authentificationcore/pkg/core/add"
	"github.com/burhon94/authentificationcore/pkg/core/token"
	"github.com/burhon94/authentificationcore/pkg/core/user"
	"github.com/burhon94/authentificationcore/pkg/jwt"
	jsonReader "github.com/burhon94/json/cmd/reader"
	jsonWriter "github.com/burhon94/json/cmd/writer"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"time"
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

type ErrorDTO struct {
	Errors []string `json:"errors"`
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleCreateToken() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var body token.RequestDTO
		err := jsonReader.ReadJSONHTTP(request, &body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.json_invalid"},
			})
			log.Print(err)
			return
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		response, err := s.tokenSvc.Generate(ctx, &body)
		if err != nil {
			switch {
			case errors.Is(err, token.ErrInvalidLogin):
				writer.WriteHeader(http.StatusBadRequest)
				err := jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
					[]string{"err.login_mismatch"},
				})
				log.Print(err)
			case errors.Is(err, token.ErrInvalidPassword):
				writer.WriteHeader(http.StatusBadRequest)
				err := jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
					[]string{"err.password_mismatch"},
				})
				log.Print(err)
			default:
				writer.WriteHeader(http.StatusBadRequest)
				err := jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
					[]string{"err.unknown"},
				})
				log.Print(err)
			}
			return
		}
		err = jsonWriter.WriteJSONHTTP(writer, &response)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) handleAddUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var body add.NewUser
		err := jsonReader.ReadJSONHTTP(request, &body)
		if err != nil {
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{[]string{"err.json_invalid"},
			})
			return
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		err = s.addUser.AddNewUser(ctx, body)
		if err != nil {
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{[]string{"error while add user"},})
		}
	}
}

func (s *Server) handleProfile() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		response, err := s.userSvc.Profile(request.Context())
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			log.Print(err)
			return
		}
		err = jsonWriter.WriteJSONHTTP(writer, &response)
		if err != nil {
			log.Print(err)
		}

	}
}

func (s *Server) handleHealth() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, err := fmt.Fprintf(writer, "Health ok")
		if err != nil {
			log.Printf("err: %v", err)
		}
	}
}
