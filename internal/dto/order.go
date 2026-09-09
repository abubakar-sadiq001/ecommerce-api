package dto

type AddToCartRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type UpdateCartItemsRequest struct {
	ID        uint               `json:"id"`
	UserID    uint               `json:"user_id"`
	CartItems []CartItemResponse `json:"cart_items"`
	Total     float64            `json:"total"`
}

type CartItemResponse struct {
	ID       uint            `json:"id"`
	Product  ProductResponse `json:"product_response"`
	Quantity int             `json:"quantity"`
	Total    float64         `json:"total"`
}

type OrderItemResponse struct {
	ID       uint            `json:"id"`
	Product  ProductResponse `json:"product_response"`
	Quantity int             `json:"quantity"`
	Price    float64         `json:"price"`
}
