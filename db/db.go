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

func Connect() {
	var err error
	pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
}

func Init() {
	if _, err := pool.Exec(context.Background(), "create table if not exists todos (id integer primary key, description varchar, done boolean)"); err != nil {
		log.Fatalf("Unable to create table: %v\n", err)
	}

	Upsert(&dto.SavedToDo{Id: 1, ToDo: dto.ToDo{Description: "write app", Done: false}})
}

func Close() {
	pool.Close()
}

func Insert(toDo *dto.ToDo) (*dto.SavedToDo, error) {
	// TODO make this concurrent
	var id int
	err := pool.QueryRow(context.Background(), "select coalesce(max(id), 0) from todos").Scan(&id)
	switch err {
	case nil:
		id = id + 1
	case pgx.ErrNoRows:
		id = 0
	default:
		log.Fatalf("Error getting max id: %v\n", err)
	}
	_, err = pool.Exec(context.Background(), "insert into todos values ($1, $2, $3)", id, toDo.Description, toDo.Done)
	return &dto.SavedToDo{Id: id, ToDo: *toDo}, err
}

func Update(toDo *dto.SavedToDo) error {
	_, err := pool.Exec(context.Background(), "update todos set description = $2, done = $3 where id = $1", toDo.Id, toDo.Description, toDo.Done)
	return err
}

func Upsert(toDo *dto.SavedToDo) error {
	_, err := pool.Exec(context.Background(), "upsert into todos values ($1, $2, $3)", toDo.Id, toDo.Description, toDo.Done)
	return err
}

func Delete(id int) error {
	_, err := pool.Exec(context.Background(), "delete from todos where id = $1", id)
	return err
}

func GetOne(id int) (*dto.SavedToDo, error) {
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
		todos = append(todos, dto.SavedToDo{Id: int(values[0].(int64)), ToDo: dto.ToDo{Description: values[1].(string), Done: values[2].(bool)}})
	}
	return &todos, nil
}
