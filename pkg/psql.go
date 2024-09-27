package pkg

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/fleimkeipa/tickets-api/models"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NewPSQLClient initializes and returns a PostgreSQL client using pg package.
func NewPSQLClient() *pg.DB {
	opts := pg.Options{
		Database: viper.GetString("database.name"),
		User:     viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		Addr:     viper.GetString("database.addr"),
	}
	db := pg.Connect(&opts)

	if err := createTables(db); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	return db
}

// createTables creates tables for the provided models.
func createTables(db *pg.DB) error {
	models := []interface{}{
		(*models.Ticket)(nil),
	}

	for _, model := range models {
		opts := &orm.CreateTableOptions{
			IfNotExists: true, // Ensures the table is created only if it doesn't exist.
		}

		if err := db.Model(model).CreateTable(opts); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

// GetTestInstance starts a PostgreSQL container for testing and returns a connected pg.DB client along with a cleanup function.
func GetTestInstance(ctx context.Context) (*pg.DB, func()) {
	const mongoVersion = "17.0"
	const port = "5432"

	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("postgres:%s", mongoVersion),
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", port)},
		WaitingFor:   wait.ForListeningPort(port), // Wait until the port is ready
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "test_db",
		},
		Cmd: []string{"postgres", "-c", "fsync=off"}, // Disable fsync for performance in tests
	}
	psqlClient, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("an error occurred while starting postgres container! error details: %v", err)
	}

	psqlPort, err := psqlClient.MappedPort(ctx, port)
	if err != nil {
		log.Fatalf("an error occurred while getting postgres port! error details: %v", err)
	}

	after, _ := strings.CutPrefix(psqlPort.Port(), "/")

	dbAddr := fmt.Sprintf("localhost:%s", after)
	opts := pg.Options{
		User:     "postgres",
		Password: "password",
		Database: "test_db",
		Addr:     dbAddr,
	}
	client := pg.Connect(&opts)

	if err := createTestTables(client); err != nil {
		log.Fatalf("Failed to create test schema: %v", err)
	}

	// Return the client and a cleanup function
	return client, func() {
		client.Close()
		psqlClient.Terminate(ctx)
	}
}

// createTestTables creates temporary test tables for the provided models.
func createTestTables(db *pg.DB) error {
	models := []interface{}{
		(*models.Ticket)(nil),
	}

	for _, model := range models {
		opts := orm.CreateTableOptions{
			Temp:        true, // Creates a temporary table for testing purposes.
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
