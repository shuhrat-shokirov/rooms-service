package app

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shuhrat-shokirov/jwt/pkg/cmd"
	"github.com/shuhrat-shokirov/new-mux/pkg/mux"
	"github.com/shuhrat-shokirov/rest/pkg/rest"
	"log"
	"net/http"
	"rooms-service/pkg/core/rooms"
	"rooms-service/pkg/core/rooms/history"
	"strconv"
)

type Server struct {
	router *mux.ExactMux
	pool   *pgxpool.Pool
	roomsSvc *rooms.Service
	secret       jwt.Secret
	historySvc *history.Service
}

func NewServer(router *mux.ExactMux, pool *pgxpool.Pool, roomsSvc *rooms.Service, secret jwt.Secret, historySvc *history.Service) *Server {
	return &Server{router: router, pool: pool, roomsSvc: roomsSvc, secret: secret, historySvc: historySvc}
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
		_, err = writer.Write([]byte("New Rooms Added!"))
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
		_, err = writer.Write([]byte("Rooms removed!"))
		if err != nil {
			log.Print(err)
		}
	}
}

func (s Server) handleLockRooms() http.HandlerFunc {
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

		err = s.roomsSvc.LockRoomByID(int64(id), s.pool)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		_, err = writer.Write([]byte("Rooms locked!"))
		if err != nil {
			log.Print(err)
		}
	}
}

func (s Server) handleUnLockRooms() http.HandlerFunc {
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

		err = s.roomsSvc.UnLockRoomByID(int64(id), s.pool)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		_, err = writer.Write([]byte("Rooms unlocked!"))
		if err != nil {
			log.Print(err)
		}
	}
}

func (s Server) handleRoomsListUnlocked() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := s.roomsSvc.AllRoomsUnlocked(s.pool)
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

func (s Server) handleRoomsListLocked() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := s.roomsSvc.AllRoomsLocked(s.pool)
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

func (s Server) handleNewRoomsHistory() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		get := request.Header.Get("Content-Type")
		if get != "application/json" {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		history := history.RoomsHistory{}
		err := rest.ReadJSONBody(request, &history)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Print(history)
		err = s.historySvc.AddNewHistory(history, s.pool)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		_, err = writer.Write([]byte("New History Added!"))
		if err != nil {
			log.Print(err)
		}
	}
}

func (s Server) handleHistoryList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := s.historySvc.AllHistory(s.pool)
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

func (s Server) handleHistoryByRoomID() http.HandlerFunc {
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

		prod, err := s.historySvc.HistoryByRoomID(int64(id), s.pool)
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

func (s Server) handleAddResultById() http.HandlerFunc {
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
		get := request.Header.Get("Content-Type")
		if get != "application/json" {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		history := history.RoomsHistory{}
		err = rest.ReadJSONBody(request, &history)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Print(history)
		err = s.historySvc.AddResultById(history, int64(id), s.pool)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		_, err = writer.Write([]byte("New Result by id Added!"))
		if err != nil {
			log.Print(err)
		}
	}
}

func (s Server) handleHistoryCurrentlyAndInThisRoom() http.HandlerFunc {
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

		prod, err := s.historySvc.HistoryCurrentlyAndInThisRoom(int64(id), s.pool)
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