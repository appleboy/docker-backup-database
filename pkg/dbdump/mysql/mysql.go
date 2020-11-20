package mysql

import (
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
func (d Dump) Exec() error {

	// Print the version number fo rht ecommand line tools
	cmd := exec.Command("mysqldump", "--version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)
	if err := cmd.Run(); err != nil {
		return err
	}

	flags := []string{"mysqldump"}
	if d.Name != "" {
		flags = append(flags, "-d", d.Name)
	}

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

	// add gzip command
	flags = append(flags, "|", "gzip", ">", "dump.sql.gz")

	envs := os.Environ()
	if d.Password != "" {
		// See the MySQL Environment Variables
		// ref: https://dev.mysql.com/doc/refman/8.0/en/environment-variables.html
		envs = append(envs, fmt.Sprintf("MYSQL_PWD=%s", d.Password))
	}

	cmd = exec.Command("bash", "-c", strings.Join(flags, " "))
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
func NewEngine(host, username, password, name, opts string) (*Dump, error) {
	return &Dump{
		Host:     host,
		Username: username,
		Password: password,
		Name:     name,
		Opts:     opts,
	}, nil
}
