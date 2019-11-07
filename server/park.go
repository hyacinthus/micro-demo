package server

import (
	"github.com/hyacinthus/x/xerr"
	"github.com/hyacinthus/micro-demo/demo"
)

// FindPark 工业园
func (s *Service) FindPark(id string) (*demo.Park, error) {
	var resp = new(demo.Park)
	err := s.db.First(resp, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// FindParks 工业园列表
func (s *Service) FindParks(query *demo.ParkQuery, offset int, limit int) ([]*demo.Park, error) {
	var resp = make([]*demo.Park, 0)
	tx := s.db
	// 条件区
	if query.ApprovalAfter != nil {
		tx = tx.Where("approval_time >= ?", query.ApprovalAfter)
	}
	if query.ApprovalBefore != nil {
		tx = tx.Where("approval_time <= ?", query.ApprovalBefore)
	}
	// 查询
	err := tx.Offset(offset).Limit(limit).Find(&resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreatePark 创建新工业园
func (s *Service) CreatePark(park *demo.Park) (*demo.Park, error) {
	// 校验数据
	if park.ID != "" {
		return nil, xerr.New(400, "InvalidData", "新建记录不接受提供ID")
	}
	// 保存
	err := s.db.Create(park).Error
	if err != nil {
		return nil, err
	}
	return park, nil
}

// UpdatePark 更新工业园
func (s *Service) UpdatePark(id string, data *demo.ParkUpdate) (*demo.Park, error) {
	var old = new(demo.Park)

	// 查询并锁定
	tx := s.db.Begin()
	err := tx.Set("gorm:query_option", "FOR UPDATE").First(old, "id = ?", id).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 逐个核对输入 有必要的话还可以校验
	if data.Code != nil {
		old.Code = *data.Code
	}
	if data.Name != nil {
		old.Name = *data.Name
	}
	if data.ApprovalTime != nil {
		old.ApprovalTime = *data.ApprovalTime
	}

	// 保存
	err = tx.Save(old).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return old, nil
}

// RemovePark 删除工业园
func (s *Service) RemovePark(id string) error {
	return s.db.Delete(&demo.Park{}, "id = ?", id).Error
}
