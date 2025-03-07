package mongo

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Dump provides dump execution arguments.
type Dump struct {
	Host     string
	Username string
	Password string
	Name     string
	Opts     string
	DumpName string
}

func getHostPort(h string) (string, string) {
	data := strings.Split(h, ":")
	host := data[0]
	port := "27017"
	if len(data) > 1 {
		port = data[1]
	}

	return host, port
}

// Exec for dump command
func (d Dump) Exec(ctx context.Context) error {
	envs := os.Environ()

	// Print the version number for the command line tools
	cmd := exec.CommandContext(ctx, "mongodump", "--version")
	cmd.Env = envs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to get mongodump version: %w", err)
	}

	flags := []string{}
	host, port := getHostPort(d.Host)
	if host != "" {
		flags = append(flags, "-h", host)
	}
	if port != "" {
		flags = append(flags, "--port", port)
	}

	if d.Username != "" {
		flags = append(flags, "-u", d.Username)
	}

	if d.Password != "" {
		envs = append(envs, "MONGO_PWD="+d.Password)
	}

	if d.Name != "" {
		flags = append(flags, "-d", d.Name)
	}

	// Compresses the output. If mongodump outputs to the dump directory,
	// the new feature compresses the individual files. The files have the suffix .gz.
	flags = append(flags, "--gzip")
	flags = append(flags, "--archive="+d.DumpName)

	if d.Opts != "" {
		flags = append(flags, d.Opts)
	}

	cmd = exec.CommandContext(ctx, "mongodump", flags...)
	cmd.Env = envs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start mongodump: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("mongodump failed: %w", err)
	}

	return nil
}

// trace prints the command to the stdout.
func trace(cmd *exec.Cmd) {
	fmt.Printf("$ %s\n", strings.Join(cmd.Args, " "))
}

// NewEngine struct
func NewEngine(host, username, password, name, dumpName, opts string) *Dump {
	return &Dump{
		Host:     host,
		Username: username,
		Password: password,
		Name:     name,
		Opts:     opts,
		DumpName: dumpName,
	}
}
