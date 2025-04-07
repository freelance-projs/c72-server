package main

import (
	"bufio"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/common/env"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{}
	cnf     = flag.String("config", "config.yaml", "Path to configuration file.")
)

func main() {
	flag.Parse()
	env.Init(env.WithFile(*cnf))

	migrate := migrateInstance(env.GetString(dto.DBLaundryDSN))

	up := &cobra.Command{
		Use:   "up",
		Short: "Migrate the database to the latest version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := prompt("up"); err != nil {
				fmt.Println(err)
				return
			}

			if err := migrate.Up(); err != nil {
				log.Fatalf("migrate up error: %v", err)
			}
		},
	}

	down := &cobra.Command{
		Use:   "down",
		Short: "Migrate the database to the previous version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := prompt("down"); err != nil {
				return
			}

			if err := migrate.Down(); err != nil {
				log.Fatalf("migrate down error: %v", err)
			}
		},
	}

	force := &cobra.Command{
		Use:   "force",
		Short: "Force the database to a specific version",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			version, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("invalid version: %v", err)
			}
			if err := prompt(fmt.Sprintf("force version %d", version)); err != nil {
				return
			}

			if err := migrate.Force(version); err != nil {
				log.Fatalf("migrate force to version %d failed, error: %v", version, err)
			}
		},
	}

	migrateToVersion := &cobra.Command{
		Use:   "to",
		Short: "Print the current migration version",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			version, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("invalid version: %v", err)
			}
			if err := prompt(fmt.Sprintf("force version %d", version)); err != nil {
				return
			}

			if err := migrate.Migrate(uint(version)); err != nil {
				log.Fatalf("migrate to version fail, error: %v", err)
			}
		},
	}

	version := &cobra.Command{
		Use:   "version",
		Short: "Print the current migration version",
		Run: func(cmd *cobra.Command, args []string) {

			version, dirty, err := migrate.Version()
			if err != nil {
				log.Fatalf("migrate to version fail, error: %v", err)
			}
			fmt.Printf("version: %d | dirty: %t\n", version, dirty)
		},
	}

	rootCmd.AddCommand(up, down, force, migrateToVersion, version)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute error: %v", err)
	}
}

func migrateInstance(dsn string) *migrate.Migrate {
	_, err := mysqldriver.ParseDSN(dsn)
	if err != nil {
		log.Fatalf("parse dsn error: %v", err)
	}
	fmt.Println("dsn: ", dsn)

	db, err := sql.Open("mysql", fmt.Sprintf("%s&multiStatements=true", dsn))
	if err != nil {
		log.Fatalf("open database error: %v", err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("create driver instance error: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatalf("create migrate instance error: %v", err)
	}

	return m
}

func prompt(action string) error {
	reader := bufio.NewReader(os.Stdin)

	cnf, err := mysqldriver.ParseDSN(env.GetString(dto.DBLaundryDSN))
	if err != nil {
		log.Fatalf("parse dsn error: %v", err)
	}

	prompt := fmt.Sprintf("action: %s | addr: %s | db name: %s \n(y/n) ", action, cnf.Addr, cnf.DBName)
	fmt.Print(prompt)

	yesno, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("read input error: %v", err)
	}
	if yesno != "y\n" {
		return errors.New("exit" + yesno)
	}
	fmt.Println("\nmigrating...")

	return nil
}
