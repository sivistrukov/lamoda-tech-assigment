package config

import (
	"os"
	"strconv"
)

const defaultDbMaxConn = 10

type App struct {
	HTTP     HTTP
	Database SQL
}

type HTTP struct {
	Host string
	Port string
}

type SQL struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	MaxConn  int
}

// New returns app config.
func New() App {
	http := newHttpConfig()
	sql := newSqlConfig()

	return App{
		HTTP:     http,
		Database: sql,
	}
}

func newHttpConfig() HTTP {
	return HTTP{
		Host: os.Getenv("SERVER_HOST"),
		Port: os.Getenv("SERVER_PORT"),
	}
}

func newSqlConfig() SQL {
	var maxConn int
	if len(os.Getenv("")) > 0 {
		maxConn, _ = strconv.Atoi(os.Getenv("PG_MAX_OPEN_CONN"))
	} else {
		maxConn = defaultDbMaxConn
	}

	return SQL{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		MaxConn:  maxConn,
	}
}
