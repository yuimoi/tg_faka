package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"tg_go_faka/internal/utils/functions"
	_type "tg_go_faka/internal/utils/type"
)

// 使用指针可以方便的置空，使用原则：必须要判断是否为空的情况
type Product struct {
	ID         int64                   `gorm:"primaryKey;not null" json:"id"` // sqlite3库有个非常奇怪的逻辑，自增是primaryKey的默认自带，但是一旦附加了autoIncrease标签就会导致自增设置失效
	Status     _type.ProductStatusType `gorm:"default:1;not null" json:"status"`
	CreateTime int64                   `gorm:"index;autoCreateTime;not null" json:"create_time"`

	Name  string          `gorm:"not null" json:"name"`
	Desc  string          `gorm:"not null" json:"desc"`
	Price decimal.Decimal `gorm:"not null" json:"price"`

	ProductItems []ProductItem `gorm:"constraint:OnDelete:CASCADE;"` //product删除要删除product_item

	InStockCount int64 `gorm:"-"`
}

func (*Product) TableName() string {
	return "product"
}
func (*Product) DefaultOrder() string {
	return "create_time DESC"
}
func (o *Product) ToDict() map[string]interface{} {
	return functions.StructToMap(o)
}

func (o *Product) BeforeCreate(tx *gorm.DB) (err error) {

	return
}
