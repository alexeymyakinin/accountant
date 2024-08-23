package main

import (
	"accountant/internal/api"
	"accountant/internal/config"
	"accountant/internal/dependencies"
	"context"
	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(context.Background(), config.Get().DbUrl)
	if err != nil {
		log.Fatal(err)
	}

	app := api.New(dependencies.New(pool, config.Get()))
	log.Fatal(app.Listen(":8000"))
}
