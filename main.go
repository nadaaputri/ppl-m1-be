package main

import (
	"log"
	"os"

	"ppl-m1-be/config"
	"ppl-m1-be/internal/app"
	v1 "ppl-m1-be/internal/controller/http/v1"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := app.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := v1.NewRouter(db, cfg).Run(":" + listenPort(cfg)); err != nil {
		log.Fatal(err)
	}
}

func listenPort(cfg config.Config) string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return cfg.AppPort
}