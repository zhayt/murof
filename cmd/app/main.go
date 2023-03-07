package main

import (
	"github.com/zhayt/clean-arch-tmp-forum/config"
	"github.com/zhayt/clean-arch-tmp-forum/internal/app"
	"log"
)

const config_path = "./config/config.json"

func main() {
	cfg, err := config.NewConfig(config_path)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
