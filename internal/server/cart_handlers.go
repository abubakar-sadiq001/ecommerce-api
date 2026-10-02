package server

import (
	"strconv"

	"github.com/abubakar-sadiq001/learning-go-shop/internal/dto"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Tags Cart
// @Summary Get user's cart
// @Router /cart [get]
// @Description Retrieve current user's shopping cart with all items
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Cart retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Cart not found"
func (s *Server) getCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	cart, err := s.cartService.GetCart(userID)
	if err != nil {
		utils.NotFoundResponse(c, "Cart not found", nil)
		return
	}

	utils.SuccessResponse(c, "Cart retrieved successfully", cart)
}

// @Tags Cart
// @Summary Add a product to the user's shopping cart
// @Router /cart/items [post]
// @Description Retrieve current user's shopping cart with all items
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddToCartRequest true "Item to add to cart"
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Cart retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid request data or insufficient stock"
// @Failure 401 {object} utils.Response "Unauthorized"
func (s *Server) addToCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	cartResponse, err := s.cartService.AddCart(userID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to add item to cart", err)
		return
	}

	utils.SuccessResponse(c, "Item added to cart successfully", cartResponse)
}

// @Summary Update cart item quantity
// @Description Update the quantity of an item in the user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Param request body dto.UpdateCartItemRequest true "New quantity"
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Cart item updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data or insufficient stock"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /cart/items/{id} [put]
func (s *Server) updateCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid cart item ID", err)
		return
	}

	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Failed to update cart items", err)
		return
	}

	cart, err := s.cartService.UpdateCartItem(userID, uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update cart item", err)
		return
	}

	utils.SuccessResponse(c, "Cart item updated successfully", cart)
}

// @Summary Remove item from cart
// @Description Remove an item from the user's shopping cart
// @Tags Cart
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Success 200 {object} utils.Response "Item removed from cart successfully"
// @Failure 400 {object} utils.Response "Invalid cart item ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /cart/items/{id} [delete]
func (s *Server) removeFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid cart item ID", err)
		return
	}

	err = s.cartService.RemoveFromCart(userID, uint(id))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove item", err)
		return
	}

	utils.SuccessResponse(c, "Item removed successfully", nil)
}
