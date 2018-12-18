package main

import (
	"net/http"

	"github.com/go-redis/redis"
	"github.com/hyacinthus/x/auth"
	"github.com/hyacinthus/x/cc"
	"github.com/hyacinthus/x/object"
	"github.com/hyacinthus/x/page"
	"github.com/hyacinthus/x/xerr"
	"github.com/hyacinthus/x/xlog"
	"github.com/hyacinthus/x/xnsq"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	nsq "github.com/nsqio/go-nsq"
)

// 定义全区变量 为了保证执行顺序 初始化均在main中执行
var (
	// config
	config = new(Config)
	// Logger
	log = xlog.Get()
	// gorm mysql db connection
	db *gorm.DB
	// redis client
	rdb *redis.Client
	// nsq producer
	producer *nsq.Producer
	// cos client
	img *object.Client
)

// the only init func
func init() {
	// config
	godotenv.Load()
	configor.Load(config)

	// debug
	if config.APP.Debug {
		xlog.Debug()
	}

	// init gorm and redis
	initDB()
	initRedis()

	// init global cache
	cc.Init(rdb)

	// init nsq
	xnsq.Init(config.NSQ.NsqdAddr, config.NSQ.NsqLookupdAddr)
	producer = xnsq.Producer()

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

func main() {
	defer clean()
	// ========== NSQ =============
	xnsq.Reg("entity_new", "ske", ReceiveEntity)

	// ========== Echo ============
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

	// status
	e.GET("/status", getStatus)
	e.GET("/", getStatus)

	// note Routes
	e.GET("/entities", getEntitys)
	e.POST("/entities", createEntity)
	e.GET("/entities/:id", getEntity)
	e.PUT("/entities/:id", updateEntity)
	e.DELETE("/entities/:id", deleteEntity)

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
