package add

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type NewUser struct {
	Username string `json:"username"`
	UserLogin string `json:"user_login"`
	Password string `json:"password"`
}

var badRequest = errors.New("bad request")

func (s *Service) AddNewUser(context context.Context, request NewUser) (err error) {
	if request.UserLogin == "" {
		return badRequest
	}

	if request.Username == "" {
		return badRequest
	}

	if request.Password == "" {
		return badRequest
	}

	passCrypt, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.New("errror while crypt pass")
		return err
	}

	avatar := "http://localhost:20000/web/media/18786405-a685-4726-951f-db336f300ef0.png"
	role := []string {"user"}

	_, err = s.pool.Exec(context, `INSERT INTO users (login, password, namesurname, avatar, roles) VALUES ($1, $2, $3, $4, $5);`, request.UserLogin, passCrypt, request.Username, avatar, role)
	if err != nil {
		return err
	}
	return
}
