package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Service ทำหน้าที่จัดการ business logic
// รู้จัก Repository แต่ไม่รู้จัก HTTP
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Register สร้าง user ใหม่
func (s *Service) Register(req RegisterRequest) (*AuthResponse, error) {
	// ตรวจว่า email ซ้ำไหม
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email นี้ถูกใช้งานแล้ว")
	}

	// Hash password ก่อนบันทึก (ห้าม store plain text เด็ดขาด)
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		Role:     RoleCustomer, // สมัครใหม่เป็น customer เสมอ — admin ต้อง provision แยก
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	// สร้าง token ให้เลย ไม่ต้อง login แยก
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: *user}, nil
}

// Login ตรวจสอบ credentials แล้วคืน token
func (s *Service) Login(req LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	// ใช้ข้อความเดียวกันทั้ง email ไม่เจอ และ password ผิด
	// เพื่อป้องกัน user enumeration attack
	if user == nil {
		return nil, errors.New("email หรือ password ไม่ถูกต้อง")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("email หรือ password ไม่ถูกต้อง")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: *user}, nil
}

// EnsureAdmin สร้าง admin คนแรกจาก env (ADMIN_EMAIL, ADMIN_PASSWORD, ADMIN_NAME)
// เรียกตอน startup — idempotent: ไม่มีก็สร้าง, มีแต่ยังไม่ใช่ admin ก็ promote, ตั้งครบแล้วก็ข้าม
func (s *Service) EnsureAdmin() error {
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	if email == "" || password == "" {
		log.Println("EnsureAdmin: ไม่ได้ตั้ง ADMIN_EMAIL/ADMIN_PASSWORD — ข้ามการสร้าง admin")
		return nil
	}

	existing, err := s.repo.FindByEmail(email)
	if err != nil {
		return err
	}

	// มี user นี้อยู่แล้ว — promote เป็น admin ถ้ายังไม่ใช่ (ไม่ยุ่งกับ password เดิม)
	if existing != nil {
		if existing.Role == RoleAdmin {
			return nil
		}
		existing.Role = RoleAdmin
		if err := s.repo.UpdateUser(existing); err != nil {
			return err
		}
		log.Printf("EnsureAdmin: promote %s เป็น admin แล้ว", email)
		return nil
	}

	// ยังไม่มี — สร้างใหม่เป็น admin
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	name := os.Getenv("ADMIN_NAME")
	if name == "" {
		name = "Admin"
	}
	admin := &User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
		Role:     RoleAdmin,
	}
	if err := s.repo.CreateUser(admin); err != nil {
		return err
	}
	log.Printf("EnsureAdmin: สร้าง admin %s สำเร็จ", email)
	return nil
}

// generateToken สร้าง JWT token
func (s *Service) generateToken(user *User) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,                            // ฝัง role ไว้ใน token ให้ middleware เช็คสิทธิ์
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // หมดอายุใน 24 ชั่วโมง
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
