/**
 * connect to database
 *
**/

// New database.go

package database

import (
	"time"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/openlog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// type global
var trueInstance gorm.DB
var instance *gorm.DB

// Connects to database
func Connect() error {
	host := archaius.GetString("database.host", "localhost")
	user := archaius.GetString("database.user", "postgres")
	dbname := archaius.GetString("database.dbname", "postgres")
	password := archaius.GetString("database.password", "root")
	port := archaius.GetString("database.port", "5432")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			utc, _ := time.LoadLocation("")
			return time.Now().In(utc)
		},
		Logger: logger.Default.LogMode(logger.Info),
	}) // <-- thread safe

	if err != nil {
		openlog.Error("error occured while connecting database")
		return err
	}

	instance = db
	trueInstance = *db
	return nil
}

// Provides the instance of the database
func GetClient() *gorm.DB {
	return instance
}

func StartTransaction() {
	*instance = *instance.Begin()
	if instance.Error != nil {
		openlog.Error("error occured while starting transaction [" + instance.Error.Error())
	}
}

func CommitTransaction() {
	instance = instance.Commit()
	if instance.Error != nil {
		openlog.Error("error occured while starting transaction [" + instance.Error.Error())
	}
	*instance = trueInstance
}

func RollbackTransaction() {
	instance = instance.Rollback()
	if instance.Error != nil {
		openlog.Error("error occured while rolling back transaction [" + instance.Error.Error())
	}
	*instance = trueInstance
}
