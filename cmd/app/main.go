package main

import (
	"github.com/mxmrykov/L0/config"
	"github.com/mxmrykov/L0/internal/app"
	"log"
)

func main() {

	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Error at starting: %v", err)
	}

	app.Run(cfg)
}
