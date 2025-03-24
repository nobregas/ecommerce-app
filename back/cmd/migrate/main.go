package main

import (
	"github.com/nobregas/ecommerce-mobile-back/config"
	"github.com/nobregas/ecommerce-mobile-back/internal/app"
	"log"
	"os"
	"strconv"

	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := app.NewMySQLStorage(mysqlCfg.Config{
		User:                 configs.Envs.DB_USER,
		Passwd:               configs.Envs.DB_PASSWORD,
		Addr:                 configs.Envs.DB_ADDRESS,
		DBName:               configs.Envs.DB_NAME,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "force" {
		if len(os.Args) < 3 {
			log.Fatal("Missing version argument for force command")
		}
		version := os.Args[2]
		v, err := strconv.Atoi(version)
		if err != nil {
			log.Fatal("Invalid version number")
		}
		if err := m.Force(v); err != nil {
			log.Fatal(err)
		}
	}
}
