package auth

import (
	"errors"
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

// generateToken สร้าง JWT token
func (s *Service) generateToken(user *User) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // หมดอายุใน 24 ชั่วโมง
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
