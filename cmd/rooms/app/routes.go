package app

import (
	"rooms-service/pkg/mux/middleware/logger"
)

func (s Server) InitRoutes() {

	s.router.GET(
		"/api/rooms",
		s.handleRoomsList(),
		logger.Logger("get list"),
	)

	s.router.GET(
		"/api/rooms/{id}",
		s.handleRoomByID(),
		logger.Logger("get rooms by id"),
	)

	s.router.POST(
		"/api/rooms/0",
		s.handleNewRooms(),
		logger.Logger("post new room"),
	)

	s.router.DELETE(
		"/api/rooms/{id}",
		s.handleDeleteRooms(),
		logger.Logger("delete room"),
	)
	s.router.POST(
		"/api/rooms/lock/{id}",
		s.handleLockRooms(),
		logger.Logger("locked room"),
		)
	s.router.POST(
		"/api/rooms/unlock/{id}",
		s.handleUnLockRooms(),
		logger.Logger("locked room"),
	)
	s.router.GET(
		"/api/rooms/unlock",
		s.handleRoomsListUnlocked(),
		logger.Logger("locked room"),
	)
	s.router.GET(
		"/api/rooms/locked",
		s.handleRoomsListLocked(),
		logger.Logger("locked room"),
	)
	s.router.POST(
		"/api/rooms/history/0",
		s.handleNewRoomsHistory(),
		logger.Logger("add new history"),
	)
	s.router.GET(
		"/api/rooms/history",
		s.handleHistoryList(),
		logger.Logger("get list"),
	)
	s.router.GET(
		"/api/rooms/history/{id}",
		s.handleHistoryByRoomID(),
		logger.Logger("get history by room_id"),
	)
	s.router.POST(
		"/api/rooms/history/add/result/{id}",
		s.handleAddResultById(),
		logger.Logger("add result in history by id"),
	)
	s.router.GET(
		"/api/history/room/{id}",
		s.handleHistoryCurrentlyAndInThisRoom(),
		logger.Logger("get history by room_id"),
	)
}