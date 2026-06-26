package auth

import (
	"errors"

	"gorm.io/gorm"
)

// Repository ทำหน้าที่ในการติดต่อกับฐานข้อมูล อย่างเดียว
// ไม่มี business logic ใดๆ อยู่ในนี้
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateUser บันทึกผู้ใช้ใหม่ลงในฐานข้อมูล
func (r *Repository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

// UpdateUser บันทึกการเปลี่ยนแปลงของ user ที่มีอยู่แล้ว
func (r *Repository) UpdateUser(user *User) error {
	return r.db.Save(user).Error
}

// FindByEmail ค้นหาผู้ใช้จากอีเมล
func (r *Repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // ไม่พบผู้ใช้
		}
		return nil, err // เกิดข้อผิดพลาดอื่นๆ
	}
	return &user, nil

}

// FindByID ค้นหาผู้ใช้จาก ID
func (r *Repository) FindByID(id uint) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // ไม่พบผู้ใช้
		}
		return nil, err // เกิดข้อผิดพลาดอื่นๆ
	}
	return &user, nil
}
