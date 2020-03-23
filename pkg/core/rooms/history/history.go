package history

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

type RoomsHistory struct {
	Id          int64  `json:"id"`
	RoomId      int64  `json:"room_id"`
	UserLogin   string  `json:"user_login"`
	NameMeeting string `json:"name_meeting"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
	Result      string `json:"result"`
}

func (s *Service) AddNewHistory(history RoomsHistory, pool *pgxpool.Pool) (err error) {
	var first, second int64
	err = pool.QueryRow(context.Background(), `SELECT id FROM rooms_history WHERE start_time <= $1 and end_time >= $1 and room_id = $2`, history.StartTime, history.RoomId).Scan(&first)
	if err == nil{
		return
	}
	err = pool.QueryRow(context.Background(), `SELECT id FROM rooms_history WHERE start_time <= $1 and end_time >= $1 and room_id = $2`, history.StartTime, history.RoomId).Scan(&second)
	if err == nil {
		return
	}
	if first == 0 && second == 0{
		_, err = pool.Exec(context.Background(), `INSERT INTO rooms_history(room_id, user_login, name_meeting, start_time, end_time)
VALUES ($1, $2, $3, $4, $5);`, history.RoomId, history.UserLogin,  history.NameMeeting, history.StartTime, history.EndTime)
		if err != nil {
			return
		}
		return nil}else{
		return errors.New("In this time meeting has have")
	}
}

func (s *Service) AllHistory (pool *pgxpool.Pool) (list []RoomsHistory ,err error)  {
	rows, err := pool.Query(context.Background(), `SELECT id, room_id, user_login, name_meeting, start_time, end_time, result FROM rooms_history;`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		history := RoomsHistory{}
		err := rows.Scan(&history.Id, &history.RoomId, &history.UserLogin, &history.NameMeeting, &history.StartTime, &history.EndTime, &history.Result)
		if err != nil {
			return nil, err
		}
		list = append(list, history)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return
}

func (s *Service) HistoryByRoomID(id int64, pool *pgxpool.Pool) (list []RoomsHistory, err error) {
	rows, err := pool.Query(context.Background(), `select id,room_id, user_login, name_meeting, start_time, end_time, result  from rooms_history where room_id=$1`, id)
	for rows.Next() {
		history := RoomsHistory{}
		err := rows.Scan(&history.Id, &history.RoomId, &history.UserLogin, &history.NameMeeting, &history.StartTime, &history.EndTime, &history.Result)
		if err != nil {
			return nil, err
		}
		list = append(list, history)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return
}

func (s *Service) AddResultById(history RoomsHistory, id int64, pool *pgxpool.Pool) (err error) {
	timestamp := time.Now().Unix()
	var first int64
	err = pool.QueryRow(context.Background(), `SELECT id FROM rooms_history WHERE end_time <= $1 and id = $2`, timestamp, id).Scan(&first)
	if err != nil{
		return
	}
	if first != 0{
		_, err = pool.Exec(context.Background(), `update rooms_history set result = $1 where id = $2`, history.Result, id)
		if err != nil {
			return
		}
		return nil}else{
		return errors.New("meeting have now end")
	}
}

func (s *Service) HistoryCurrentlyAndInThisRoom(id int64, pool *pgxpool.Pool) (history RoomsHistory, err error) {
	timestamp := time.Now().Unix()
	err = pool.QueryRow(context.Background(),
		`SELECT id, room_id, user_login, name_meeting, start_time, end_time, result FROM rooms_history WHERE start_time <= $1 and end_time >= $1 and room_id = $2`,
		timestamp, id).Scan(&history.Id, &history.RoomId, &history.UserLogin, &history.NameMeeting, &history.StartTime, &history.EndTime, &history.Result)
	if err != nil {
		return RoomsHistory{}, nil
	}
	return

}