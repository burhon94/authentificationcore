package user

import (
	"context"
	"errors"
	"github.com/burhon94/authentificationcore/pkg/core/token"
	"github.com/burhon94/authentificationcore/pkg/middleware/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type ResponseDTO struct {
	Id int64 `json:"id"`
	Login string `json:"login"`
	NameSurname string `json:"name_surname"`
	Avatar string `json:"avatar"`
	Role []string `json:"role"`
}

type ResponseChangeDTO struct {
	Id int64 `json:"id"`
	Password string `json:"password"`
	NameSurname string `json:"name_surname"`
	Avatar string `json:"avatar"`
}

type ResponseCheckPass struct {
	Id   int64  `json:"id"`
	Pass string `json:"pass"`
}

func (s *Service) Profile(ctx context.Context) (response ResponseDTO, err error) {
	auth, ok := jwt.FromContext(ctx).(*token.UserStruct)
	if !ok {
		return ResponseDTO{}, errors.New("error ")
	}

	var userData ResponseDTO
	err = s.pool.QueryRow(ctx, `SELECT login, namesurname, avatar, roles FROM users WHERE id = $1;`, auth.Id).Scan(&userData.Login, &userData.NameSurname, &userData.Avatar, &userData.Role)
	if err != nil {
		return ResponseDTO{}, errors.New("error ")
	}
	userData.Id = auth.Id

	return userData, nil
}

func (s *Service) UpdateUser(ctx context.Context, id int64, userName string) (err error) {
	_, err = s.pool.Exec(ctx, `UPDATE users SET namesurname = $2 WHERE id = $1;`, id, userName)
	if err != nil {
		return errors.New("error ")
	}

	return nil
}

func (s *Service) CheckPass(ctx context.Context, id int64) (response ResponseCheckPass, err error) {
	var userData ResponseCheckPass
	err = s.pool.QueryRow(ctx, `SELECT password from users WHERE id = $1;`, id).Scan(&userData.Pass)
	if err != nil {
		return ResponseCheckPass{}, errors.New("error ")
	}
	userData.Id = id

	return userData, nil
}

func (s *Service) UpdatePass(ctx context.Context, id int64, pass string) (err error) {
	_, err = s.pool.Exec(ctx, `UPDATE users SET password = $2 WHERE id = $1;`, id, pass)
	if err != nil {
		return errors.New("error ")
	}

	return nil
}