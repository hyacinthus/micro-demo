package demo

import (
	"time"

	"github.com/hyacinthus/x/model"
)

// Park 工业园/开发区
type Park struct {
	model.Entity
	Code         string    `json:"code" gorm:"type:varchar(20)"`   // 发改委的编码 子园区和二期之类的没有编码
	Name         string    `json:"name"`                           // 核准名称 应该不会重复
	ApprovalTime time.Time `json:"approval_time" gorm:"type:date"` // 批准时间 只精确到月份
}

// ParkUpdate 工业园更新
type ParkUpdate struct {
	Code         *string    `json:"code"`          // 发改委的编码 子园区和二期之类的没有编码
	Name         *string    `json:"name"`          // 核准名称 应该不会重复
	ApprovalTime *time.Time `json:"approval_time"` // 批准时间 只精确到月份
}

// ParkQuery 工业园筛选条件
type ParkQuery struct {
	ApprovalAfter        *time.Time `query:"-"`
	ApprovalAfterString  string     `query:"approval_after,omitempty"` // 批准晚于 YYYYMM
	ApprovalBefore       *time.Time `query:"-"`
	ApprovalBeforeString string     `query:"approval_before,omitempty"` // 批准早于 YYYYMM
}

// Blur 去除敏感信息
func (p *Park) Blur() {
	p.Code = ""
}
