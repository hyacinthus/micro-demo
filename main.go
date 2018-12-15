package main

import (
	"net/http"

	"github.com/go-redis/redis"
	"github.com/hyacinthus/x/auth"
	"github.com/hyacinthus/x/cc"
	"github.com/hyacinthus/x/object"
	"github.com/hyacinthus/x/page"
	"github.com/hyacinthus/x/xerr"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

// 定义全区变量 为了保证执行顺序 初始化均在main中执行
var (
	// config
	config = new(Config)
	// Logger
	log *logrus.Logger
	// gorm mysql db connection
	db *gorm.DB
	// redis client
	rdb *redis.Client
	// cos client
	img *object.Client
)

// the only init func
func init() {
	// config
	godotenv.Load()
	configor.Load(config)

	// logger
	log = logrus.New()
	if config.APP.Debug {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "06-01-02 15:04:05.00",
		})
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	// init gorm and redis
	initDB()
	initRedis()

	// init global cache
	cc.Init(rdb)

	// init cos client
	img = object.New(&object.Config{
		AppID:     config.QCloud.AppID,
		Region:    config.QCloud.Region,
		Bucket:    "image", // 注意修改这里
		SecretID:  config.QCloud.SecretID,
		SecretKey: config.QCloud.SecretKey,
	})

	// async create tables
	go createTables()

}

// @title RESTful API DEMO by Golang & Echo
// @version 1.0
// @description This is a demo server.

// @contact.name Muninn
// @contact.email hyacinthus@gmail.com

// @license.name MIT
// @license.url https://github.com/hyacinthus/restdemo/blob/master/LICENSE

// @host demo.crandom.com
// @BasePath /
func main() {
	defer clean()
	// init echo
	e := echo.New()

	// error handler
	e.HTTPErrorHandler = xerr.ErrorHandler

	// common middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(page.Middleware(config.APP.PageSize)) // 分页参数解析，在 pagination.go 定义

	// Echo debug setting
	if config.APP.Debug {
		e.Debug = true
	}

	// auth group
	a := e.Group("/")
	a.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:     []byte(config.APP.JWTSecret),
		SuccessHandler: auth.ParseJWT,
	}))

	// init mysql and redis
	initDB()
	defer db.Close()
	initRedis()
	defer rdb.Close()

	// init global cache
	cc.Init(rdb)

	// async create tables
	go createTables()

	// status
	e.GET("/status", getStatus)

	// note Routes
	e.GET("/notes", getEntitys)
	e.POST("/notes", createEntity)
	e.GET("/notes/:id", getEntity)
	e.PUT("/notes/:id", updateEntity)
	e.DELETE("/notes/:id", deleteEntity)

	// Start echo server
	e.Logger.Fatal(e.Start(config.APP.Host + ":" + config.APP.Port))
}

// API状态 成功204 失败500
func getStatus(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func clean() {
	db.Close()
	rdb.Close()
}
