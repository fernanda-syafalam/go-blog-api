package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

var k = koanf.New(".")

func main() {
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		log.Fatalf("error loading config.yaml: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Missing migrate command. Example: up, down, version, force")
	}
	command := os.Args[1]
	migrationPath := "db/migrations"

	username := k.String("database.username")
	password := k.String("database.password")
	host := k.String("database.host")
	port := k.String("database.port")
	database := k.String("database.name")
	sslmode := k.String("database.sslmode")

	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", username, password, host, port, database, sslmode)
	args := []string{
		"-path", migrationPath,
		"-database", dsn,
		command,
	}

	if len(os.Args) > 2 {
		args = append(args, os.Args[2:]...)
	}
	
	log.Printf("Using Args: %s", args)

	cmd := exec.Command("migrate", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("migrate command failed: %v", err)
	}
}
