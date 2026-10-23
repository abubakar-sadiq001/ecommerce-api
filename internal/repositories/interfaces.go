package repositories

import "github.com/abubakar-sadiq001/learning-go-shop/internal/models"

type UserRepositoryInterface interface {
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmailAndActive(email string, isActive bool) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error

	CreateRefreshToken(token *models.RefreshToken) error
	GetValidRefreshToken(token string) (*models.RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokenByID(id uint) error
	DeleteExpiredTokenByID(id uint) error
}

type CartRepositoryInterface interface {
	GetByUserID(userID uint) (*models.Cart, error)
	Create(cart *models.Cart) error
	Update(cart *models.Cart) error
	Delete(id uint) error
}

type ProductRepositoryInterface interface {
	Create(category *models.Category) error
	FindByIsActive() ([]models.Category, error)
	FindByID(id uint) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
	InsertNewProduct(product *models.Product) error
	ProductImageCount(productID uint) (*int64, error)
	InsertProductImage(image *models.ProductImage) error
	DeleteProduct(id uint) error
}
