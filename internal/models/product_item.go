package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tg_go_faka/internal/utils/functions"
	_type "tg_go_faka/internal/utils/type"
)

type ProductItem struct {
	ID          uuid.UUID               `gorm:"primaryKey;not null" json:"id"`
	Status      _type.ProductStatusType `gorm:"default:1;not null" json:"status"`
	CreateTime  int64                   `gorm:"index;autoCreateTime;not null" json:"create_time"`
	EndLockTime int64                   `gorm:"index;not null" json:"end_lock_time"`

	Content string `gorm:"not null" json:"content"`

	ProductID uuid.UUID `json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID"`

	Orders []Order `gorm:"constraint:OnDelete:CASCADE;"` //product_item删除要删除订单

}

func (*ProductItem) TableName() string {
	return "product_item"
}
func (*ProductItem) DefaultOrder() string {
	return "create_time DESC"
}
func (o *ProductItem) ToDict() map[string]interface{} {
	return functions.StructToMap(o)
}

func (o *ProductItem) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
