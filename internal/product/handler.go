package product

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
	api := r.Group("/api")

	// Public routes — ดูสินค้า/หมวดหมู่ได้โดยไม่ต้อง login
	products := api.Group("/products")
	{
		products.GET("", h.GetAll)
		products.GET("/:id", h.GetByID)
	}
	api.GET("/categories", h.GetAllCategories)

	// Admin-only routes — ต้อง login และเป็น admin เท่านั้น (เพิ่ม/แก้/ลบสินค้า, สร้าง category)
	admin := api.Group("")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		admin.POST("/products", h.Create)
		admin.PUT("/products/:id", h.Update)
		admin.DELETE("/products/:id", h.Delete)

		admin.POST("/categories", h.CreateCategory)
	}
}

// GET /api/products
func (h *Handler) GetAll(c *gin.Context) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "ดึงข้อมูลสินค้าสำเร็จ", products)
}

// GET /api/products/:id
func (h *Handler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id ไม่ถูกต้อง")
		return
	}

	product, err := h.service.GetProductByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "ดึงข้อมูลสินค้าสำเร็จ", product)
}

// POST /api/products
func (h *Handler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.service.CreateProduct(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "เพิ่มสินค้าสำเร็จ", product)
}

// PUT /api/products/:id
func (h *Handler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id ไม่ถูกต้อง")
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.service.UpdateProduct(id, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "อัปเดตสินค้าสำเร็จ", product)
}

// DELETE /api/products/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id ไม่ถูกต้อง")
		return
	}

	if err := h.service.DeleteProduct(id); err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "ลบสินค้าสำเร็จ", nil)
}

// POST /api/categories
func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.CreateCategory(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "เพิ่ม category สำเร็จ", category)
}

// GET /api/categories
func (h *Handler) GetAllCategories(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "ดึงข้อมูล category สำเร็จ", categories)
}

// parseID แปลง string param เป็น uint
func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	return uint(id), err
}
