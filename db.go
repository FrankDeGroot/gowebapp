package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDb() {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	Pool = pool

	if _, err := Pool.Exec(context.Background(), "create table if not exists book (id integer primary key, title varchar)"); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
	}

	if _, err := Pool.Exec(context.Background(), "upsert into book values (1, 'some book')"); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert row: %v\n", err)
	}
}
