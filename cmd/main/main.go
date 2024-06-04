package main

import (
	"flag"
	"fmt"
	"log"
	"wbLvL0/internal/app"
	"wbLvL0/internal/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "cfg",
		"./internal/config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatal(fmt.Errorf("error while reading config: %s", err.Error()))
	}

	app.Run(cfg)
}
