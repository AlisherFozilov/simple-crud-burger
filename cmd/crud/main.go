package main

import (
	"context"
	"flag"
	"github.com/AlisherFozilov/crud/cmd/crud/app"
	"github.com/AlisherFozilov/crud/pkg/crud/services/burgers"
	"github.com/jackc/pgx/v4/pgxpool"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var (
	host = flag.String("host", "0.0.0.0", "Server host")
	port = flag.String("port", "9999", "Server port")
	dsn  = flag.String("dsn", "postgres://user:pass@localhost:5432/app", "Postgres DSN")
)

func main() {
	flag.Parse()
	portEnv, ok := os.LookupEnv("PORT")
	if ok {
		*port = portEnv
	}
	dsnEnv, ok := os.LookupEnv("DATABASE_URL")
	if ok {
		*dsn = dsnEnv
	}

	addr := net.JoinHostPort(*host, *port)
	start(addr, *dsn)
}

func start(addr string, dsn string) {
	router := app.NewExactMux()
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(err)
	}

	burgersSvc := burgers.NewBurgersSvc(pool)
	server := app.NewServer(
		router,
		pool,
		burgersSvc, // DI + Containers
		filepath.Join("web", "templates"),
		filepath.Join("web", "assets"),
	)
	server.InitRoutes()

	panic(http.ListenAndServe(addr, server))
}
