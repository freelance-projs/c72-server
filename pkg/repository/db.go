package repository

import (
	"database/sql"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DB(sqlConn *sql.DB) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlConn,
	}), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			Colorful:             true,
			ParameterizedQueries: false,
			LogLevel:             logger.Info,
		}),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
