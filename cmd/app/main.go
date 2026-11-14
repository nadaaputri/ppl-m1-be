package main

import (
	"log"

	"ppl-m1-be/config"
	"ppl-m1-be/internal/app"
	v1 "ppl-m1-be/internal/controller/http/v1"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	db, err := app.ConnectDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Println("Koneksi PostgreSQL berhasil")
	address := ":" + cfg.AppPort
	log.Printf("Server berjalan di http://localhost%s", address)
	return v1.NewRouter(db, cfg).Run(address)
}
