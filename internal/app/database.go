package app

import (
	"database/sql"
	"errors"
	"fmt"
	"lamoda-tech-assigment/internal/config"
	"log"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

func newPostgresqlConnect(cfg config.SQL) (*sql.DB, error) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("host=%s port=%s ", cfg.Host, cfg.Port))
	builder.WriteString(fmt.Sprintf("user=%s password=%s ", cfg.User, cfg.Password))
	builder.WriteString(fmt.Sprintf("dbname=%s ", cfg.DBName))
	builder.WriteString("sslmode=disable ")

	params := builder.String()

	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxConn)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func runMigrations(cfg config.SQL) {
	databaseUrl := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	migrationsDir := "file://migrations"

	m, err := migrate.New(migrationsDir, databaseUrl)
	if err != nil {
		log.Print(err)
		return
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Print(err)
		return
	}

	version, dirty, err := m.Version()
	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("Applied migrations: %d, Diry: %t\n", version, dirty)
}
