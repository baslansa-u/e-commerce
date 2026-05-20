package product

import "errors"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// --- Product ---

func (s *Service) CreateProduct(req CreateProductRequest) (*Product, error) {
	// ตรวจว่า category มีอยู่จริง
	cat, err := s.repo.FindCategoryByID(req.CategoryID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, errors.New("ไม่พบ category นี้")
	}

	product := &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	// Preload category ก่อน return
	product.Category = *cat
	return product, nil
}

func (s *Service) GetAllProducts() ([]Product, error) {
	return s.repo.FindAll()
}

func (s *Service) GetProductByID(id uint) (*Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("ไม่พบสินค้านี้")
	}
	return product, nil
}

func (s *Service) UpdateProduct(id uint, req UpdateProductRequest) (*Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("ไม่พบสินค้านี้")
	}

	// อัปเดตเฉพาะ field ที่ส่งมา
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.CategoryID > 0 {
		// ตรวจว่า category ใหม่มีอยู่จริง
		cat, err := s.repo.FindCategoryByID(req.CategoryID)
		if err != nil {
			return nil, err
		}
		if cat == nil {
			return nil, errors.New("ไม่พบ category นี้")
		}
		product.CategoryID = req.CategoryID
	}

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Service) DeleteProduct(id uint) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("ไม่พบสินค้านี้")
	}
	return s.repo.Delete(id)
}

// --- Category ---

func (s *Service) CreateCategory(req CreateCategoryRequest) (*Category, error) {
	category := &Category{Name: req.Name}
	if err := s.repo.CreateCategory(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *Service) GetAllCategories() ([]Category, error) {
	return s.repo.FindAllCategories()
}
