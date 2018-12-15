package main

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)

// ID 实体共用字段
type ID struct {
	ID string `json:"id" gorm:"type:varchar(20);primary_key"`
}

// Time 实体共用字段
type Time struct {
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	// 最后更新时间
	UpdatedAt time.Time `json:"updated_at"`
	// 软删除
	DeletedAt *time.Time `json:"-"`
}

// BeforeCreate GORM hook
func (id *ID) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", xid.New().String())
	return nil
}

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
	db.AutoMigrate(&Entity{})
}
