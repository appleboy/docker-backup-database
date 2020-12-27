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
func NewEngine(config config.Config) (backup Backup, err error) {
	switch config.Database.Driver {
	case "postgres":
		return postgres.NewEngine(
			config.Database.Host,
			config.Database.Username,
			config.Database.Password,
			config.Database.Name,
			config.Storage.DumpName,
			config.Database.Opts,
		)
	case "mysql":
		return mysql.NewEngine(
			config.Database.Host,
			config.Database.Username,
			config.Database.Password,
			config.Database.Name,
			config.Storage.DumpName,
			config.Database.Opts,
		)
	case "mongo":
		return mongo.NewEngine(
			config.Database.Host,
			config.Database.Username,
			config.Database.Password,
			config.Database.Name,
			config.Storage.DumpName,
			config.Database.Opts,
		)
	}

	return nil, errors.New("We don't support Databaser Dirver: " + config.Database.Driver)
}
