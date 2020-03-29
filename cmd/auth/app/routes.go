package app

import (
	"github.com/burhon94/authentificationcore/pkg/core/token"
	"github.com/burhon94/authentificationcore/pkg/middleware/authenticated"
	"github.com/burhon94/authentificationcore/pkg/middleware/jwt"
	"github.com/burhon94/authentificationcore/pkg/middleware/logger"
	"reflect"
)

func (s *Server) InitRoutes() {
	s.router.GET(
		"/api/health",
		s.handleHealth(),
		logger.Logger("HEALTH"),
		)

	s.router.POST(
		"/api/tokens",
		s.handleCreateToken(),
		logger.Logger("TOKEN"),
	)

	s.router.POST(
		"/api/users/0",
		s.handleAddUser(),
		logger.Logger("REGISTRATION"),
	)

	s.router.GET(
		"/api/users/me",
		s.handleProfile(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.UserStruct)(nil)).Elem(), s.secret),
		logger.Logger("USERS"),
	)

	s.router.POST(
		"/api/users/{id}/edit",
		s.handleUpdateUser(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.UserStruct)(nil)).Elem(), s.secret),
		logger.Logger("USERS_EDIT"),
		)

	s.router.POST(
		"/api/users/{id}/pass",
		s.handleCheckPass(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.UserStruct)(nil)).Elem(), s.secret),
		logger.Logger("USERS_EDIT"),
	)

	s.router.POST(
		"/api/users/{id}/edit/pass",
		s.handleUpdatePass(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.UserStruct)(nil)).Elem(), s.secret),
		logger.Logger("USERS_CHECK_PASS"),
		)
}