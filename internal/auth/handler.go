package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sq/ecommerce/pkg/response"
)

// Handler ทำหน้าที่รับ HTTP request และส่ง response
// รู้จัก Service แต่ไม่รู้จัก DB
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes ลงทะเบียน routes ทั้งหมดของ auth
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}

// Register godoc
// POST /api/auth/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	// Bind + validate ใน step เดียว (binding:"required" ใน model.go)
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Register(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "ลงทะเบียนสำเร็จ", result)
}

// Login godoc
// POST /api/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Login(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "เข้าสู่ระบบสำเร็จ", result)
}
