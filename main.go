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
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	obj *object.Client
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
	obj = object.New(&object.Config{
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

	// user group
	// 需要用户登录后才能调用
	user := e.Group("")
	user.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:     []byte(config.APP.JWTSecret),
		SuccessHandler: auth.ParseJWT,
	}))

	// admin group
	// 需要管理员登录后才能调用
	// TODO: 现在在登录环节调用不同的登录接口 将来考虑这里换个parse方法解析管理员详细权限
	admin := e.Group("admin")
	admin.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:     []byte(config.APP.JWTSecret),
		SuccessHandler: auth.ParseJWT,
	}))

	// sys group
	// 微服务间内部调用，不验证，不在中间件暴露到外界
	sys := e.Group("sys")

	// status
	e.GET("/status", getStatus)
	e.GET("/", getStatus)

	// entity Routes
	user.GET("/entities", userGetEntitys)
	user.POST("/entities", userCreateEntity)
	user.GET("/entities/:id", userGetEntity)
	user.PUT("/entities/:id", userUpdateEntity)
	user.DELETE("/entities/:id", userDeleteEntity)
	sys.GET("/entities/:id", sysGetEntity)

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
