package config

type (
	// Config provides the system configuration.
	Config struct {
		Server   Server
		Database Database
		Storage  Storage
	}

	// Storage config
	Storage struct {
		Endpoint  string `envconfig:"APP_STORAGE_ENDPOINT" default:"s3-ap-northeast-1.amazonaws.com"`
		AccessID  string `envconfig:"AWS_ACCESS_KEY_ID"`
		SecretKey string `envconfig:"AWS_SECRET_ACCESS_KEY"`
		SSL       bool   `envconfig:"APP_STORAGE_SSL"`
		Region    string `envconfig:"APP_STORAGE_REGION" default:"ap-northeast-1"`
		Bucket    string `envconfig:"APP_STORAGE_BUCKET" default:"test"`
		Path      string `envconfig:"APP_STORAGE_PATH" default:"data"`
		Driver    string `envconfig:"APP_STORAGE_DRIVER" default:"s3"`
	}

	// Server provides the server configuration.
	Server struct {
		Addr  string `envconfig:"APP_SERVER_ADDR"`
		Host  string `envconfig:"APP_SERVER_Host"`
		Proto string `envconfig:"APP_SERVER_Proto" default:"http"`
		Port  string `envconfig:"APP_SERVER_PORT" default:"8080"`
		Pprof bool   `envconfig:"APP_SERVER_PPROF"`
		Root  string `envconfig:"APP_SERVER_ROOT" default:"/"`
		Debug bool   `envconfig:"APP_SERVER_DEBUG"`
	}

	// Database config
	Database struct {
		Driver        string `envconfig:"APP_DATABASE_DRIVER" default:"sqlite3"`
		Username      string `envconfig:"APP_DATABASE_USERNAME" default:"root"`
		Password      string `envconfig:"APP_DATABASE_PASSWORD" default:"root"`
		Name          string `envconfig:"APP_DATABASE_NAME" default:"db"`
		Host          string `envconfig:"APP_DATABASE_HOST" default:"localhost:3306"`
		Schema        string `envconfig:"APP_DATABASE_SCHEMA" default:"public"`
		UseSQLite3    bool   `envconfig:"APP_DATABASE_USE_SQLITE3"`
		UseMySQL      bool   `envconfig:"APP_DATABASE_USE_MYSQL"`
		UseMSSQL      bool   `envconfig:"APP_DATABASE_USE_MSSQL"`
		UsePostgreSQL bool   `envconfig:"APP_DATABASE_USE_POSTGRESQL"`
	}
)

func configureDatabase(c *Config) {
	switch c.Database.Driver {
	case "postgres":
		c.Database.UsePostgreSQL = true
	case "mssql":
		c.Database.UseMSSQL = true
	}
}

func defaultAddress(c *Config) {
	c.Server.Addr = c.Server.Proto + "://" + c.Server.Host
}
