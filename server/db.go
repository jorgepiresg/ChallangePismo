package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func (s *server) createSqlConn() *sqlx.DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := s.config.DB

	db, err := sqlx.ConnectContext(ctx, cfg.DriverName, fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password))
	if err != nil {
		log.Fatal("createSqlConn connection: ", err.Error())
	}

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal("createSqlConn ping: ", err.Error())
	}

	db.SetMaxIdleConns(30)
	db.SetMaxOpenConns(300)
	db.SetConnMaxLifetime(30 * time.Minute)

	s.runMigrationsUp(db)

	log.Println("database started")
	return db
}

func (s *server) runMigrationsUp(db *sqlx.DB) {

	migrationFile := s.config.DB.MigrationFile

	entries, err := os.ReadDir(migrationFile)
	if err != nil {
		log.Fatal(err)
	}

	tx := db.MustBegin()

	for _, e := range entries {

		if !isUp(e.Name()) {
			continue
		}

		bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", migrationFile, e.Name()))
		if err != nil {
			continue
		}

		query := string(bytes)

		tx.MustExec(query)
	}

	tx.Commit()

}

func isUp(fileName string) bool {
	return strings.Contains(fileName, "up")
}

func isDown(fileName string) bool {
	return strings.Contains(fileName, "down")
}
