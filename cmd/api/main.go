package main

import (
	"net/http"

	"github.com/jinzhu/configor"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/hyacinthus/x/auth"
	"github.com/hyacinthus/x/page"
	"github.com/hyacinthus/x/xcc"
	"github.com/hyacinthus/x/xdb"
	"github.com/hyacinthus/x/xerr"
	"github.com/hyacinthus/x/xkv"
	"github.com/hyacinthus/x/xlog"
	"github.com/hyacinthus/x/xmq"
	"github.com/hyacinthus/x/xobj"

	"github.com/hyacinthus/micro-demo/demo"
	"github.com/hyacinthus/micro-demo/server"
)

// @title 雪豹资料库 RESTful API
// @version 1.0
// @description 雪豹知识库后台接口

// @contact.name Muninn
// @contact.email hyacinthus@gmail.com

// @schemes https
// @BasePath /
//
// Securities
// @securitydefinitions.apikey Auth.Bearer
// @in header
// @name Authorization
func main() {
	// ============ Config ============
	var config = new(Config)
	godotenv.Load()
	configor.Load(config)
	if config.App.Debug {
		xlog.Debug()
	}

	// ============ Iaas ============
	// gorm mysql db connection
	var db = xdb.New(config.DB)
	defer db.Close()
	if config.App.Debug {
		db.LogMode(true)
	}
	// redis client
	var kv = xkv.New(config.Redis)
	defer kv.Close()
	// cache client
	var cc = xcc.New(kv)
	// init cos image client
	var img = xobj.New(xobj.ProviderCOS, "image", config.Object)
	// init cos static client
	var obj = xobj.New(xobj.ProviderCOS, "static", config.Object)
	// init mq
	var mq = xmq.New(xmq.ProviderNSQ, config.MQ)
	defer mq.Close()

	// async create tables
	go db.AutoMigrate(&demo.Park{})

	// ============ Service ============
	var s = server.NewService(db, cc, mq, img, obj)
	var h = server.NewHandler(s)

	// ========== Event Owner ===========
	// mq.CreateTopic(mp.TopicUserSubs)

	// ========= Event Worker ===========
	// mq.Sub(mp.TopicTempMsgSent, "wechat", mpw.TempMsgLogProcessor)

	// ============ Echo ============
	e := echo.New()

	if config.App.Debug {
		e.Debug = true
	}

	// error handler
	e.HTTPErrorHandler = xerr.ErrorHandler

	// common middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(page.Middleware(config.App.PageSize)) // 分页参数解析，在 pagination.go 定义

	// sys group
	// sys := e.Group("/sys")

	// admin group
	admin := e.Group("/demo/admin")
	admin.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:     []byte(config.App.JWTSecret),
		SuccessHandler: auth.ParseJWT,
	}))

	// auth group , debug 时总是当作 debug 用户
	user := e.Group("/demo")
	user.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:     []byte(config.App.JWTSecret),
		SuccessHandler: auth.ParseJWT,
	}))

	e.GET("/status", getStatus)

	// mp users
	// admin.POST("/mp/user/sync", mph.PostUserSync)

	// park
	user.GET("/parks/:id", h.GetPark)
	user.GET("/parks", h.GetParks)

	// park admin
	admin.GET("/parks/:id", h.AdminGetPark)
	admin.GET("/parks", h.AdminGetParks)
	admin.POST("/parks", h.AdminPostPark)
	admin.PUT("/parks/:id", h.AdminPutPark)
	admin.DELETE("/parks/:id", h.AdminDeletePark)

	// Start echo server
	e.Logger.Fatal(e.Start(config.App.Host + ":" + config.App.Port))
}

// API状态 成功204 失败500
func getStatus(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
