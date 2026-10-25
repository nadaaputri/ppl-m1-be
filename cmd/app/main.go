package main

import (
	"log"

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

	log.Println("Koneksi PostgreSQL berhasil")

	router := v1.NewRouter(db, cfg)

	address := ":" + cfg.AppPort
	log.Printf("Server berjalan di http://localhost%s", address)

	if err := router.Run(address); err != nil {
		log.Fatal(err)
	}
}