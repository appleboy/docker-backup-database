package dbdump

import (
	"errors"

	"github.com/appleboy/docker-backup-database/pkg/config"
	"github.com/appleboy/docker-backup-database/pkg/dbdump/mongo"
	"github.com/appleboy/docker-backup-database/pkg/dbdump/mysql"
	"github.com/appleboy/docker-backup-database/pkg/dbdump/postgres"
)

// Backup database interface
type Backup interface {
	// Exec backup database
	Exec() error
}

// NewEngine return storage interface
func NewEngine(cfg config.Config) (backup Backup, err error) {
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
	}

	return nil, errors.New("We don't support Databaser Dirver: " + cfg.Database.Driver)
}
