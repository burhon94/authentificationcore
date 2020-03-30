package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/burhon94/authentificationcore/pkg/core/add"
	"github.com/burhon94/authentificationcore/pkg/core/token"
	jsonReader "github.com/burhon94/json/cmd/reader"
	jsonWriter "github.com/burhon94/json/cmd/writer"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type ErrorDTO struct {
	Errors []string `json:"errors"`
}

type userId struct {
	Id int64 `json:"id"`
}

type userUpdateName struct {
	Id          int64  `json:"id"`
	NameSurname string `json:"name_surname"`
}

type userUpdatePass struct {
	Id   int64  `json:"id"`
	Pass string `json:"pass"`
}

type userUpdateAvatar struct {
	Id   int64  `json:"id"`
	AvatarUrl string `json:"avatar_url"`
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
		err := jsonReader.ReadJSONHTTP(request, &userId{})
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		response, err := s.userSvc.Profile(ctx)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		err = jsonWriter.WriteJSONHTTP(writer, &response)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
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

func (s *Server) handleUpdateUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var dataRequest userUpdateName
		err := jsonReader.ReadJSONHTTP(request, &dataRequest)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		id := dataRequest.Id

		err = s.userSvc.UpdateUser(ctx, id, dataRequest.NameSurname)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		err = jsonWriter.WriteJSONHTTP(writer, struct {
			Status string `json:"status"`
		}{
			"ok",
		})
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}
	}
}

func (s *Server) handleCheckPass() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var dataRequest userId
		err := jsonReader.ReadJSONHTTP(request, &dataRequest)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		id := dataRequest.Id
		userData, err := s.userSvc.CheckPass(ctx, id)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		err = jsonWriter.WriteJSONHTTP(writer, &userData)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}
	}
}

func (s *Server) handleUpdatePass() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var dataRequest userUpdatePass
		err := jsonReader.ReadJSONHTTP(request, &dataRequest)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		id := dataRequest.Id
		pass, err := bcrypt.GenerateFromPassword([]byte(dataRequest.Pass), bcrypt.DefaultCost)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		err = s.userSvc.UpdatePass(ctx, id, string(pass))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		err = jsonWriter.WriteJSONHTTP(writer, struct {
			Status string `json:"status"`
		}{
			"ok",
		})
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}
	}
}

func (s *Server) handleUpdateAvatar() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var dataRequest userUpdateAvatar
		err := jsonReader.ReadJSONHTTP(request, &dataRequest)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		id := dataRequest.Id
		avatarUrl := dataRequest.AvatarUrl

		err = s.userSvc.UpdateAvatar(ctx, id, avatarUrl)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}

		err = jsonWriter.WriteJSONHTTP(writer, struct {
			Status string `json:"status"`
		}{
			"ok",
		})
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = jsonWriter.WriteJSONHTTP(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			return
		}
	}
}
