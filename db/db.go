package db

import (
	"github.com/anonyindian/logger"
	"github.com/desterlib/backend-go/config"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var SESSION *gorm.DB

func LoadDB(log *logger.Logger) {
	log = log.Create("DATABASE")
	defer func() {
		log.ChangeLevel(logger.LevelMain)
		log.Println("LOADED")
	}()
	conn, err := pq.ParseURL(config.ValueOf.DatabaseURI)
	if err != nil {
		log.ChangeLevel(logger.LevelError).Printlnf("failed to parse DB URI: %s", err.Error())
	}
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 glogger.Default.LogMode(glogger.Error),
	})
	if err != nil {
		log.Printlnf("failed to connect to DB: %s", err.Error())
	}
	SESSION = db

	dB, _ := db.DB()
	dB.SetMaxOpenConns(100)
	log.ChangeLevel(logger.LevelInfo)
	log.Println("Database connected")

	// Create tables if they don't exist
	SESSION.AutoMigrate(&ConfigDB{})
	log.Println("Auto-migrated database schema")
}
