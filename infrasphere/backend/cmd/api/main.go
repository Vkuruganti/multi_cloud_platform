package main

import (
	"log"
	"net/http"

	"github.com/infrasphere/control-plane/backend/internal/api"
	"github.com/infrasphere/control-plane/backend/internal/config"
	"github.com/infrasphere/control-plane/backend/internal/database"
	"github.com/infrasphere/control-plane/backend/internal/providers/mock"
)

func main() {
	cfg := config.Load()
	store := database.NewSeededStore(mock.New())
	log.Printf("InfraSphere API listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, api.New(cfg, store)))
}

