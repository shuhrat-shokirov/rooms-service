package app

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shuhrat-shokirov/jwt/pkg/cmd"
	"github.com/shuhrat-shokirov/mux/pkg/mux"
	"github.com/shuhrat-shokirov/rest/pkg/rest"
	"log"
	"net/http"
	"rooms-service/pkg/core/rooms"
	"strconv"
)

type Server struct {
	router *mux.ExactMux
	pool   *pgxpool.Pool
	roomsSvc *rooms.Service
	secret       jwt.Secret
}

func NewServer(router *mux.ExactMux, pool *pgxpool.Pool, roomsSvc *rooms.Service, secret jwt.Secret) *Server {
	return &Server{router: router, pool: pool, roomsSvc: roomsSvc, secret: secret}
}




func (s Server) ServeHTTP(writer http.ResponseWriter,request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s Server) Start() {
	s.InitRoutes()
}

func (s Server) handleRoomsList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := s.roomsSvc.AllRooms(s.pool)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		err = rest.WriteJSONBody(writer, &list)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}

func (s Server) handleRoomByID() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromCTX, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer,http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idFromCTX)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		prod, err := s.roomsSvc.RoomByID(int64(id), s.pool)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		err = rest.WriteJSONBody(writer, &prod)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}

func (s Server) handleNewRooms() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		get := request.Header.Get("Content-Type")
		if get != "application/json" {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		room := rooms.Rooms{}
		err := rest.ReadJSONBody(request, &room)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Print(room)
		err = s.roomsSvc.AddNewRooms(room, s.pool)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		_, err = writer.Write([]byte("New Product Added!"))
		if err != nil {
			log.Print(err)
		}
	}
}

func (s Server) handleDeleteRooms() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromCTX, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer,http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idFromCTX)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = s.roomsSvc.RemoveByID(int64(id), s.pool)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		_, err = writer.Write([]byte("Product removed!"))
		if err != nil {
			log.Print(err)
		}
	}
}