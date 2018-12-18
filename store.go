package main

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	"github.com/hyacinthus/ske/ske"
)

func initRedis() {
	// redis conn
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	log.Info("Redis connect successful.")
}

func initDB() {
	var err error
	// mysql conn
	for {
		db, err = gorm.Open("mysql", config.DB.User+":"+config.DB.Password+
			"@tcp("+config.DB.Host+":"+config.DB.Port+")/"+config.DB.Name+
			"?charset=utf8mb4&parseTime=True&loc=Local&timeout=90s")
		if err != nil {
			log.WithError(err).Warn("waiting for connect to db")
			time.Sleep(time.Second * 2)
			continue
		}
		db.DB().SetConnMaxLifetime(time.Duration(config.DB.Lifetime) * time.Second)
		log.Info("Mysql connect successful.")
		break
	}

	// gorm debug log
	if config.APP.Debug {
		db.LogMode(true)
	}
}

// createTable gorm auto migrate tables
func createTables() {
	db.AutoMigrate(&ske.Entity{})
}
