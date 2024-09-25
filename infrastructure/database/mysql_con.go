package database

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// DB :nodoc:
var DB *gorm.DB

func init() {
	godotenv.Load()
}

// DBInit :nodoc:
func DBInit() (*gorm.DB, error) {
	mysqlCon := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True",
		os.Getenv("DB_MYSQL_USERNAME"),
		os.Getenv("DB_MYSQL_PASSWORD"),
		os.Getenv("DB_MYSQL_HOST"),
		os.Getenv("DB_MYSQL_PORT"),
		os.Getenv("DB_MYSQL_DATABASE"),
	)
	var err error
	DB, err = gorm.Open("mysql", mysqlCon)
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed connected to database %s", mysqlCon))
		return DB, err
	}
	logrus.Info("Successfully connected to database")
	DB.DB().SetConnMaxLifetime(35 * time.Minute)
	DB.DB().SetMaxIdleConns(20)
	DB.DB().SetMaxOpenConns(200)
	DB.LogMode(true)
	DB.SetLogger(log.New(os.Stdout, "\r\n", 0))
	return DB, err
}

// GetConnection :nodoc:
func GetConnection() *gorm.DB {
	if DB == nil {
		logrus.Info("No Active Connection Found")
		DB, _ = DBInit()
	}
	return DB
}
