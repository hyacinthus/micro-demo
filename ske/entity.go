package ske

import (
	"github.com/hyacinthus/x/model"
	"github.com/levigross/grequests"
)

// ============== 模型 ==============

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

// ============== 方法 ==============

// SomeMethod 实体附带的方法
func (e *Entity) SomeMethod() {
}

// ============== SDK ==============

var entityURL = "http://ske/"

// SetEntityURL 供使用者修改
func SetEntityURL(url string) {
	entityURL = url
}

// GetEntity 拉取指定实体
func GetEntity(id string) (*Entity, error) {
	var entity = new(Entity)
	r, err := grequests.Get(entityURL+"entities/"+id, nil)
	if err != nil {
		return nil, err
	}
	err = r.JSON(entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
