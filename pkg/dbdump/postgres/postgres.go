package postgres

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
	port := "5432"
	if len(data) > 1 {
		port = data[1]
	}

	return host, port
}

// Exec for dump command
func (d Dump) Exec(ctx context.Context) error {
	envs := os.Environ()

	// Print the version number fo rht ecommand line tools
	cmd := exec.CommandContext(ctx, "pg_dump", "--version")
	cmd.Env = envs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)
	if err := cmd.Run(); err != nil {
		return err
	}

	flags := []string{"pg_dump"}
	host, port := getHostPort(d.Host)
	if host != "" {
		flags = append(flags, "-h", host)
	}
	if port != "" {
		flags = append(flags, "-p", port)
	}

	if d.Username != "" {
		flags = append(flags, "-U", d.Username)
	}

	if d.Opts != "" {
		flags = append(flags, d.Opts)
	}

	if d.Name != "" {
		flags = append(flags, d.Name)
	}

	// add gzip command
	flags = append(flags, "|", "gzip", ">", d.DumpName)

	if d.Password != "" {
		envs = append(envs, fmt.Sprintf("PGPASSWORD=%s", d.Password))
	}

	cmd = exec.CommandContext(ctx, "bash", "-c", strings.Join(flags, " ")) //nolint:gosec
	cmd.Env = envs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)
	return cmd.Run()
}

// trace prints the command to the stdout.
func trace(cmd *exec.Cmd) {
	fmt.Printf("$ %s\n", strings.Join(cmd.Args, " "))
}

// NewEngine struct
func NewEngine(host, username, password, name, dumpName, opts string) (*Dump, error) {
	return &Dump{
		Host:     host,
		Username: username,
		Password: password,
		Name:     name,
		Opts:     opts,
		DumpName: dumpName,
	}, nil
}
