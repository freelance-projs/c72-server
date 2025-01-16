package main

import (
	"database/sql"
	"fmt"
	"strings"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
)

func main() {
}

func migrateUp() error {
	dsn := "root:secret@tcp(host.docker.internal:3306)/tag_scan?charset=utf8mb4&loc=Local&parseTime=True"
	dsn = strings.Split(dsn, "?")[0]
	// fmt.Println("dsn", dsn, strings.Split(dsn, "?")[0])
	_, err := mysqlDriver.ParseDSN(strings.Split(dsn, "?")[0])
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s?multiStatements=true", dsn))
	if err != nil {
		return err
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}
