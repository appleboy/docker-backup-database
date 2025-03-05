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
			Usage:       "Database type to back up (postgres or mysql)",
			EnvVars:     []string{"PLUGIN_DATABASE_DRIVER", "INPUT_DATABASE_DRIVER", "DATABASE_DRIVER"},
			Destination: &cfg.Database.Driver,
		},
		&cli.StringFlag{
			Name:        "database.username",
			Usage:       "Username for database authentication",
			EnvVars:     []string{"PLUGIN_DATABASE_USERNAME", "INPUT_DATABASE_USERNAME", "DATABASE_USERNAME"},
			Destination: &cfg.Database.Username,
		},
		&cli.StringFlag{
			Name:        "database.password",
			Usage:       "Password for database authentication",
			EnvVars:     []string{"PLUGIN_DATABASE_PASSWORD", "INPUT_DATABASE_PASSWORD", "DATABASE_PASSWORD"},
			Destination: &cfg.Database.Password,
		},
		&cli.StringFlag{
			Name:        "database.name",
			Usage:       "Name of the database to back up",
			Value:       "postgres",
			EnvVars:     []string{"PLUGIN_DATABASE_NAME", "INPUT_DATABASE_NAME", "DATABASE_NAME"},
			Destination: &cfg.Database.Name,
		},
		&cli.StringFlag{
			Name:        "database.host",
			Value:       "localhost:5432",
			Usage:       "Database host address with port (e.g., localhost:5432)",
			EnvVars:     []string{"PLUGIN_DATABASE_HOST", "INPUT_DATABASE_HOST", "DATABASE_HOST"},
			Destination: &cfg.Database.Host,
		},
		&cli.StringFlag{
			Name:        "database.opts",
			Usage:       "Additional database connection options",
			EnvVars:     []string{"PLUGIN_DATABASE_OPTS", "INPUT_DATABASE_OPTS", "DATABASE_OPTS"},
			Destination: &cfg.Database.Opts,
		},

		// storage
		&cli.StringFlag{
			Name:        "storage.driver",
			Value:       "s3",
			Usage:       "Storage service type (s3 or gcs)",
			EnvVars:     []string{"PLUGIN_STORAGE_DRIVER", "INPUT_STORAGE_DRIVER", "STORAGE_DRIVER"},
			Destination: &cfg.Storage.Driver,
		},
		&cli.StringFlag{
			Name:        "storage.access_id",
			Usage:       "Access key ID for storage authentication",
			EnvVars:     []string{"PLUGIN_ACCESS_KEY_ID", "INPUT_ACCESS_KEY_ID", "ACCESS_KEY_ID"},
			Destination: &cfg.Storage.AccessID,
		},
		&cli.StringFlag{
			Name:        "storage.secret_key",
			Usage:       "Secret access key for storage authentication",
			EnvVars:     []string{"PLUGIN_SECRET_ACCESS_KEY", "INPUT_SECRET_ACCESS_KEY", "SECRET_ACCESS_KEY"},
			Destination: &cfg.Storage.SecretKey,
		},
		&cli.StringFlag{
			Name:        "storage.endpoint",
			Value:       "s3.amazonaws.com",
			Usage:       "Storage service endpoint URL",
			EnvVars:     []string{"PLUGIN_STORAGE_ENDPOINT", "INPUT_STORAGE_ENDPOINT", "STORAGE_ENDPOINT"},
			Destination: &cfg.Storage.Endpoint,
		},
		&cli.StringFlag{
			Name:        "storage.bucket",
			Usage:       "Bucket name to store backup files",
			EnvVars:     []string{"PLUGIN_STORAGE_BUCKET", "INPUT_STORAGE_BUCKET", "STORAGE_BUCKET"},
			Destination: &cfg.Storage.Bucket,
		},
		&cli.StringFlag{
			Name:        "storage.region",
			Value:       "ap-northeast-1",
			Usage:       "Region where your storage bucket is located",
			EnvVars:     []string{"PLUGIN_STORAGE_REGION", "INPUT_STORAGE_REGION", "STORAGE_REGION"},
			Destination: &cfg.Storage.Region,
		},
		&cli.StringFlag{
			Name:        "storage.path",
			Value:       "backup",
			Usage:       "Path/folder within the bucket to store backups",
			EnvVars:     []string{"PLUGIN_STORAGE_PATH", "INPUT_STORAGE_PATH", "STORAGE_PATH"},
			Destination: &cfg.Storage.Path,
		},
		&cli.BoolFlag{
			Name:        "storage.ssl",
			Usage:       "Use SSL for secure storage connection",
			EnvVars:     []string{"PLUGIN_STORAGE_SSL", "INPUT_STORAGE_SSL", "STORAGE_SSL"},
			Destination: &cfg.Storage.SSL,
		},
		&cli.StringFlag{
			Name:        "storage.dump_name",
			Usage:       "Filename for the database dump",
			EnvVars:     []string{"PLUGIN_STORAGE_DUMP_NAME", "INPUT_STORAGE_DUMP_NAME", "STORAGE_DUMP_NAME"},
			Destination: &cfg.Storage.DumpName,
			Value:       "dump.sql.gz",
		},
		&cli.BoolFlag{
			Name:  "storage.insecure_skip_verify",
			Usage: "Skip SSL certificate verification (not recommended for production)",
			EnvVars: []string{
				"PLUGIN_STORAGE_INSECURE_SKIP_VERIFY",
				"INPUT_STORAGE_INSECURE_SKIP_VERIFY",
				"STORAGE_INSECURE_SKIP_VERIFY",
			},
			Destination: &cfg.Storage.SkipVerify,
		},
		&cli.IntFlag{
			Name:        "storage.days",
			Usage:       "Number of days to keep backup files (0 = keep forever)",
			EnvVars:     []string{"PLUGIN_STORAGE_DAYS", "INPUT_STORAGE_DAYS", "STORAGE_DAYS"},
			Destination: &cfg.Storage.Days,
			Value:       0,
		},

		// SCHEDULE
		&cli.StringFlag{
			Name:        "time.schedule",
			Usage:       "Backup schedule in cron format (e.g., @daily, @hourly, or '0 0 * * *')",
			EnvVars:     []string{"PLUGIN_TIME_SCHEDULE", "INPUT_TIME_SCHEDULE", "TIME_SCHEDULE"},
			Destination: &cfg.Server.Schedule,
		},
		&cli.StringFlag{
			Name:        "time.location",
			Usage:       "Timezone for scheduling (e.g., Asia/Tokyo, America/New_York)",
			EnvVars:     []string{"PLUGIN_TIME_LOCATION", "INPUT_TIME_LOCATION", "TIME_LOCATION"},
			Destination: &cfg.Server.Location,
		},

		// File Format
		&cli.StringFlag{
			Name:        "file.prefix",
			Usage:       "Text to add before the backup filename",
			EnvVars:     []string{"PLUGIN_FILE_PREFIX", "INPUT_FILE_PREFIX", "FILE_PREFIX"},
			Destination: &cfg.File.Prefix,
		},
		&cli.StringFlag{
			Name:        "file.suffix",
			Usage:       "Text to add after the backup filename",
			EnvVars:     []string{"PLUGIN_FILE_SUFFIX", "INPUT_FILE_SUFFIX", "FILE_SUFFIX"},
			Destination: &cfg.File.Suffix,
		},
		&cli.StringFlag{
			Name:        "file.format",
			Usage:       "Date/time format for backup filenames (Go time format)",
			Value:       "20060102150405",
			EnvVars:     []string{"PLUGIN_FILE_FORMAT", "INPUT_FILE_FORMAT", "FILE_FORMAT"},
			Destination: &cfg.File.Format,
		},
		&cli.StringFlag{
			Name:        "webhook.url",
			Usage:       "URL to notify when backup completes",
			EnvVars:     []string{"PLUGIN_WEBHOOK_URL", "INPUT_WEBHOOK_URL", "WEBHOOK_URL"},
			Destination: &cfg.Webhook.URL,
		},
		&cli.BoolFlag{
			Name:        "webhook.insecure",
			Usage:       "Allow insecure (HTTP) webhook connections",
			EnvVars:     []string{"PLUGIN_WEBHOOK_INSECURE", "INPUT_WEBHOOK_INSECURE", "WEBHOOK_INSECURE"},
			Destination: &cfg.Webhook.Insecure,
		},
	}
}
