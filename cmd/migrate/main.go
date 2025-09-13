package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var direction string
	flag.StringVar(&direction, "direction", "up", "migration direction (up or down)")
	flag.Parse()

	// 環境変数からDB接続情報を取得
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		dbPass = "password"
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "english_app"
	}

	// DB接続URL
	dbURL := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// マイグレーションインスタンスの作成
	m, err := migrate.New(
		"file://migrations/mysql",
		dbURL,
	)
	if err != nil {
		log.Fatal("Migration initialization failed:", err)
	}
	defer m.Close()

	// マイグレーションの実行
	if direction == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration up failed:", err)
		}
		log.Println("Migration up completed successfully")
	} else if direction == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration down failed:", err)
		}
		log.Println("Migration down completed successfully")
	} else {
		log.Fatal("Invalid direction. Use 'up' or 'down'")
	}
}
