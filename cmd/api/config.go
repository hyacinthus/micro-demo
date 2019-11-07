package main

import (
	"github.com/hyacinthus/x/xdb"
	"github.com/hyacinthus/x/xkv"
	"github.com/hyacinthus/x/xmq"
	"github.com/hyacinthus/x/xobj"
)

// Config is the global settings of this project
type Config struct {
	App struct {
		Debug     bool   `default:"false"`
		Host      string `default:"0.0.0.0"`
		Port      string `default:"80"`
		PageSize  int    `default:"20"`
		JWTSecret string `default:"secret"`
		BaseURL   string `default:"https://api.xuebaox.com/"`
	}

	// 数据库
	DB xdb.Config

	// Redis
	Redis xkv.Config

	// 对象存储
	Object xobj.Config

	// 队列
	MQ xmq.Config
}
