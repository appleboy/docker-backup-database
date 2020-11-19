package main

import (
	"backup/pkg/config"

	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "database.driver",
			Value:       "postgres",
			Usage:       "database driver",
			EnvVars:     []string{"PLUGIN_DATABASE_DRIVER", "INPUT_DATABASE_DRIVER", "DATABASE_DRIVER"},
			Destination: &cfg.Database.Driver,
		},
		&cli.StringFlag{
			Name:        "database.username",
			Usage:       "database username",
			EnvVars:     []string{"PLUGIN_DATABASE_USERNAME", "INPUT_DATABASEUSERNAME", "DATABASE_USERNAME"},
			Destination: &cfg.Database.Username,
		},
		&cli.StringFlag{
			Name:        "database.password",
			Usage:       "database password",
			EnvVars:     []string{"PLUGIN_DATABASE_PASSWORD", "INPUT_DATABASE_PASSWORD", "DATABASE_PASSWORD"},
			Destination: &cfg.Database.Password,
		},
		&cli.StringFlag{
			Name:        "database.name",
			Usage:       "database name",
			EnvVars:     []string{"PLUGIN_DATABASE_NAME", "INPUT_DATABASE_NAME", "DATABASE_NAME"},
			Destination: &cfg.Database.Name,
		},
		&cli.StringFlag{
			Name:        "database.host",
			Value:       "localhost:5432",
			Usage:       "database host",
			EnvVars:     []string{"PLUGIN_DATABASE_HOST", "INPUT_DATABASE_HOST", "DATABASE_HOST"},
			Destination: &cfg.Database.Host,
		},
		&cli.StringFlag{
			Name:        "database.schema",
			Value:       "public",
			Usage:       "database schema",
			EnvVars:     []string{"PLUGIN_DATABASE_SCHEMA", "INPUT_DATABASE_SCHEMA", "DATABASE_SCHEMA"},
			Destination: &cfg.Database.Schema,
		},

		// storage
		&cli.StringFlag{
			Name:        "storage.driver",
			Value:       "s3",
			Usage:       "storage driver",
			EnvVars:     []string{"PLUGIN_STORAGE_DRIVER", "INPUT_STORAGE_DRIVER", "STORAGE_DRIVER"},
			Destination: &cfg.Storage.Driver,
		},
		&cli.StringFlag{
			Name:        "storage.access_id",
			Value:       "s3",
			Usage:       "storage access id",
			EnvVars:     []string{"PLUGIN_STORAGE_ACCESSID", "INPUT_STORAGE_ACCESSID", "STORAGE_ACCESSID"},
			Destination: &cfg.Storage.AccessID,
		},
		&cli.StringFlag{
			Name:        "storage.secret_key",
			Value:       "s3",
			Usage:       "storage secret key",
			EnvVars:     []string{"PLUGIN_STORAGE_SECRETKEY", "INPUT_STORAGE_SECRETKEY", "STORAGE_SECRETKEY"},
			Destination: &cfg.Storage.SecretKey,
		},
		&cli.StringFlag{
			Name:        "storage.Endpoint",
			Value:       "s3",
			Usage:       "storage Endpoint",
			EnvVars:     []string{"PLUGIN_STORAGE_ENDPOINT", "INPUT_STORAGE_ENDPOINT", "STORAGE_ENDPOINT"},
			Destination: &cfg.Storage.Endpoint,
		},
		&cli.StringFlag{
			Name:        "storage.Bucket",
			Value:       "s3",
			Usage:       "storage Bucket",
			EnvVars:     []string{"PLUGIN_STORAGE_BUCKET", "INPUT_STORAGE_BUCKET", "STORAGE_BUCKET"},
			Destination: &cfg.Storage.Bucket,
		},
		&cli.StringFlag{
			Name:        "storage.Region",
			Value:       "s3",
			Usage:       "storage Region",
			EnvVars:     []string{"PLUGIN_STORAGE_REGION", "INPUT_STORAGE_REGION", "STORAGE_REGION"},
			Destination: &cfg.Storage.Region,
		},
		&cli.StringFlag{
			Name:        "storage.Path",
			Value:       "s3",
			Usage:       "storage Path",
			EnvVars:     []string{"PLUGIN_STORAGE_PATH", "INPUT_STORAGE_PATH", "STORAGE_PATH"},
			Destination: &cfg.Storage.Path,
		},
		&cli.BoolFlag{
			Name:        "storage.SSL",
			Usage:       "storage SSL",
			EnvVars:     []string{"PLUGIN_STORAGE_SSL", "INPUT_STORAGE_SSL", "STORAGE_SSL"},
			Destination: &cfg.Storage.SSL,
		},
	}
}
