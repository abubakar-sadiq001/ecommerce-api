package server

import (
	"strconv"

	"github.com/abubakar-sadiq001/learning-go-shop/internal/dto"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) getCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	cart, err := s.cartService.GetCart(userID)
	if err != nil {
		utils.NotFoundResponse(c, "Cart not found", nil)
		return
	}

	utils.SuccessResponse(c, "Cart retrieved successfully", cart)
}

func (s *Server) addToCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	cartResponse, err := s.cartService.AddCart(userID, &req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to add item to cart", err)
		return
	}

	utils.SuccessResponse(c, "Item added to cart successfully", cartResponse)
}

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
		utils.BadRequestResponse(c, "Failed to update cart item", err)
		return
	}

	utils.SuccessResponse(c, "Cart item updated successfully", cart)
}

func (s *Server) removeFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid cart item ID", err)
		return
	}

	err = s.cartService.RemoveFromCart(userID, uint(id))
	if err != nil {
		utils.BadRequestResponse(c, "Failed to remove item", err)
		return
	}

	utils.SuccessResponse(c, "Item removed successfully", nil)
}
