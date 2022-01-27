package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDb() *gorm.DB {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			Colorful:      true,
		},
	)

	config := &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 gormLogger,
	}

	err := godotenv.Load(".env")
	if err != nil {
		panic(".env is not loaded properly")
	}

	var DB_USERNAME = os.Getenv(`DB_USERNAME`)
	var DB_PASSWORD = os.Getenv(`DB_PASSWORD`)
	var DB_NAME = os.Getenv(`DB_NAME`)
	var DB_HOST = os.Getenv(`DB_HOST`)
	var DB_PORT = os.Getenv(`DB_PORT`)

	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local`,
		DB_USERNAME,
		DB_PASSWORD,
		DB_HOST,
		DB_PORT,
		DB_NAME,
	)

	openConnection, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("error connecting to database, error=", err)
		return nil
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: openConnection,
	}), config)

	if err != nil {
		fmt.Println("error connecting to database, error=", err)
		return nil
	}

	return db
}
