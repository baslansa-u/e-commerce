package order

import (
	"time"

	"github.com/sq/ecommerce/internal/product"
)

// OrderStatus คือสถานะของ order
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"   // รอชำระเงิน
	StatusPaid      OrderStatus = "paid"      // ชำระแล้ว
	StatusShipped   OrderStatus = "shipped"   // จัดส่งแล้ว
	StatusDelivered OrderStatus = "delivered" // ส่งถึงแล้ว
	StatusCancelled OrderStatus = "cancelled" // ยกเลิก
)

// allowedTransitions นิยามว่าจากสถานะหนึ่งไปสถานะใดได้บ้าง
// (delivered / cancelled เป็น terminal — เปลี่ยนต่อไม่ได้)
var allowedTransitions = map[OrderStatus][]OrderStatus{
	StatusPending: {StatusPaid, StatusCancelled},
	StatusPaid:    {StatusShipped, StatusCancelled},
	StatusShipped: {StatusDelivered},
}

// CanTransitionTo เช็คว่าเปลี่ยนจากสถานะปัจจุบันไป next ได้ไหม
func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	for _, allowed := range allowedTransitions[s] {
		if allowed == next {
			return true
		}
	}
	return false
}

// Order คือออเดอร์หลัก 1 ออเดอร์ต่อ 1 user
type Order struct {
	ID         uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint        `gorm:"not null" json:"user_id"`
	Status     OrderStatus `gorm:"not null;default:'pending'" json:"status"`
	TotalPrice float64     `gorm:"not null" json:"total_price"`
	Items      []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// OrderItem คือสินค้าแต่ละรายการในออเดอร์
// เก็บ snapshot ของราคา ณ ตอนที่สั่ง (ราคาอาจเปลี่ยนในอนาคต)
type OrderItem struct {
	ID        uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   uint            `gorm:"not null" json:"order_id"`
	ProductID uint            `gorm:"not null" json:"product_id"`
	Product   product.Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int             `gorm:"not null" json:"quantity"`
	Price     float64         `gorm:"not null" json:"price"` // snapshot ราคา ณ ตอนสั่ง
	CreatedAt time.Time       `json:"created_at"`
}

// --- Request/Response structs ---

// CheckoutItem คือสินค้าแต่ละชิ้นที่ส่งมาตอน checkout
type CheckoutItem struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// CheckoutRequest คือ body ที่รับตอน checkout
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items" binding:"required,min=1"`
}

// UpdateStatusRequest คือ body ที่ admin ส่งมาเปลี่ยนสถานะ order
type UpdateStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required,oneof=pending paid shipped delivered cancelled"`
}
