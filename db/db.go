package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"todo-app/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

var (
	ErrNotFound    = errors.New("not found")
	ErrTooManyRows = errors.New("too many rows")
	ErrOther       = errors.New("other error")
)

func Connect() {
	var err error
	pool, err = retry(10, 1, func() (*pgxpool.Pool, error) {
		return pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	if _, err := pool.Exec(context.Background(), "create table if not exists todos (id serial primary key, description varchar, done boolean)"); err != nil {
		log.Fatalf("Unable to create table: %v\n", err)
	}
}

func retry[T any](attempts int, sleep int, f func() (T, error)) (result T, err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Println("retrying after error:", err)
			time.Sleep(time.Duration(sleep) * time.Second)
			sleep *= 2
		}
		result, err = f()
		if err == nil {
			return result, nil
		}
	}
	return result, fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

func Close() {
	pool.Close()
}

func Insert(toDo *dto.ToDo) (*dto.SavedToDo, error) {
	var id string
	if err := pool.QueryRow(context.Background(), "insert into todos(description, done) values ($1, $2) returning id", toDo.Description, toDo.Done).Scan(&id); err != nil {
		return nil, err
	}
	return &dto.SavedToDo{Id: id, ToDo: *toDo}, nil
}

func Update(toDo *dto.SavedToDo) error {
	_, err := pool.Exec(context.Background(), "update todos set description = $2, done = $3 where id = $1", toDo.Id, toDo.Description, toDo.Done)
	return err
}

func Upsert(toDo *dto.SavedToDo) error {
	_, err := pool.Exec(context.Background(), "upsert into todos values ($1, $2, $3)", toDo.Id, toDo.Description, toDo.Done)
	return err
}

func Delete(id string) error {
	_, err := pool.Exec(context.Background(), "delete from todos where id = $1", id)
	return err
}

func GetOne(id string) (*dto.SavedToDo, error) {
	var description string
	var done bool
	err := pool.QueryRow(context.Background(), "select description, done from todos where id=$1", id).Scan(&description, &done)
	switch err {
	case nil:
		return &dto.SavedToDo{Id: id, ToDo: dto.ToDo{Description: description, Done: done}}, nil
	case pgx.ErrNoRows:
		return nil, ErrNotFound
	case pgx.ErrTooManyRows:
		return nil, ErrTooManyRows
	default:
		return nil, ErrOther
	}
}

func GetAll() (*[]dto.SavedToDo, error) {
	rows, err := pool.Query(context.Background(), "select id, description, done from todos")
	if err != nil {
		log.Fatalf("Error getting todos: %v\n", err)
	}
	todos := make([]dto.SavedToDo, 0)
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Printf("Error getting values: %v\n", err)
			return nil, err
		}
		todos = append(todos, dto.SavedToDo{Id: strconv.FormatInt(values[0].(int64), 10), ToDo: dto.ToDo{Description: values[1].(string), Done: values[2].(bool)}})
	}
	return &todos, nil
}
