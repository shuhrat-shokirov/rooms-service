package rooms

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

type Rooms struct {
	Id     int64 `json:"id"`
	Status bool `json:"status"`
	TimeStart string `json:"time_start"`
	TimeStop string `json:"time_stop"`
	FileName string `json:"file_name"`
	Removed string `json:"removed"`
}

func (s *Service) AllRooms (pool *pgxpool.Pool) (list []Rooms ,err error)  {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), `SELECT id, status, timestart, timestop, filename FROM rooms where removed=false;`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := Rooms{}
		err := rows.Scan(&item.Id, &item.Status, &item.TimeStart, &item.TimeStop, &item.FileName)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return
}

func (s *Service) AddNewRooms(room Rooms, pool *pgxpool.Pool) (err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `INSERT INTO rooms(status, timestart, timestop, filename)
VALUES ($1, $2, $3, $4);`,room.Status, room.TimeStart, room.TimeStop, room.FileName)
	if err != nil {
		return
	}
	return nil
}

func (s *Service) RemoveByID(id int64, pool *pgxpool.Pool) (err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return errors.New("can't connect to database!")
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `update rooms set removed = true where id = $1`, id)
	if err != nil {
		return errors.New(fmt.Sprintf("can't remove from database product (id: %d)!", id))
	}
	return nil
}

func (s *Service) RoomByID(id int64, pool *pgxpool.Pool) (room Rooms, err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return Rooms{}, errors.New("can't connect to database!")
	}
	defer conn.Release()
	err = conn.QueryRow(context.Background(), `select id, status, timestart, timestop, filename  from room where id=$1`,
		id).Scan(&room.Id, &room.Status, &room.TimeStart, &room.TimeStop, &room.FileName)
	if err != nil {
		return Rooms{}, errors.New(fmt.Sprintf("can't remove from database burger (id: %d)!", id))
	}
	return
}