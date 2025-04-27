package db

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"
	"todo-app/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

var (
	ErrNotFound    = errors.New("not found")
	ErrTooManyRows = errors.New("too many rows")
	ErrOther       = errors.New("other error")
)

func Open() (*DB, error) {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	db := &DB{pool: pool}

	if _, err := retry(10, 1*time.Second, func() (pgconn.CommandTag, error) {
		return db.pool.Exec(context.Background(), "create table if not exists tasks (id serial primary key, description varchar, done boolean)")
	}); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) Insert(task *dto.Task) (*dto.SavedTask, error) {
	var id string
	if err := db.pool.QueryRow(context.Background(), "insert into tasks(description, done) values ($1, $2) returning id", task.Description, task.Done).Scan(&id); err != nil {
		return nil, err
	}
	return &dto.SavedTask{Id: id, Task: *task}, nil
}

func (db *DB) Update(task *dto.SavedTask) error {
	_, err := db.pool.Exec(context.Background(), "update tasks set description = $2, done = $3 where id = $1", task.Id, task.Description, task.Done)
	return err
}

func (db *DB) Upsert(task *dto.SavedTask) error {
	_, err := db.pool.Exec(context.Background(), "upsert into tasks values ($1, $2, $3)", task.Id, task.Description, task.Done)
	return err
}

func (db *DB) Delete(id string) error {
	_, err := db.pool.Exec(context.Background(), "delete from tasks where id = $1", id)
	return err
}

func (db *DB) GetOne(id string) (*dto.SavedTask, error) {
	var description string
	var done bool
	err := db.pool.QueryRow(context.Background(), "select description, done from tasks where id = $1", id).Scan(&description, &done)
	switch err {
	case nil:
		return &dto.SavedTask{Id: id, Task: dto.Task{Description: description, Done: done}}, nil
	case pgx.ErrNoRows:
		return nil, ErrNotFound
	case pgx.ErrTooManyRows:
		return nil, ErrTooManyRows
	default:
		return nil, ErrOther
	}
}

func (db *DB) GetAll() (*[]dto.SavedTask, error) {
	rows, err := db.pool.Query(context.Background(), "select id, description, done from tasks")
	if err != nil {
		return nil, err
	}
	tasks := make([]dto.SavedTask, 0)
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, dto.SavedTask{Id: strconv.FormatInt(values[0].(int64), 10), Task: dto.Task{Description: values[1].(string), Done: values[2].(bool)}})
	}
	return &tasks, nil
}
