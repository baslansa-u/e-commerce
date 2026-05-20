package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sq/ecommerce/internal/auth"
	"github.com/sq/ecommerce/internal/order"
	"github.com/sq/ecommerce/internal/product"
	"github.com/sq/ecommerce/pkg/database"
)

func main() {
	// โหลด .env ไฟล์
	if err := godotenv.Load(); err != nil {
		log.Fatal("Fail load .env file: ", err)
	}

	// เชื่อมต่อฐานข้อมูล
	db := database.Connect()

	// Auto Migrate สร้างตารางอัตโนมัติจาก struct
	db.AutoMigrate(&auth.User{}, &product.Category{}, &product.Product{}, &order.Order{}, &order.OrderItem{})

	// Init layers แบบ manual dependency injection
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	// Setup router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Rigister routes
	authHandler.RegisterRoutes(r)

	productRepo := product.NewRepository(db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)
	productHandler.RegisterRoutes(r)

	orderRepo := order.NewRepository(db)
	orderService := order.NewService(orderRepo, productRepo, db)
	orderHandler := order.NewHandler(orderService)
	orderHandler.RegisterRoutes(r)

	log.Println("Server start 8080")
	r.Run(":8080")
}
