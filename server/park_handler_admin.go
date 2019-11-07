package server

import (
	"github.com/hyacinthus/x/page"
	"github.com/hyacinthus/micro-demo/demo"
	"github.com/labstack/echo/v4"

	_ "github.com/hyacinthus/x/xerr" // for swagger
)

// AdminGetPark godoc
// @Summary 获取工业园
// @Description 目前只有发改委批复的国家级和省级
// @ID admin-get-park
// @Tags park-admin
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} demo.Park
// @Failure 404 {object} xerr.Error
// @Router /demo/admin/parks/{id} [get]
// @Security Auth.Bearer
func (h *Handler) AdminGetPark(c echo.Context) error {
	id := c.Param("id")
	resp, err := h.s.FindPark(id)
	if err != nil {
		return err
	}
	return c.JSON(200, resp)
}

// AdminGetParks godoc
// @Summary 获取工业园列表
// @Description
// @ID admin-get-parks
// @Tags park-admin
// @Produce json
// @Param approval_before query string false "YYYYMM"
// @Param approval_after query string false "YYYYMM"
// @Param page query int false "第几页"
// @Param per_page query int false "每页几条"
// @Success 200 {array} demo.Park
// @Failure 400 {object} xerr.Error
// @Router /demo/admin/parks [get]
// @Security Auth.Bearer
func (h *Handler) AdminGetParks(c echo.Context) error {
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
	return c.JSON(200, resp)
}

// AdminPostPark godoc
// @Summary 新建工业园
// @Description
// @ID admin-post-park
// @Tags park-admin
// @Accept  json
// @Produce json
// @Param park body demo.Park true "Add Park"
// @Success 200 {object} demo.Park
// @Failure 400 {object} xerr.Error
// @Router /demo/admin/parks [post]
// @Security Auth.Bearer
func (h *Handler) AdminPostPark(c echo.Context) error {
	var req = new(demo.Park)
	err := c.Bind(req)
	resp, err := h.s.CreatePark(req)
	if err != nil {
		return err
	}
	return c.JSON(200, resp)
}

// AdminPutPark godoc
// @Summary 修改工业园
// @Description
// @ID admin-put-park
// @Tags park-admin
// @Accept  json
// @Produce json
// @Param id path string true "id"
// @Param park body demo.ParkUpdate true "Change Park"
// @Success 200 {object} demo.Park
// @Failure 400 {object} xerr.Error
// @Router /demo/admin/parks/{id} [put]
// @Security Auth.Bearer
func (h *Handler) AdminPutPark(c echo.Context) error {
	id := c.Param("id")
	var req = new(demo.ParkUpdate)
	err := c.Bind(req)
	resp, err := h.s.UpdatePark(id, req)
	if err != nil {
		return err
	}
	return c.JSON(200, resp)
}

// AdminDeletePark godoc
// @Summary 修改工业园
// @Description
// @ID admin-delete-park
// @Tags park-admin
// @Param id path string true "id"
// @Success 204
// @Failure 400 {object} xerr.Error
// @Router /demo/admin/parks/{id} [delete]
// @Security Auth.Bearer
func (h *Handler) AdminDeletePark(c echo.Context) error {
	id := c.Param("id")
	err := h.s.RemovePark(id)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}
