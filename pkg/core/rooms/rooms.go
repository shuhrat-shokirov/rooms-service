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
	Name   string `json:"name"`
	Status bool `json:"status"`
	FileName string `json:"file_name"`
	Removed string `json:"removed"`
}

func (s *Service) AllRooms (pool *pgxpool.Pool) (list []Rooms ,err error)  {
	rows, err := pool.Query(context.Background(), `SELECT id, name,  status, filename FROM rooms where removed=false;`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := Rooms{}
		err := rows.Scan(&item.Id, &item.Name, &item.Status, &item.FileName)
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
	_, err = pool.Exec(context.Background(), `INSERT INTO rooms(name, filename)
VALUES ($1, $2);`,room.Name, room.FileName)
	if err != nil {
		return
	}
	return nil
}

func (s *Service) RemoveByID(id int64, pool *pgxpool.Pool) (err error) {
	_, err = pool.Exec(context.Background(), `update rooms set removed = true where id = $1`, id)
	if err != nil {
		return errors.New(fmt.Sprintf("can't remove from database product (id: %d)!", id))
	}
	return nil
}

func (s *Service) RoomByID(id int64, pool *pgxpool.Pool) (room Rooms, err error) {
	err = pool.QueryRow(context.Background(), `select name, status, filename  from room where id=$1`,
		id).Scan(&room.Name, &room.Status, &room.FileName)
	if err != nil {
		return Rooms{}, errors.New(fmt.Sprintf("can't list from database rooms (id: %d)!", id))
	}
	return
}

func (s *Service) LockRoomByID(id int64, pool *pgxpool.Pool) (err error) {
	_, err = pool.Exec(context.Background(), `update rooms set status = true where id = $1`, id)
	if err != nil {
		return  errors.New(fmt.Sprintf("can't remove from database rooms (id: %d)!", id))
	}
	return
}

func (s *Service) UnLockRoomByID(id int64, pool *pgxpool.Pool) (err error) {
	_, err = pool.Exec(context.Background(), `update rooms set status = false where id = $1`, id)
	if err != nil {
		return  errors.New(fmt.Sprintf("can't remove from database rooms (id: %d)!", id))
	}
	return
}

func (s *Service) AllRoomsUnlocked (pool *pgxpool.Pool) (list []Rooms ,err error)  {
	rows, err := pool.Query(context.Background(), `SELECT name, status, filename FROM rooms where removed=false and status = false;`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := Rooms{}
		err := rows.Scan(&item.Name, &item.Status, &item.FileName)
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

func (s *Service) AllRoomsLocked (pool *pgxpool.Pool) (list []Rooms ,err error)  {
	rows, err := pool.Query(context.Background(), `SELECT name, status, filename FROM rooms where removed=false and status != false;`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := Rooms{}
		err := rows.Scan(&item.Id, &item.Status,  &item.FileName)
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

