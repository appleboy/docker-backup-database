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
		Endpoint  string
		AccessID  string
		SecretKey string
		SSL       bool
		Region    string
		Bucket    string
		Path      string
		Driver    string
	}

	// Server provides the server configuration.
	Server struct {
		Addr  string
		Host  string
		Proto string
		Port  string
		Pprof bool
		Root  string
		Debug bool
	}

	// Database config
	Database struct {
		Driver        string
		Username      string
		Password      string
		Name          string
		Host          string
		Opts          string
		UseSQLite3    bool
		UseMySQL      bool
		UseMSSQL      bool
		UsePostgreSQL bool
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
