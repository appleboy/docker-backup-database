package main

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/appleboy/docker-backup-database/pkg/config"
	"github.com/appleboy/docker-backup-database/pkg/dbdump"

	"github.com/appleboy/go-storage"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
)

// Version set at compile-time
var (
	Version string
)

func main() {
	// Load env-file if it exists first
	if filename, found := os.LookupEnv("PLUGIN_ENV_FILE"); found {
		_ = godotenv.Load(filename)
	}

	if _, err := os.Stat("/run/drone/env"); err == nil {
		_ = godotenv.Overload("/run/drone/env")
	}

	cfg := &config.Config{}
	app := &cli.App{
		Name:      "docker-backup-datavase",
		Usage:     "Docker image to periodically backup your database",
		Copyright: "Copyright (c) " + strconv.Itoa(time.Now().Year()) + " Bo-Yi Wu",
		Version:   Version,
		Flags:     settingsFlags(cfg),
		Action:    run(cfg),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(cfg *config.Config) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if cfg.Server.Schedule != "" {
			c := cron.New()
			if cfg.Server.Location != "" {
				loc, err := time.LoadLocation(cfg.Server.Location)
				if err != nil {
					log.Fatal("crontab location error: " + err.Error())
				}
				c = cron.New(cron.WithLocation(loc))
			}

			if _, err := c.AddFunc(cfg.Server.Schedule, func() {
				log.Println("start backup database now")
				if err := backupDB(ctx.Context, cfg); err != nil {
					log.Fatal("can't backup database: " + err.Error())
				}
				log.Println("backup database successfully")
			}); err != nil {
				log.Fatal("crontab Schedule error: " + err.Error())
			}
			c.Start()
			// Register shutdown signal notifications
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			<-sig
			log.Println("shutting down backup service")
		}
		return backupDB(ctx.Context, cfg)
	}
}

func backupDB(ctx context.Context, cfg *config.Config) error {
	// initial storage interface
	s3, err := storage.NewEngine(storage.Config{
		Endpoint:  cfg.Storage.Endpoint,
		AccessID:  cfg.Storage.AccessID,
		SecretKey: cfg.Storage.SecretKey,
		SSL:       cfg.Storage.SSL,
		Region:    cfg.Storage.Region,
		Path:      cfg.Storage.Path,
		Bucket:    cfg.Storage.Bucket,
		Addr:      cfg.Server.Addr,
		Driver:    cfg.Storage.Driver,
	})
	if err != nil {
		return err
	}

	// initial database dump interface
	backup, err := dbdump.NewEngine(*cfg)
	if err != nil {
		return err
	}

	// check bucket exist
	if exist, err := s3.BucketExists(ctx, cfg.Storage.Bucket); !exist {
		if err != nil {
			return errors.New("bucket not exist or you don't have permission: " + err.Error())
		}

		// create new bucket
		if err := s3.CreateBucket(ctx, cfg.Storage.Bucket, cfg.Storage.Region); err != nil {
			return errors.New("can't create bucket: " + err.Error())
		}
	}

	if err := backup.Exec(ctx); err != nil {
		return err
	}

	// upload file to s3
	content, err := ioutil.ReadFile(cfg.Storage.DumpName)
	if err != nil {
		return errors.New("can't open the gzip file: " + err.Error())
	}

	filename := []string{}
	if cfg.File.Prefix == "" {
		cfg.File.Prefix = cfg.Database.Driver
	}

	filename = append(filename, cfg.File.Prefix)

	defaultFormat := time.Now().Format(cfg.File.Format)
	if cfg.Server.Location != "" {
		loc, _ := time.LoadLocation(cfg.Server.Location)
		defaultFormat = time.Now().In(loc).Format("20060102150405")
	}

	filename = append(filename, defaultFormat)
	if cfg.File.Suffix != "" {
		filename = append(filename, cfg.File.Suffix)
	}

	filePath := path.Join(cfg.Storage.Path, strings.Join(filename, "-")+".sql.gz")

	// backup database
	return s3.UploadFile(ctx, cfg.Storage.Bucket, filePath, content, nil)
}
