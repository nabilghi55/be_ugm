package main

import (
	"log"

	"samplegadai-backend/internal/config"
	"samplegadai-backend/internal/database"
)

func main() {
	log.Println("=== MEMULAI RESET DATABASE GADAI SYARIAH DARI 0 ===")
	cfg := config.LoadConfig()
	db := database.InitDB(cfg)

	database.ResetAndSeed(db)
	log.Println("=== RESET DATABASE SELESAI DENGAN SUKSES! DATA BERSIH & TERVALIDASI ===")
}
