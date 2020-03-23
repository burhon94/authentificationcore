package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/burhon94/alifMux/pkg/mux"
	"github.com/burhon94/authentificationcore/cmd/auth/app"
	"github.com/burhon94/authentificationcore/pkg/core/add"
	"github.com/burhon94/authentificationcore/pkg/core/token"
	"github.com/burhon94/authentificationcore/pkg/core/user"
	"github.com/burhon94/authentificationcore/pkg/jwt"
	"github.com/burhon94/bdi/pkg/di"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
)
//-host 0.0.0.0 -port 9999 -dsn postgres://user:pass@localhost:5401/auth -key alifkey
var (
	host = flag.String("host", "", "Server host")
	port = flag.String("port", "", "Server port")
	dsn  = flag.String("dsn", "", "Postgres DSN")
	secret = flag.String("key", "", "key")
)

type DSN string

func main() {
	flag.Parse()
	addr := net.JoinHostPort(*host, *port)
	keySecret := jwt.Secret(*secret)
	start(addr, *dsn, keySecret)
}

func start(addr string, dsn string, secret jwt.Secret) {
	container := di.NewContainer()

	err := container.Provide(
		app.NewServer,
		mux.NewExactMux,
		func() jwt.Secret { return secret },
		func() DSN { return DSN(dsn) },
		func(dsn DSN) *pgxpool.Pool {
			pool, err := pgxpool.Connect(context.Background(), string(dsn))
			if err != nil {
				panic(fmt.Errorf("can't create pool: %w", err))
			}
			return pool
		},
		token.NewService,
		user.NewService,
		add.NewService,
	)
	if err != nil {
		panic(fmt.Sprintf("can't set provide: %v", err))
	}

	container.Start()
	var appServer *app.Server
	container.Component(&appServer)

	log.Printf("authSvc listinig: ... %s", addr)
	panic(http.ListenAndServe(addr, appServer))
}

