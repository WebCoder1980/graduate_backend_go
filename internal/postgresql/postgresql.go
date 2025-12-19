package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"strconv"
)

func getConnStr() string {
	port, err := strconv.Atoi(os.Getenv("postgresql_port"))
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("postgresql_host"),
		port,
		os.Getenv("postgresql_user"),
		os.Getenv("postgresql_password"),
		os.Getenv("postgresql_dbname"),
		"disable",
	)
}

func getConnection() *sql.DB {
	connStr := getConnStr()
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
