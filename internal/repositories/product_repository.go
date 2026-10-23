package repositories

import (
	"github.com/abubakar-sadiq001/learning-go-shop/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(category *models.Category) error {
	if err := r.db.Create(&category).Error; err != nil {
		return err
	}
	return nil
}
func (r *ProductRepository) FindByIsActive() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Where("is_active = ?", true).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *ProductRepository) Update(category *models.Category) error {
	if err := r.db.Save(&category).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *ProductRepository) InsertNewProduct(product *models.Product) error {
	if err := r.db.Create(&product).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) ProductImageCount(productID uint) (*int64, error) {
	var count int64
	if err := r.db.Model(&models.ProductImage{}).Where("product_id = ?", productID).Count(&count).Error; err != nil {
		return nil, err
	}
	return &count, nil
}

func (r *ProductRepository) InsertProductImage(image *models.ProductImage) error {
	return r.db.Create(&image).Error
}

func (r *ProductRepository) DeleteProduct(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}
