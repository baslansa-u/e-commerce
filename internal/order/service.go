package order

import (
	"errors"
	"fmt"

	"github.com/sq/ecommerce/internal/product"
	"gorm.io/gorm"
)

type Service struct {
	repo        *Repository
	productRepo *product.Repository // ใช้ดึงข้อมูลสินค้า
	db          *gorm.DB            // ใช้สำหรับ transaction
}

func NewService(repo *Repository, productRepo *product.Repository, db *gorm.DB) *Service {
	return &Service{repo: repo, productRepo: productRepo, db: db}
}

// Checkout สร้าง order จาก items ที่ส่งมา
func (s *Service) Checkout(userID uint, req CheckoutRequest) (*Order, error) {
	var orderItems []OrderItem
	var totalPrice float64

	// ใช้ transaction ทั้งหมด: ตรวจ stock → ลด stock → สร้าง order
	// ถ้า step ใดผิดพลาด rollback ทั้งหมด
	err := s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := NewRepository(tx)

		for _, item := range req.Items {
			// ดึงข้อมูลสินค้า
			p, err := s.productRepo.FindByID(item.ProductID)
			if err != nil {
				return err
			}
			if p == nil {
				return errors.New("ไม่พบสินค้า id: " + itoa(item.ProductID))
			}

			// ตรวจ stock
			if p.Stock < item.Quantity {
				return errors.New("สินค้า \"" + p.Name + "\" stock ไม่เพียงพอ")
			}

			// ลด stock
			if err := txRepo.DecreaseStock(item.ProductID, item.Quantity); err != nil {
				return err
			}

			orderItems = append(orderItems, OrderItem{
				ProductID: item.ProductID,
				Product:   *p,
				Quantity:  item.Quantity,
				Price:     p.Price, // snapshot ราคา ณ ตอนนี้
			})

			totalPrice += p.Price * float64(item.Quantity)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// สร้าง order หลังจาก transaction stock สำเร็จ
	order := &Order{
		UserID:     userID,
		Status:     StatusPending,
		TotalPrice: totalPrice,
		Items:      orderItems,
	}

	if err := s.repo.CreateOrder(order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetMyOrders ดึง order ทั้งหมดของ user
func (s *Service) GetMyOrders(userID uint) ([]Order, error) {
	return s.repo.FindByUserID(userID)
}

// GetOrderByID ดึง order โดย id และตรวจว่าเป็นของ user นั้นจริง
func (s *Service) GetOrderByID(orderID, userID uint) (*Order, error) {
	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("ไม่พบออเดอร์นี้")
	}
	// ป้องกัน user อื่นดู order ของคนอื่น
	if order.UserID != userID {
		return nil, errors.New("ไม่มีสิทธิ์ดูออเดอร์นี้")
	}
	return order, nil
}

// MockPayment จำลองการชำระเงิน (เปลี่ยน status pending → paid)
func (s *Service) MockPayment(orderID, userID uint) (*Order, error) {
	order, err := s.GetOrderByID(orderID, userID)
	if err != nil {
		return nil, err
	}
	if order.Status != StatusPending {
		return nil, errors.New("ออเดอร์นี้ไม่อยู่ในสถานะรอชำระเงิน")
	}

	if err := s.repo.UpdateStatus(orderID, StatusPaid); err != nil {
		return nil, err
	}

	order.Status = StatusPaid
	return order, nil
}

// itoa แปลง uint เป็น string แบบง่าย
func itoa(n uint) string {
	return fmt.Sprintf("%d", n)
}
