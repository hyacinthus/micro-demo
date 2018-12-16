package main

import (
	"encoding/json"
	"net/http"

	"github.com/hyacinthus/x/xerr"
	"github.com/labstack/echo"
	nsq "github.com/nsqio/go-nsq"
)

// ============= 模型部分 ==============

// Entity 实体样例
type Entity struct {
	ID
	// 标题
	Title string `json:"title"`
	Time
}

// EntityUpdate 更新请求结构体，用指针可以判断是否有请求这个字段
type EntityUpdate struct {
	// 标题
	Title *string `json:"title"`
}

// ============= 业务部分 ==============

func findEntityByID(id string) (*Entity, error) {
	var r = new(Entity)
	if err := db.Where("id = ?", id).First(r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

// ============= 事件处理部分 ==============

// ReceiveEntity 处理收到的对象
// Topic: entity_new Channel: ske
// Body: Entity
func ReceiveEntity(msg *nsq.Message) error {
	entity := new(Entity)
	err := json.Unmarshal(msg.Body, entity)
	if err != nil {
		log.WithError(err).Errorf("接收到的消息格式错误: %+v", msg)
		return err
	}
	log.Info(entity)
	return nil
}

// ============= REST 部分 ==============

// createEntity 新建实体
// @Tags 实体
// @Summary 新建实体
// @Description 新建一条实体
// @Accept  json
// @Produce  json
// @Param data body main.Entity true "实体内容"
// @Success 201 {object} main.Entity
// @Failure 400 {object} xerr.Error
// @Failure 401 {object} xerr.Error
// @Failure 500 {object} xerr.Error
// @Security ApiKeyAuth
// @Router /entities [post]
func createEntity(c echo.Context) error {
	// 输入
	var r = new(Entity)
	if err := c.Bind(r); err != nil {
		return err
	}
	// 校验
	if r.Title == "" {
		return xerr.New(400, "BadRequest", "Empty title")
	}
	// 保存
	if err := db.Create(r).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, r)
}

// updateEntity 更新实体
// @Tags 实体
// @Summary 更新实体
// @Description 更新指定id的实体
// @Accept  json
// @Produce  json
// @Param data body main.EntityUpdate true "更新内容"
// @Success 200 {object} main.Entity
// @Failure 400 {object} xerr.Error
// @Failure 401 {object} xerr.Error
// @Failure 404 {object} xerr.Error
// @Failure 500 {object} xerr.Error
// @Security ApiKeyAuth
// @Router /entities/{id} [put]
func updateEntity(c echo.Context) error {
	// 获取URL中的ID
	id := c.Param("id")
	var n = new(EntityUpdate)
	if err := c.Bind(n); err != nil {
		return err
	}
	old, err := findEntityByID(id)
	if err != nil {
		return err
	}
	// 利用指针检查是否有请求这个字段
	if n.Title != nil {
		if *n.Title == "" {
			return xerr.New(400, "BadRequest", "Empty title")
		}
		old.Title = *n.Title
	}

	if err := db.Save(old).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, old)
}

// deleteEntity 删除实体
// @Tags 实体
// @Summary 删除实体
// @Description 删除指定id的实体
// @Accept  json
// @Produce  json
// @Param id path int true "实体编号"
// @Success 204
// @Failure 400 {object} xerr.Error
// @Failure 401 {object} xerr.Error
// @Failure 404 {object} xerr.Error
// @Failure 500 {object} xerr.Error
// @Security ApiKeyAuth
// @Router /entities/{id} [delete]
func deleteEntity(c echo.Context) error {
	id := c.Param("id")
	// 删除数据库对象
	if err := db.Where("id = ?", id).Delete(&Entity{}).Error; err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// getEntity 获取实体
// @Tags 实体
// @Summary 获取实体
// @Description 获取指定id的实体
// @Accept  json
// @Produce  json
// @Param id path int true "实体编号"
// @Success 200 {object} main.Entity
// @Failure 400 {object} xerr.Error
// @Failure 401 {object} xerr.Error
// @Failure 500 {object} xerr.Error
// @Security ApiKeyAuth
// @Router /entities/{id} [get]
func getEntity(c echo.Context) error {
	id := c.Param("id")
	r, err := findEntityByID(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

// getEntitys 获取实体列表
// @Tags 实体
// @Summary 获取实体列表
// @Description 获取用户的全部实体，有分页，默认一页10条。
// @Accept  json
// @Produce  json
// @Param page query int false "页码"
// @Param per_page query int false "每页几条"
// @Success 200 {array} main.Entity
// @Failure 400 {object} xerr.Error
// @Failure 401 {object} xerr.Error
// @Failure 500 {object} xerr.Error
// @Security ApiKeyAuth
// @Router /entities [get]
func getEntitys(c echo.Context) error {
	// 提前make可以让查询没有结果的时候返回空列表
	var ns = make([]*Entity, 0)
	// 分页信息
	limit := c.Get("limit").(int)
	offset := c.Get("offset").(int)
	err := db.Order("updated_at desc").
		Offset(offset).Limit(limit).Find(&ns).Error
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ns)
}
