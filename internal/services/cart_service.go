package services

import (
	"errors"

	"github.com/abubakar-sadiq001/learning-go-shop/internal/dto"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/models"
	"gorm.io/gorm"
)

type CartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{
		db: db,
	}
}

func (s *CartService) GetCart(userID uint) (*dto.CartResponse, error) {
	var cart models.Cart
	if err := s.db.Preload("CartItems.Product.Category").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}

	return s.convertToCartResponse(&cart), nil
}

func (s *CartService) AddCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	// check if products exists
	var product models.Product
	if err := s.db.First(&product, req.ProductID).Error; err != nil {
		return nil, err
	}

	// check for stock availability
	if req.Quantity > product.Stock {
		return nil, errors.New("insufficient stock")
	}

	// search for existing cart. if not found, it creates a new cart
	var cart models.Cart
	if err := s.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		cart := models.Cart{UserID: userID}
		if err := s.db.Create(&cart).Error; err != nil {
			return nil, err
		}
	}

	// search cart item. if it exists, update, and if not, create
	var cartItem models.CartItem
	if err := s.db.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error; err != nil {
		// Create new cart item
		cartItem = models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := s.db.Create(&cartItem).Error; err != nil {
			return nil, err
		}
		if err := s.db.Save(&cartItem).Error; err != nil {
			return nil, err
		}

	} else {
		// Update existing cart item
		cartItem.Quantity += req.Quantity
		if cartItem.Quantity > product.Stock {
			return nil, errors.New("insufficient stock")
		}
		s.db.Save(&cartItem)
	}

	return s.GetCart(userID)
}

func (s *CartService) UpdateCartItem(userID, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	var cartItem models.CartItem
	if err := s.db.Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		First(&cartItem).Error; err != nil {
		return nil, errors.New("cart item not found")
	}

	var product models.Product
	if err := s.db.First(&product, cartItem.ProductID).Error; err != nil {
		return nil, errors.New("product not found")
	}

	if req.Quantity > product.Stock {
		return nil, errors.New("insufficient stock")
	}

	cartItem.Quantity = req.Quantity
	if err := s.db.Save(&cartItem).Error; err != nil {
		return nil, err
	}

	return s.GetCart(userID)
}

func (s *CartService) RemoveFromCart(userID, itemID uint) error {
	return s.db.Where("id = ? AND cart_id IN (?)", itemID,
		s.db.Select("id").Table("carts").
			Where("user_id = ?", userID)).
		Delete(&models.CartItem{}).Error
}

func (s *CartService) convertToCartResponse(cart *models.Cart) *dto.CartResponse {
	cartItems := make([]dto.CartItemResponse, len(cart.CartItems))
	var total float64

	for i := range cart.CartItems {
		subtotal := float64(cart.CartItems[i].Quantity) * cart.CartItems[i].Product.Price
		total += subtotal

		cartItems[i] = dto.CartItemResponse{
			ID: cart.CartItems[i].ID,
			Product: dto.ProductResponse{
				ID:          cart.CartItems[i].Product.ID,
				CategoryID:  cart.CartItems[i].Product.CategoryID,
				Name:        cart.CartItems[i].Product.Name,
				Description: cart.CartItems[i].Product.Description,
				Price:       cart.CartItems[i].Product.Price,
				Stock:       cart.CartItems[i].Product.Stock,
				SKU:         cart.CartItems[i].Product.SKU,
				IsActive:    cart.CartItems[i].Product.IsActive,
				CreatedAt:   cart.CartItems[i].Product.CreatedAt,
				UpdatedAt:   cart.CartItems[i].Product.UpdatedAt,
				Category: dto.CategoryResponse{
					ID:          cart.CartItems[i].Product.Category.ID,
					Name:        cart.CartItems[i].Product.Category.Name,
					Description: cart.CartItems[i].Product.Category.Description,
					IsActive:    cart.CartItems[i].Product.Category.IsActive,
				},
			},
			Quantity: cart.CartItems[i].Quantity,
			Total:    total,
		}
	}

	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItems,
		Total:     total,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}
