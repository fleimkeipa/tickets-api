package tests

import (
	"log"

	"github.com/fleimkeipa/tickets-api/config"

	"github.com/go-pg/pg"
	_ "github.com/lib/pq"
)

var (
	test_db     *pg.DB
	terminateDB = func() {}
)

func loadEnv() {
	if err := config.LoadEnv("../"); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	log.Println("configuration loaded successfully")
}

func addTempData(data interface{}) error {
	_, err := test_db.Model(data).Insert()
	if err != nil {
		return err
	}

	return nil
}

func clearTable() error {
	_, err := test_db.Exec("TRUNCATE tickets; DELETE FROM tickets")
	if err != nil {
		return err
	}

	return nil
}
