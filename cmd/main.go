package main

import (
	_ "lamoda-tech-assigment/docs"
	"lamoda-tech-assigment/internal/app"
	"lamoda-tech-assigment/internal/config"
	"log"
)

// @title		Lamoda Tech assigment
// @version	1.0
// @host		localhost:8080
// @BasePath	/api
func main() {
	cfg := config.New()

	a, finalizeFunc := app.MustNew(cfg)
	defer finalizeFunc()

	if err := a.Run(); err != nil {
		log.Println(err)
	}
}
