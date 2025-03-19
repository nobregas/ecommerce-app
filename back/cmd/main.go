package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/nobregas/ecommerce-mobile-back/cmd/api"
	configs "github.com/nobregas/ecommerce-mobile-back/config"
	"github.com/nobregas/ecommerce-mobile-back/db"
)

func main() {
	cfg := mysql.Config{
		User:                 configs.Envs.DB_USER,
		Passwd:               configs.Envs.DB_PASSWORD,
		Addr:                 configs.Envs.DB_ADDRESS,
		DBName:               configs.Envs.DB_NAME,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := db.NewMySQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected")
}
