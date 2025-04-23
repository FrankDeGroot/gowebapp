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
		return db.pool.Exec(context.Background(), "create table if not exists todos (id serial primary key, description varchar, done boolean)")
	}); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) Insert(todo *dto.Todo) (*dto.SavedTodo, error) {
	var id string
	if err := db.pool.QueryRow(context.Background(), "insert into todos(description, done) values ($1, $2) returning id", todo.Description, todo.Done).Scan(&id); err != nil {
		return nil, err
	}
	return &dto.SavedTodo{Id: id, Todo: *todo}, nil
}

func (db *DB) Update(todo *dto.SavedTodo) error {
	_, err := db.pool.Exec(context.Background(), "update todos set description = $2, done = $3 where id = $1", todo.Id, todo.Description, todo.Done)
	return err
}

func (db *DB) Upsert(todo *dto.SavedTodo) error {
	_, err := db.pool.Exec(context.Background(), "upsert into todos values ($1, $2, $3)", todo.Id, todo.Description, todo.Done)
	return err
}

func (db *DB) Delete(id string) error {
	_, err := db.pool.Exec(context.Background(), "delete from todos where id = $1", id)
	return err
}

func (db *DB) GetOne(id string) (*dto.SavedTodo, error) {
	var description string
	var done bool
	err := db.pool.QueryRow(context.Background(), "select description, done from todos where id = $1", id).Scan(&description, &done)
	switch err {
	case nil:
		return &dto.SavedTodo{Id: id, Todo: dto.Todo{Description: description, Done: done}}, nil
	case pgx.ErrNoRows:
		return nil, ErrNotFound
	case pgx.ErrTooManyRows:
		return nil, ErrTooManyRows
	default:
		return nil, ErrOther
	}
}

func (db *DB) GetAll() (*[]dto.SavedTodo, error) {
	rows, err := db.pool.Query(context.Background(), "select id, description, done from todos")
	if err != nil {
		return nil, err
	}
	todos := make([]dto.SavedTodo, 0)
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		todos = append(todos, dto.SavedTodo{Id: strconv.FormatInt(values[0].(int64), 10), Todo: dto.Todo{Description: values[1].(string), Done: values[2].(bool)}})
	}
	return &todos, nil
}
