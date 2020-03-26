package token

import (
	"context"
	"errors"
	"github.com/burhon94/authentificationcore/pkg/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service struct {
	secret jwt.Secret
	pool *pgxpool.Pool
}

func NewService(secret jwt.Secret, pool *pgxpool.Pool) *Service {
	return &Service{secret: secret, pool: pool}
}


type UserStruct struct {
	Id          int64    `json:"id"`
	Exp         int64    `json:"exp"`
	NameSurname string   `json:"name_surname"`
	Roles       []string `json:"roles"`
}

type RequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseDTO struct {
	Token string `json:"token"`
}

var ErrInvalidLogin = errors.New("invalid login")
var ErrInvalidPassword = errors.New("invalid password")

func (s *Service) Generate(context context.Context, request *RequestDTO) (response ResponseDTO, err error) {
	var id int64
	var userPass, userName string
	var roles []string
	err = s.pool.QueryRow(context, `SELECT password, namesurname, roles, id FROM users WHERE login = $1;
`, request.Username).Scan(&userPass, &userName, &roles, &id)
	if err != nil {
		err = ErrInvalidLogin
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPass), []byte(request.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		err = ErrInvalidPassword
		return
	}

	response.Token, err = jwt.Encode(UserStruct{
		Id:          id,
		Exp:         time.Now().Add(time.Hour).Unix(),
		NameSurname: userName,
		Roles:       roles,
	}, s.secret)
	if err != nil {
		return ResponseDTO{}, err
	}

	return
}
