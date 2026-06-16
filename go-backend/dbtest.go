package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

func main() {
	cfg, err := pgconn.ParseConfig("postgres://admin:password123@127.0.0.1:5432/ainyx?sslmode=disable")
	if err != nil {
		fmt.Println("PARSE ERROR:", err)
		return
	}
	fmt.Printf("User: %q\n", cfg.User)
	fmt.Printf("Password: %q\n", cfg.Password)
	fmt.Printf("Host: %q\n", cfg.Host)
	fmt.Printf("Port: %d\n", cfg.Port)
	fmt.Printf("Database: %q\n", cfg.Database)

	conn, err := pgconn.ConnectConfig(context.Background(), cfg)
	if err != nil {
		fmt.Println("CONNECT ERROR:", err)
		return
	}
	defer conn.Close(context.Background())
	fmt.Println("SUCCESS!")
}
