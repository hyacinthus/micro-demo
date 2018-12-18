package ske

import (
	"github.com/hyacinthus/x/model"
)

// Entity 实体样例
type Entity struct {
	model.Entity
	// 标题
	Title string `json:"title"`
}

// EntityUpdate 更新请求结构体，用指针可以判断是否有请求这个字段
type EntityUpdate struct {
	// 标题
	Title *string `json:"title"`
}
