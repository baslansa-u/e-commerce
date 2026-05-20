package auth

import "time"

// User คือ struct หลักที่ mapping กับตาราง users ในฐานข้อมูล
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // json:"-" = ไม่ส่งรหัสผ่านกลับใน response
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
