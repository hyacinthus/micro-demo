package server

import (
	"github.com/jinzhu/gorm"

	"github.com/hyacinthus/x/xcc"
	"github.com/hyacinthus/x/xlog"
	"github.com/hyacinthus/x/xmq"
	"github.com/hyacinthus/x/xobj"
)

// 日志
var log = xlog.Get()

// Service 业务代码
type Service struct {
	db  *gorm.DB
	cc  xcc.Client
	mq  xmq.Client
	img xobj.Client
	obj xobj.Client
}

// NewService 新服务
func NewService(db *gorm.DB, cc xcc.Client, mq xmq.Client, img, obj xobj.Client) *Service {
	var s = &Service{
		db:  db,
		cc:  cc,
		mq:  mq,
		img: img,
		obj: obj,
	}
	return s
}

// Handler http handler
type Handler struct {
	s *Service
}

// NewHandler get all handler
func NewHandler(s *Service) *Handler {
	return &Handler{
		s: s,
	}
}
