package server

import (
	"github.com/hyacinthus/x/auth"
	"github.com/hyacinthus/x/page"
	"github.com/hyacinthus/micro-demo/demo"
	"github.com/labstack/echo/v4"

	_ "github.com/hyacinthus/x/xerr" // for swagger
)

// 普通用户只需要查询功能，对于付费和免费用户使用不同的策略

// GetPark godoc
// @Summary 获取工业园
// @Description 目前只有发改委批复的国家级和省级
// @ID get-park
// @Tags park
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} demo.Park
// @Failure 404 {object} xerr.Error
// @Router /demo/parks/{id} [get]
// @Security Auth.Bearer
func (h *Handler) GetPark(c echo.Context) error {
	id := c.Param("id")
	resp, err := h.s.FindPark(id)
	if err != nil {
		return err
	}

	// 对免费用户隐藏关键信息
	oid := auth.GetOID(c)
	if oid == "" {
		resp.Blur()
	}

	return c.JSON(200, resp)
}

// GetParks godoc
// @Summary 获取工业园列表
// @Description
// @ID get-parks
// @Tags park
// @Produce json
// @Param approval_before query string false "YYYYMM"
// @Param approval_after query string false "YYYYMM"
// @Param page query int false "第几页"
// @Param per_page query int false "每页几条"
// @Success 200 {array} demo.Park
// @Failure 400 {object} xerr.Error
// @Router /demo/parks [get]
// @Security Auth.Bearer
func (h *Handler) GetParks(c echo.Context) error {
	var query = new(demo.ParkQuery)
	offset, limit := page.Parse(c)
	err := c.Bind(query)
	if err != nil {
		return err
	}
	resp, err := h.s.FindParks(query, offset, limit)
	if err != nil {
		return err
	}

	// 对免费用户隐藏关键信息
	oid := auth.GetOID(c)
	if oid == "" {
		for _, item := range resp {
			item.Blur()
		}
	}

	return c.JSON(200, resp)
}
