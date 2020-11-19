package main

import (
	"log"
	"os"
	"strconv"
	"time"

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

	app := cli.NewApp()
	app.Name = "Backup Database"
	app.Usage = "Docker image to periodically backup a your database"
	app.Copyright = "Copyright (c) " + strconv.Itoa(time.Now().Year()) + " Bo-Yi Wu"
	app.Authors = []*cli.Author{
		{
			Name:  "Bo-Yi Wu",
			Email: "appleboy.tw@gmail.com",
		},
	}
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "host, H",
			Usage:   "Server host",
			EnvVars: []string{"PLUGIN_HOST", "SCP_HOST", "SSH_HOST", "HOST", "INPUT_HOST"},
		},
		&cli.StringFlag{
			Name:    "port, P",
			Value:   "22",
			Usage:   "Server port, default to 22",
			EnvVars: []string{"PLUGIN_PORT", "SCP_PORT", "SSH_PORT", "PORT", "INPUT_PORT"},
		},
		&cli.StringFlag{
			Name:    "username, u",
			Usage:   "Server username",
			EnvVars: []string{"PLUGIN_USERNAME", "PLUGIN_USER", "SCP_USERNAME", "SSH_USERNAME", "USERNAME", "INPUT_USERNAME"},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	return nil
}
