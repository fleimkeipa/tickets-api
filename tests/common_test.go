package tests

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/fleimkeipa/tickets-api/config"
	"github.com/fleimkeipa/tickets-api/models"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var test_db *pg.DB

const testDBName = "test_db"

func loadEnv() {
	if err := config.LoadEnv("../"); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	log.Println("configuration loaded successfully")
}

func init() {
	loadEnv()

	addr := strings.Split(viper.GetString("database.addr"), ":")
	if len(addr) < 2 {
		viper.Set("database.host", "localhost")
		viper.Set("database.port", "5432")
	} else {
		viper.Set("database.host", addr[0])
		viper.Set("database.port", addr[1])
	}

	createTestDB()

	test_db = initTestDBClient()

	if err := createTestSchema(test_db); err != nil {
		log.Fatal(err)
	}
}

func initTestDBClient() *pg.DB {
	opts := pg.Options{
		Database: testDBName,
		User:     viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		Addr:     viper.GetString("database.addr"),
	}
	db := pg.Connect(&opts)

	return db
}

func createTestDB() {
	conninfo := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetString("database.port"),
	)
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		log.Fatal(err)
	}

	query := fmt.Sprintf(`
		DROP DATABASE IF EXISTS %s;
	`, testDBName)
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = fmt.Sprintf(`
		CREATE DATABASE %s;
	`, testDBName)
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

// createSchema creates database schema for Ticket
func createTestSchema(db *pg.DB) error {
	models := []interface{}{
		(*models.Ticket)(nil),
	}

	for _, model := range models {
		opts := orm.CreateTableOptions{
			Temp:        true,
			IfNotExists: true,
		}
		err := db.
			Model(model).
			CreateTable(&opts)
		if err != nil {
			return err
		}
	}

	return nil
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
