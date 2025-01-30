package config

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func ConnectDB() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connStr := "postgres://postgres:Exxomatic%402024@213.199.41.58:5000/exxomatic_ev"
	DB, err = pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal("Unable to connect to the database:", err)
	}
}
