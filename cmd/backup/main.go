package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"backup/pkg/config"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

// Version set at compile-time
var (
	Version string
)

func main() {
	// Load env-file if it exists first
	if filename, found := os.LookupEnv("PLUGIN_ENV_FILE"); found {
		_ = godotenv.Load(filename)
	}

	if _, err := os.Stat("/run/drone/env"); err == nil {
		godotenv.Overload("/run/drone/env")
	}

	cfg := &config.Config{}
	app := &cli.App{
		Name:      "docker-backup-datavase",
		Usage:     "Docker image to periodically backup a your database",
		Copyright: "Copyright (c) " + strconv.Itoa(time.Now().Year()) + " Bo-Yi Wu",
		Version:   Version,
		Flags:     settingsFlags(cfg),
		Action:    run(cfg),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(cfg *config.Config) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		return nil
	}
}
