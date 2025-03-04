package config

type (
	// Config provides the system configuration.
	Config struct {
		Server   Server
		Database Database
		Storage  Storage
		File     File
		Webhook  Webhook
	}

	Webhook struct {
		URL string
	}

	// Storage config
	Storage struct {
		Endpoint   string
		AccessID   string
		SecretKey  string
		SSL        bool
		Region     string
		Bucket     string
		Path       string
		Driver     string
		DumpName   string
		SkipVerify bool
		Days       int
	}

	// Server provides the server configuration.
	Server struct {
		Addr     string
		Host     string
		Proto    string
		Port     string
		Pprof    bool
		Root     string
		Debug    bool
		Schedule string
		Location string
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

	// File struct
	File struct {
		Prefix string
		Format string
		Suffix string
	}
)
