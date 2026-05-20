package order

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sq/ecommerce/pkg/middleware"
	"github.com/sq/ecommerce/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// ทุก route ของ order ต้อง login
	api := r.Group("/api").Use(middleware.AuthRequired())
	{
		api.POST("/orders/checkout", h.Checkout)
		api.GET("/orders", h.GetMyOrders)
		api.GET("/orders/:id", h.GetByID)
		api.POST("/orders/:id/pay", h.MockPayment)
	}
}

// POST /api/orders/checkout
func (h *Handler) Checkout(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.service.Checkout(userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "สั่งซื้อสำเร็จ", order)
}

// GET /api/orders
func (h *Handler) GetMyOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	orders, err := h.service.GetMyOrders(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ดึงข้อมูลออเดอร์สำเร็จ", orders)
}

// GET /api/orders/:id
func (h *Handler) GetByID(c *gin.Context) {
	userID := c.GetUint("user_id")
	orderID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id ไม่ถูกต้อง")
		return
	}

	order, err := h.service.GetOrderByID(orderID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ดึงข้อมูลออเดอร์สำเร็จ", order)
}

// POST /api/orders/:id/pay  — mock payment
func (h *Handler) MockPayment(c *gin.Context) {
	userID := c.GetUint("user_id")
	orderID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id ไม่ถูกต้อง")
		return
	}

	order, err := h.service.MockPayment(orderID, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ชำระเงินสำเร็จ", order)
}

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	return uint(id), err
}
