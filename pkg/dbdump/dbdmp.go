package dbdump

import (
	"context"
	"errors"

	"github.com/appleboy/docker-backup-database/pkg/config"
	"github.com/appleboy/docker-backup-database/pkg/dbdump/mongo"
	"github.com/appleboy/docker-backup-database/pkg/dbdump/mysql"
	"github.com/appleboy/docker-backup-database/pkg/dbdump/postgres"
)

// Backup database interface
type Backup interface {
	// Exec backup database
	Exec(context.Context) error
}

// NewEngine return storage interface
func NewEngine(cfg config.Config) Backup {
	switch cfg.Database.Driver {
	case "postgres":
		return postgres.NewEngine(
			cfg.Database.Host,
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Storage.DumpName,
			cfg.Database.Opts,
		)
	case "mysql":
		return mysql.NewEngine(
			cfg.Database.Host,
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Storage.DumpName,
			cfg.Database.Opts,
		)
	case "mongo":
		return mongo.NewEngine(
			cfg.Database.Host,
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Storage.DumpName,
			cfg.Database.Opts,
		)
	default:
		panic(errors.New("unsupported database driver"))
	}
}
