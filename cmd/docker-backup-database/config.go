package main

import (
	"github.com/appleboy/docker-backup-database/pkg/config"

	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "database.driver",
			Value:       "postgres",
			Usage:       "Specify the database driver (e.g., postgres, mysql)",
			EnvVars:     []string{"PLUGIN_DATABASE_DRIVER", "INPUT_DATABASE_DRIVER", "DATABASE_DRIVER"},
			Destination: &cfg.Database.Driver,
		},
		&cli.StringFlag{
			Name:        "database.username",
			Usage:       "Specify the database username",
			EnvVars:     []string{"PLUGIN_DATABASE_USERNAME", "INPUT_DATABASE_USERNAME", "DATABASE_USERNAME"},
			Destination: &cfg.Database.Username,
		},
		&cli.StringFlag{
			Name:        "database.password",
			Usage:       "Specify the database password",
			EnvVars:     []string{"PLUGIN_DATABASE_PASSWORD", "INPUT_DATABASE_PASSWORD", "DATABASE_PASSWORD"},
			Destination: &cfg.Database.Password,
		},
		&cli.StringFlag{
			Name:        "database.name",
			Usage:       "Specify the database name",
			Value:       "postgres",
			EnvVars:     []string{"PLUGIN_DATABASE_NAME", "INPUT_DATABASE_NAME", "DATABASE_NAME"},
			Destination: &cfg.Database.Name,
		},
		&cli.StringFlag{
			Name:        "database.host",
			Value:       "localhost:5432",
			Usage:       "Specify the database host and port",
			EnvVars:     []string{"PLUGIN_DATABASE_HOST", "INPUT_DATABASE_HOST", "DATABASE_HOST"},
			Destination: &cfg.Database.Host,
		},
		&cli.StringFlag{
			Name:        "database.opts",
			Usage:       "Specify additional database options",
			EnvVars:     []string{"PLUGIN_DATABASE_OPTS", "INPUT_DATABASE_OPTS", "DATABASE_OPTS"},
			Destination: &cfg.Database.Opts,
		},

		// storage
		&cli.StringFlag{
			Name:        "storage.driver",
			Value:       "s3",
			Usage:       "Specify the storage driver (e.g., s3, gcs)",
			EnvVars:     []string{"PLUGIN_STORAGE_DRIVER", "INPUT_STORAGE_DRIVER", "STORAGE_DRIVER"},
			Destination: &cfg.Storage.Driver,
		},
		&cli.StringFlag{
			Name:        "storage.access_id",
			Usage:       "Specify the storage access ID",
			EnvVars:     []string{"PLUGIN_ACCESS_KEY_ID", "INPUT_ACCESS_KEY_ID", "ACCESS_KEY_ID"},
			Destination: &cfg.Storage.AccessID,
		},
		&cli.StringFlag{
			Name:        "storage.secret_key",
			Usage:       "Specify the storage secret key",
			EnvVars:     []string{"PLUGIN_SECRET_ACCESS_KEY", "INPUT_SECRET_ACCESS_KEY", "SECRET_ACCESS_KEY"},
			Destination: &cfg.Storage.SecretKey,
		},
		&cli.StringFlag{
			Name:        "storage.endpoint",
			Value:       "s3.amazonaws.com",
			Usage:       "Specify the storage endpoint URL",
			EnvVars:     []string{"PLUGIN_STORAGE_ENDPOINT", "INPUT_STORAGE_ENDPOINT", "STORAGE_ENDPOINT"},
			Destination: &cfg.Storage.Endpoint,
		},
		&cli.StringFlag{
			Name:        "storage.bucket",
			Usage:       "Specify the storage bucket name",
			EnvVars:     []string{"PLUGIN_STORAGE_BUCKET", "INPUT_STORAGE_BUCKET", "STORAGE_BUCKET"},
			Destination: &cfg.Storage.Bucket,
		},
		&cli.StringFlag{
			Name:        "storage.region",
			Value:       "ap-northeast-1",
			Usage:       "Specify the storage region",
			EnvVars:     []string{"PLUGIN_STORAGE_REGION", "INPUT_STORAGE_REGION", "STORAGE_REGION"},
			Destination: &cfg.Storage.Region,
		},
		&cli.StringFlag{
			Name:        "storage.path",
			Value:       "backup",
			Usage:       "Specify the storage path",
			EnvVars:     []string{"PLUGIN_STORAGE_PATH", "INPUT_STORAGE_PATH", "STORAGE_PATH"},
			Destination: &cfg.Storage.Path,
		},
		&cli.BoolFlag{
			Name:        "storage.ssl",
			Usage:       "Enable or disable SSL for storage",
			EnvVars:     []string{"PLUGIN_STORAGE_SSL", "INPUT_STORAGE_SSL", "STORAGE_SSL"},
			Destination: &cfg.Storage.SSL,
		},
		&cli.StringFlag{
			Name:        "storage.dump_name",
			Usage:       "Specify the name of the storage dump file",
			EnvVars:     []string{"PLUGIN_STORAGE_DUMP_NAME", "INPUT_STORAGE_DUMP_NAME", "STORAGE_DUMP_NAME"},
			Destination: &cfg.Storage.DumpName,
			Value:       "dump.sql.gz",
		},
		&cli.BoolFlag{
			Name:  "storage.insecure_skip_verify",
			Usage: "Skip SSL certificate verification for storage",
			EnvVars: []string{
				"PLUGIN_STORAGE_INSECURE_SKIP_VERIFY",
				"INPUT_STORAGE_INSECURE_SKIP_VERIFY",
				"STORAGE_INSECURE_SKIP_VERIFY",
			},
			Destination: &cfg.Storage.SkipVerify,
		},
		&cli.IntFlag{
			Name:        "storage.days",
			Usage:       "Set the number of days to retain storage files",
			EnvVars:     []string{"PLUGIN_STORAGE_DAYS", "INPUT_STORAGE_DAYS", "STORAGE_DAYS"},
			Destination: &cfg.Storage.Days,
			Value:       0,
		},

		// SCHEDULE
		&cli.StringFlag{
			Name:        "time.schedule",
			Usage:       "Specify a pre-defined schedule or cron expression",
			EnvVars:     []string{"PLUGIN_TIME_SCHEDULE", "INPUT_TIME_SCHEDULE", "TIME_SCHEDULE"},
			Destination: &cfg.Server.Schedule,
		},
		&cli.StringFlag{
			Name:        "time.location",
			Usage:       "Specify the time zone for scheduling (default is local time zone)",
			EnvVars:     []string{"PLUGIN_TIME_LOCATION", "INPUT_TIME_LOCATION", "TIME_LOCATION"},
			Destination: &cfg.Server.Location,
		},

		// File Format
		&cli.StringFlag{
			Name:        "file.prefix",
			Usage:       "Specify the prefix for the file name (default is storage driver name)",
			EnvVars:     []string{"PLUGIN_FILE_PREFIX", "INPUT_FILE_PREFIX", "FILE_PREFIX"},
			Destination: &cfg.File.Prefix,
		},
		&cli.StringFlag{
			Name:        "file.suffix",
			Usage:       "Specify the suffix for the file name",
			EnvVars:     []string{"PLUGIN_FILE_SUFFIX", "INPUT_FILE_SUFFIX", "FILE_SUFFIX"},
			Destination: &cfg.File.Suffix,
		},
		&cli.StringFlag{
			Name:        "file.format",
			Usage:       "Specify the format string for the file name (default is 20060102150405)",
			Value:       "20060102150405",
			EnvVars:     []string{"PLUGIN_FILE_FORMAT", "INPUT_FILE_FORMAT", "FILE_FORMAT"},
			Destination: &cfg.File.Format,
		},
		&cli.StringFlag{
			Name:        "webhook.url",
			Usage:       "Specify the webhook URL",
			EnvVars:     []string{"PLUGIN_WEBHOOK_URL", "INPUT_WEBHOOK_URL", "WEBHOOK_URL"},
			Destination: &cfg.Webhook.URL,
		},
		&cli.BoolFlag{
			Name:        "webhook.insecure",
			Usage:       "Enable or disable insecure mode for webhook",
			EnvVars:     []string{"PLUGIN_WEBHOOK_INSECURE", "INPUT_WEBHOOK_INSECURE", "WEBHOOK_INSECURE"},
			Destination: &cfg.Webhook.Insecure,
		},
	}
}
