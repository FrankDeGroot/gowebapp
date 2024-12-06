package db

import (
	"context"
	"errors"
	"log"
	"os"
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

func Init() {
	var err error
	pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	if _, err := pool.Exec(context.Background(), "create table if not exists todo (id integer primary key, description varchar, done boolean)"); err != nil {
		log.Fatalf("Unable to create table: %v\n", err)
	}

	Upsert(&dto.ToDo{Id: 1, Description: "write app", Done: false})
}

func Close() {
	pool.Close()
}

func Upsert(toDo *dto.ToDo) error {
	_, err := pool.Exec(context.Background(), "upsert into todo values ($1, $2, $3)", toDo.Id, toDo.Description, toDo.Done)
	return err
}

func Get(id int) (*dto.ToDo, error) {
	var description string
	var done bool
	err := pool.QueryRow(context.Background(), "select description, done from todo where id=$1", id).Scan(&description, &done)
	switch err {
	case nil:
		return &dto.ToDo{Id: id, Description: description, Done: done}, nil
	case pgx.ErrNoRows:
		return nil, ErrNotFound
	case pgx.ErrTooManyRows:
		return nil, ErrTooManyRows
	default:
		return nil, ErrOther
	}
}

func GetAll() (*[]dto.ToDo, error) {
	rows, err := pool.Query(context.Background(), "select id, description, done from todo")
	if err != nil {
		log.Fatalf("Error getting todos: %v\n", err)
	}

	todos := make([]dto.ToDo, 0)
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Printf("Error getting values: %v\n", err)
			return nil, err
		}
		todos = append(todos, dto.ToDo{Id: int(values[0].(int64)), Description: values[1].(string), Done: values[2].(bool)})
	}
	return &todos, nil
}
