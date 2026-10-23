package services

import (
	"github.com/abubakar-sadiq001/learning-go-shop/internal/dto"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/models"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/repositories"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/utils"
	"gorm.io/gorm"
)

var _ ProductServiceInterface = (*ProductService)(nil)

type ProductService struct {
	db          *gorm.DB
	productRepo repositories.ProductRepositoryInterface
}

func NewProductService(db *gorm.DB, productRepo repositories.ProductRepositoryInterface) *ProductService {
	return &ProductService{
		db:          db,
		productRepo: productRepo,
	}
}

func (s *ProductService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	err := s.productRepo.Create(&category)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil

}

func (s *ProductService) GetCategories() ([]dto.CategoryResponse, error) {

	categories, err := s.productRepo.FindByIsActive()
	if err != nil {
		return nil, err
	}

	response := make([]dto.CategoryResponse, len(categories))
	for i := range categories {
		response[i] = dto.CategoryResponse{
			ID:          categories[i].ID,
			Name:        categories[i].Name,
			Description: categories[i].Description,
			IsActive:    categories[i].IsActive,
			CreatedAt:   categories[i].CreatedAt,
			UpdatedAt:   categories[i].UpdatedAt,
		}
	}

	return response, nil
}

func (s *ProductService) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	err = s.productRepo.Update(category)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (s *ProductService) DeleteCategory(id uint) error {
	err := s.productRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}

	if err := s.productRepo.InsertNewProduct(&product); err != nil {
		return nil, err
	}

	return s.GetProduct(product.ID)
}

func (s *ProductService) GetProducts(page, pageLimit int) ([]dto.ProductResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}

	if pageLimit < 1 {
		pageLimit = 10
	}

	offset := (page - 1) * pageLimit
	var products []models.Product

	var total int64
	if err := s.db.Model(&models.Product{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, nil, err
	}

	if err := s.db.Preload("Category").Preload("Images").Where("is_active = ?", true).Offset(offset).Limit(pageLimit).Find(&products).Error; err != nil {
		return nil, nil, err
	}

	response := make([]dto.ProductResponse, len(products))
	for i := range products {
		response[i] = s.convertToProductResponse(&products[i])
	}

	totalPages := int((total + int64(pageLimit) - 1) / int64(pageLimit))
	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      pageLimit,
		Total:      total,
		TotalPages: int(totalPages),
	}

	return response, meta, nil
}

func (s *ProductService) GetProduct(id uint) (*dto.ProductResponse, error) {
	var product models.Product
	if err := s.db.Preload("Category").Preload("Images").First(&product, id).Error; err != nil {
		return nil, err
	}

	response := s.convertToProductResponse(&product)

	return &response, nil
}

func (s *ProductService) UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.db.Save(&product).Error; err != nil {
		return nil, err
	}

	response := &dto.ProductResponse{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		IsActive:    product.IsActive,
	}

	return response, nil
}

// AddProductImage saves the image data as new record to the [product_images] table.
func (s *ProductService) AddProductImage(productID uint, url, altText string) error {
	count, err := s.productRepo.ProductImageCount(productID)
	if err != nil {
		return err
	}

	image := models.ProductImage{
		ProductID: productID,
		URL:       url,
		AltText:   altText,
		IsPrimary: int(*count) == 0,
	}

	if err := s.productRepo.InsertProductImage(&image); err != nil {
		return err
	}
	return nil
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.productRepo.DeleteProduct(id)
}

func (s *ProductService) convertToProductResponse(product *models.Product) dto.ProductResponse {
	images := make([]dto.ProductImageResponse, len(product.Images))
	for i := range product.Images {
		images[i] = dto.ProductImageResponse{
			ID:        product.Images[i].ID,
			URL:       product.Images[i].URL,
			AltText:   product.Images[i].AltText,
			IsPrimary: product.Images[i].IsPrimary,
			CreatedAt: product.Images[i].CreatedAt,
		}
	}

	return dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		Category: dto.CategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Name,
			Description: product.Description,
			IsActive:    product.IsActive,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		},
		Images:    images,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}
