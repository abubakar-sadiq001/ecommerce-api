package server

import (
	"strconv"

	"github.com/abubakar-sadiq001/learning-go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	order, err := s.orderService.CreateOrder(userID)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create order", err)
		return
	}

	utils.CreatedResponse(c, "Order created successfully", order)
}

func (s *Server) getOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, meta, err := s.orderService.GetOrders(userID, page, pageLimit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve orders", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Orders successfully retrieved ", orders, *meta)
}

func (s *Server) getOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	orderID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid order ID", err)
		return
	}

	order, err := s.orderService.GetOrder(userID, uint(orderID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Order not found", err)
		return
	}

	utils.SuccessResponse(c, "Order successfully retrieved ", order)
}
