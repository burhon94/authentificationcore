package user

import (
	"context"
	"errors"
	"github.com/burhon94/authentificationcore/pkg/core/token"
	"github.com/burhon94/authentificationcore/pkg/middleware/jwt"
)

type Service struct { }

func NewService() *Service {
	return &Service{}
}

type ResponseDTO struct {
	Id int64 `json:"id"`
	NameSurname string `json:"name_surname"`
	Avatar string `json:"avatar"`
	Role []string `json:"role"`
}

func (s *Service) Profile(ctx context.Context) (response ResponseDTO, err error) {
	auth, ok := jwt.FromContext(ctx).(*token.UserStruct)
	if !ok {
		return ResponseDTO{}, errors.New("error ")
	}

	return ResponseDTO{
		Id: auth.Id,
		NameSurname: auth.NameSurname,
		Avatar: "https://i.pravatar.cc/50",
		Role: auth.Roles,
	}, nil
}
