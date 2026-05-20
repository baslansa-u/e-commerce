package product

import (
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// --- Product ---

func (r *Repository) Create(p *Product) error {
	return r.db.Create(p).Error
}

// FindAll ดึงสินค้าทั้งหมด พร้อม preload Category
func (r *Repository) FindAll() ([]Product, error) {
	var products []Product
	err := r.db.Preload("Category").Find(&products).Error
	return products, err
}

func (r *Repository) FindByID(id uint) (*Product, error) {
	var product Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *Repository) Update(product *Product) error {
	return r.db.Save(product).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Product{}, id).Error
}

// --- Category ---

func (r *Repository) CreateCategory(c *Category) error {
	return r.db.Create(c).Error
}

func (r *Repository) FindAllCategories() ([]Category, error) {
	var categories []Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *Repository) FindCategoryByID(id uint) (*Category, error) {
	var category Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}
