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
			Value:       "postgres",
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
			Name:        "database.opts",
			Usage:       "database options",
			EnvVars:     []string{"PLUGIN_DATABASE_OPTS", "INPUT_DATABASE_OPTS", "DATABASE_OPTS"},
			Destination: &cfg.Database.Opts,
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
			Usage:       "storage access id",
			EnvVars:     []string{"PLUGIN_ACCESS_KEY_ID", "INPUT_ACCESS_KEY_ID", "ACCESS_KEY_ID"},
			Destination: &cfg.Storage.AccessID,
		},
		&cli.StringFlag{
			Name:        "storage.secret_key",
			Usage:       "storage secret key",
			EnvVars:     []string{"PLUGIN_SECRET_ACCESS_KEY", "INPUT_SECRET_ACCESS_KEY", "SECRET_ACCESS_KEY"},
			Destination: &cfg.Storage.SecretKey,
		},
		&cli.StringFlag{
			Name:        "storage.endpoint",
			Value:       "s3.amazonaws.com",
			Usage:       "storage endpoint",
			EnvVars:     []string{"PLUGIN_STORAGE_ENDPOINT", "INPUT_STORAGE_ENDPOINT", "STORAGE_ENDPOINT"},
			Destination: &cfg.Storage.Endpoint,
		},
		&cli.StringFlag{
			Name:        "storage.bucket",
			Usage:       "storage bucket",
			EnvVars:     []string{"PLUGIN_STORAGE_BUCKET", "INPUT_STORAGE_BUCKET", "STORAGE_BUCKET"},
			Destination: &cfg.Storage.Bucket,
		},
		&cli.StringFlag{
			Name:        "storage.region",
			Value:       "ap-northeast-1",
			Usage:       "storage region",
			EnvVars:     []string{"PLUGIN_STORAGE_REGION", "INPUT_STORAGE_REGION", "STORAGE_REGION"},
			Destination: &cfg.Storage.Region,
		},
		&cli.StringFlag{
			Name:        "storage.path",
			Value:       "backup",
			Usage:       "storage path",
			EnvVars:     []string{"PLUGIN_STORAGE_PATH", "INPUT_STORAGE_PATH", "STORAGE_PATH"},
			Destination: &cfg.Storage.Path,
		},
		&cli.BoolFlag{
			Name:        "storage.ssl",
			Usage:       "storage ssl",
			EnvVars:     []string{"PLUGIN_STORAGE_SSL", "INPUT_STORAGE_SSL", "STORAGE_SSL"},
			Destination: &cfg.Storage.SSL,
		},
		&cli.StringFlag{
			Name:        "storage.dump_name",
			Usage:       "storage dump name",
			EnvVars:     []string{"PLUGIN_STORAGE_DUMP_NAME", "INPUT_STORAGE_DUMP_NAME", "STORAGE_DUMP_NAME"},
			Destination: &cfg.Storage.DumpName,
			Value:       "dump.sql.gz",
		},
		&cli.BoolFlag{
			Name:        "storage.insecure_skip_verify",
			Usage:       "storage insecure skip verify",
			EnvVars:     []string{"PLUGIN_STORAGE_INSECURE_SKIP_VERIFY", "INPUT_STORAGE_INSECURE_SKIP_VERIFY", "STORAGE_INSECURE_SKIP_VERIFY"},
			Destination: &cfg.Storage.InsecureSkipVerify,
		},

		// SCHEDULE
		&cli.StringFlag{
			Name:        "schedule",
			Usage:       "cron spec format",
			EnvVars:     []string{"PLUGIN_schedule", "INPUT_schedule", "schedule"},
			Destination: &cfg.Server.Schedule,
		},
	}
}
