// cmd/auth-server/main.go
package main

import (
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/Roflan4eg/auth-serivce/internal/app"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"github.com/Roflan4eg/auth-serivce/internal/storage"
)

func main() {
	cfg, err := config.LoadFromFile("config.yaml")
	if err != nil {
		panic(err)
	}
	log := logger.NewLogger(cfg.App)

	if cfg.App.AutoMigrate {
		if err = storage.RunMigrations(cfg.Database, log); err != nil {
			panic(err)
		}
	}

	newApp := app.New(cfg, log)

	if err = newApp.Setup(); err != nil {
		log.Error(err.Error())
		panic(err)
	}
	if err = newApp.Start(); err != nil {
		log.Error(err.Error())
		panic(err)
	}
}
