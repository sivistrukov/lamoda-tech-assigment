package main

import (
	"lamoda-tech-assigment/internal/app"
	"lamoda-tech-assigment/internal/config"
	"log"
)

func main() {
	cfg := config.New()

	a, finalizeFunc := app.MustNew(cfg)
	defer finalizeFunc()

	if err := a.Run(); err != nil {
		log.Println(err)
	}
}
