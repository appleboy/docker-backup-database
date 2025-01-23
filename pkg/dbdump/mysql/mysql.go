package mysql

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
	port := "3306"
	if len(data) > 1 {
		port = data[1]
	}

	return host, port
}

// Exec for dump command
func (d Dump) Exec(ctx context.Context) error {
	envs := os.Environ()

	// Print the version number for the command line tools
	cmd := exec.CommandContext(ctx, "mysqldump", "--version")
	cmd.Env = envs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to get mysqldump version: %w", err)
	}

	flags := []string{}
	host, port := getHostPort(d.Host)
	if host != "" {
		flags = append(flags, "-h", host)
	}
	if port != "" {
		flags = append(flags, "-P", port)
	}

	if d.Username != "" {
		flags = append(flags, "-u", d.Username)
	}

	if d.Opts != "" {
		flags = append(flags, d.Opts)
	}

	if d.Name != "" {
		flags = append(flags, d.Name)
	}

	if d.Password != "" {
		envs = append(envs, "MYSQL_PWD="+d.Password)
	}

	cmd = exec.CommandContext(ctx, "mysqldump", flags...)
	cmd.Env = envs

	// Use a pipe to gzip the output
	gzipCmd := exec.CommandContext(ctx, "gzip")
	gzipCmd.Stdin, _ = cmd.StdoutPipe()
	gzipCmd.Stdout = os.Stdout
	gzipCmd.Stderr = os.Stderr

	trace(cmd)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start mysqldump: %w", err)
	}

	trace(gzipCmd)
	if err := gzipCmd.Start(); err != nil {
		return fmt.Errorf("failed to start gzip: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("mysqldump failed: %w", err)
	}

	if err := gzipCmd.Wait(); err != nil {
		return fmt.Errorf("gzip failed: %w", err)
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
