package dbdump

import (
	"errors"

	"backup/pkg/config"
	"backup/pkg/dbdump/mongo"
	"backup/pkg/dbdump/mysql"
	"backup/pkg/dbdump/postgres"
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
