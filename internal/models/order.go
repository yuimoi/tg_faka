package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"tg_go_faka/internal/utils/functions"
	_type "tg_go_faka/internal/utils/type"
)

// 使用指针可以方便的置空，使用原则：必须要判断是否为空的情况
type Order struct {
	//ID         int64             `gorm:"primaryKey;not null" json:"id"` // sqlite3库有个非常奇怪的逻辑，自增是primaryKey的默认自带，但是一旦附加了autoIncrease标签就会导致自增设置失效
	ID         uuid.UUID             `gorm:"primaryKey;not null" json:"id"`
	Status     _type.OrderStatusType `gorm:"default:0;not null" json:"status"`
	CreateTime int64                 `gorm:"index;autoCreateTime;not null" json:"create_time"`
	EndTime    int64                 `gorm:"index;not null" json:"end_time"`

	Price decimal.Decimal `gorm:"not null" json:"price"`

	TgID      int64 `gorm:"index;not null" json:"tg_id"` // 不能给unique，一个tg_id会创建多个订单
	MessageID int   `gorm:"index;not null" json:"message_id"`

	ProductItemID uuid.UUID   `json:"product_item_id"`
	ProductItem   ProductItem `gorm:"foreignKey:ProductItemID"`
}

func (*Order) TableName() string {
	return "order"
}
func (*Order) DefaultOrder() string {
	return "create_time DESC"
}
func (o *Order) ToDict() map[string]interface{} {
	return functions.StructToMap(o)
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}

func NewOrder(price decimal.Decimal, tgID int64, messageID int, endTime int64, productItemID uuid.UUID) *Order {
	order := &Order{
		Price:         price,
		TgID:          tgID,
		MessageID:     messageID,
		EndTime:       endTime,
		ProductItemID: productItemID,
	}
	return order
}
