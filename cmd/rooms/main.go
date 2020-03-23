package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shuhrat-shokirov/jwt/pkg/cmd"
	"github.com/shuhrat-shokirov/new-mux/pkg/mux"
	"net"
	"net/http"
	"os"
	"rooms-service/cmd/rooms/app"
	"rooms-service/pkg/core/rooms"
	"rooms-service/pkg/core/rooms/history"
)

var (
	host = flag.String("host", "", "Server host")
	port = flag.String("port", "", "Server port")
	dsn  = flag.String("dsn", "", "Postgres DSN")
)
//-host 0.0.0.0 -port 9999 -dsn postgres://user:pass@localhost:5430/product
const (
	envHost = "HOST"
	envPort = "PORT"
	envDSN  = "DATABASE_URL"
)

type DSN string

func main() {
	flag.Parse()
	serverHost := checkENV(envHost, *host)
	serverPort := checkENV(envPort, *port)
	serverDsn := checkENV(envDSN, *dsn)
	addr := net.JoinHostPort(serverHost, serverPort)
	secret := jwt.Secret("secret")
	start(addr, serverDsn, secret)
}
func checkENV(env string, loc string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		return loc
	}
	return str
}
func start(addr string, dsn string,  secret jwt.Secret) {
	pool, err := pgxpool.Connect(context.Background(), string(dsn))
	if err != nil {
		panic(fmt.Errorf("can't create pool: %w", err))
	}

	roomsSvc := rooms.NewService()
	exactMux := mux.NewExactMux()
	historySvc := history.NewService()
	server := app.NewServer(exactMux, pool, roomsSvc, secret, historySvc)
	server.Start()
	panic(http.ListenAndServe(addr, server))

}