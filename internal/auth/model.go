package auth

import "time"

// Role คือสิทธิ์ของผู้ใช้ในระบบ
type Role string

const (
	RoleCustomer Role = "customer" // ลูกค้าทั่วไป — ดูสินค้า, สั่งซื้อ, ดู order ตัวเอง
	RoleAdmin    Role = "admin"    // ผู้ดูแล — จัดการสินค้า/หมวดหมู่ และ order ทั้งหมด
)

// User คือ struct หลักที่ mapping กับตาราง users ในฐานข้อมูล
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`                          // json:"-" = ไม่ส่งรหัสผ่านกลับใน response
	Role      Role      `gorm:"type:varchar(20);not null;default:customer" json:"role"` // สิทธิ์: customer (default) หรือ admin
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterRequest คือ struct ที่ใช้สำหรับรับข้อมูลจาก client เมื่อมีการลงทะเบียน
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest คือ struct ที่ใช้สำหรับรับข้อมูลจาก client เมื่อมีการเข้าสู่ระบบ
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse คือ struct ที่ใช้สำหรับส่งข้อมูลกลับไปยัง client หลังจากการเข้าสู่ระบบสำเร็จ
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
