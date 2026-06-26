package order

import (
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateOrder บันทึก order พร้อม items ในครั้งเดียว (transaction)
func (r *Repository) CreateOrder(order *Order) error {
	// ใช้ transaction ป้องกันข้อมูลหาย ถ้า step ใด step หนึ่ง fail
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

// FindByID ดึง order พร้อม items และ product ของแต่ละ item
func (r *Repository) FindByID(id uint) (*Order, error) {
	var order Order
	err := r.db.
		Preload("Items.Product.Category").
		First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// FindByUserID ดึง order ทั้งหมดของ user คนนั้น
func (r *Repository) FindByUserID(userID uint) ([]Order, error) {
	var orders []Order
	err := r.db.
		Preload("Items.Product").
		Where("user_id = ?", userID).
		Order("created_at DESC"). // ล่าสุดขึ้นก่อน
		Find(&orders).Error
	return orders, err
}

// FindAll ดึง order ทั้งหมดในระบบ (สำหรับ admin)
func (r *Repository) FindAll() ([]Order, error) {
	var orders []Order
	err := r.db.
		Preload("Items.Product").
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// UpdateStatus เปลี่ยนสถานะ order
func (r *Repository) UpdateStatus(id uint, status OrderStatus) error {
	return r.db.Model(&Order{}).Where("id = ?", id).Update("status", status).Error
}

// DecreaseStock ลด stock สินค้า (เรียกใน transaction เดียวกับ CreateOrder)
func (r *Repository) DecreaseStock(productID uint, quantity int) error {
	result := r.db.Exec(
		"UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ?",
		quantity, productID, quantity,
	)
	if result.Error != nil {
		return result.Error
	}
	// ถ้า RowsAffected = 0 แปลว่า stock ไม่พอ
	if result.RowsAffected == 0 {
		return errors.New("stock ไม่เพียงพอ")
	}
	return nil
}
