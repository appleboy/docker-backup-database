package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
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
	"github.com/appleboy/go-storage/core"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
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
		slog.Error("can't run app", "err", err.Error())
	}
}

func run(cfg *config.Config) cli.ActionFunc {
	return func(ctx *cli.Context) error {
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

		// get context
		appCtx := ctx.Context

		// check bucket exist
		if exist, err := s3.BucketExists(appCtx, cfg.Storage.Bucket); !exist {
			if err != nil {
				return errors.New("bucket not exist or you don't have permission: " + err.Error())
			}

			// create new bucket
			if err := s3.CreateBucket(appCtx, cfg.Storage.Bucket, cfg.Storage.Region); err != nil {
				return errors.New("can't create bucket: " + err.Error())
			}
		}

		// Set lifecycle on bucket or an object prefix.
		if cfg.Storage.Days > 0 && cfg.Storage.Path != "" {
			if err := s3.SetLifeCycle(appCtx, cfg.Storage.Bucket, &core.LifecycleConfig{
				Days:   cfg.Storage.Days,
				Prefix: cfg.Storage.Path,
			}); err != nil {
				return errors.New("can't set bucket lifecycle: " + err.Error())
			}
			slog.Info("set bucket lifecycle successfully",
				"days", cfg.Storage.Days,
				"prefix", cfg.Storage.Path,
				"bucket", cfg.Storage.Bucket,
			)
		}

		if cfg.Server.Schedule == "" {
			slog.Info("no schedule found, backup database now")
			return backupDB(appCtx, cfg, s3)
		}

		// start cron job
		c := cron.New()
		if cfg.Server.Location != "" {
			if loc, err := time.LoadLocation(cfg.Server.Location); err == nil {
				c = cron.New(cron.WithLocation(loc))
			} else {
				return err
			}
		}

		// backup task
		backupTask := func() {
			slog.Info("start backup database now", "schedule", cfg.Server.Schedule)
			if err := backupDB(appCtx, cfg, s3); err != nil {
				slog.Error("can't backup database", "err", err.Error())
				return
			}
			slog.Info("backup database successfully")

			// call webhook if configured
			if cfg.Webhook.URL != "" {
				if err := callWebhook(appCtx, cfg.Webhook.URL, cfg.Webhook.Insecure); err != nil {
					slog.Error("can't call webhook", "err", err.Error())
					return
				}
				slog.Info("call webhook successfully")
			}
		}

		if _, err := c.AddFunc(cfg.Server.Schedule, backupTask); err != nil {
			return fmt.Errorf("crontab schedule error: %w", err)
		}
		c.Start()

		// Register shutdown signal notifications
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		slog.Info("shutting down backup service")

		// stop cron job
		c.Stop()

		return nil
	}
}

func backupDB(ctx context.Context, cfg *config.Config, s3 core.Storage) error {
	// Initialize database dump interface
	backup := dbdump.NewEngine(*cfg)
	if err := backup.Exec(ctx); err != nil {
		return err
	}

	// Read the dump file
	content, err := os.ReadFile(cfg.Storage.DumpName)
	if err != nil {
		return errors.New("can't open the gzip file: " + err.Error())
	}

	// Construct the filename
	filenameParts := []string{}
	if cfg.File.Prefix == "" {
		cfg.File.Prefix = cfg.Database.Driver
	}
	filenameParts = append(filenameParts, cfg.File.Prefix)

	timeFormat := time.Now().Format(cfg.File.Format)
	if cfg.Server.Location != "" {
		loc, _ := time.LoadLocation(cfg.Server.Location)
		timeFormat = time.Now().In(loc).Format("20060102150405")
	}
	filenameParts = append(filenameParts, timeFormat)

	if cfg.File.Suffix != "" {
		filenameParts = append(filenameParts, cfg.File.Suffix)
	}

	filePath := path.Join(cfg.Storage.Path, strings.Join(filenameParts, "-")+".sql.gz")

	// Upload the file to S3
	return s3.UploadFile(ctx, cfg.Storage.Bucket, filePath, content, nil)
}

// callWebhook sends a POST request to the specified webhook URL.
// It accepts a context for request cancellation, the target URL, and a boolean flag to disable SSL certificate verification.
//
// Parameters:
//   - ctx: context.Context - The context to control the request lifetime.
//   - target: string - The webhook URL to send the POST request to.
//   - insecure: bool - If true, disables SSL certificate verification.
//
// Returns:
//   - error: An error if the request fails or the response status code is not 200 OK.
func callWebhook(ctx context.Context, target string, insecure bool) error {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return fmt.Errorf("invalid webhook URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, parsedURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set a useful User-Agent header
	req.Header.Set("User-Agent", "Docker-Backup-Database/"+Version)

	transport := http.DefaultTransport.(*http.Transport).Clone()
	if insecure {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // #nosec
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-success status code: %d", resp.StatusCode)
	}

	return nil
}
