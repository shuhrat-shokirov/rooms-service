package app

import (
	"context"
	"errors"
	"github.com/shuhrat-shokirov/mux/pkg/mux/middleware/authenticated"
	"github.com/shuhrat-shokirov/mux/pkg/mux/middleware/jwt"
	"github.com/shuhrat-shokirov/mux/pkg/mux/middleware/logger"
	"reflect"
	"rooms-service/pkg/core/token"
)

func (s Server) InitRoutes() {

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		panic(errors.New("can't create database"))
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `
Create table if not exists rooms
(
    id       Bigserial primary key,
    status   bool,
    timestart text NOT NULL ,
    timestop text NOT NULL ,
    filename  text NOT NULL ,
    removed   bool DEFAULT FALSE
);
`)
	if err != nil {
		panic(errors.New("can't create database"))
	}
	s.router.GET(
		"/api/rooms",
		s.handleRoomsList(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("get list"),
	)

	s.router.GET(
		"/api/rooms/{id}",
		s.handleRoomByID(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("get product by id"),
	)

	s.router.POST(
		"/api/rooms/new",
		s.handleNewRooms(),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("post new product"),
	)

	s.router.DELETE(
		"/api/rooms/{id}",
		s.handleDeleteRooms(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("delete product"),
	)


}