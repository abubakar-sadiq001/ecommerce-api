package server

import (
	"strconv"

	"github.com/abubakar-sadiq001/learning-go-shop/internal/dto"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/services"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) getCategories(c *gin.Context) {
	productService := services.NewProductServices(s.db)
	response, err := productService.GetCategories()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve categories", err)
		return
	}

	utils.SuccessResponse(c, "Categories", response)
}

func (s *Server) updateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if id != 0 {
		utils.BadRequestResponse(c, "Invalid category id", err)
		return
	}

	var req dto.UpdatedCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	productService := services.NewProductServices(s.db)
	response, err := productService.GetCategories()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update category", err)
		return
	}

	utils.SuccessResponse(c, "Category has been updated successfully", response)
}

func (s *Server) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if id != 0 {
		utils.BadRequestResponse(c, "Invalid category id", err)
		return
	}

	productService := services.NewProductServices(s.db)
	err = productService.DeleteCategory(uint(id))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete category", err)
		return
	}

	utils.SuccessResponse(c, "Category has been deleted successfully", nil)
}

func (s *Server) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	productService := services.NewProductServices(s.db)
	response, err := productService.CreateCategory(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create category", err)
		return
	}

	utils.SuccessResponse(c, "Category created successfuly", response)

}

func (s *Server) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	productService := services.NewProductServices(s.db)
	response, err := productService.CreateProduct(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create product", err)
		return
	}

	utils.SuccessResponse(c, "Product has been created successfully", response)
}

func (s *Server) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	productService := services.NewProductServices(s.db)
	response, meta, err := productService.GetProducts(page, pageLimit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve products", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Product retrieved successfully", response, *meta)
}

func (s *Server) getProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if id != 0 {
		utils.BadRequestResponse(c, "Invalid category id", err)
		return
	}

	productService := services.NewProductServices(s.db)
	response, err := productService.GetProduct(uint(id))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve product", err)
		return
	}

	utils.SuccessResponse(c, "Product retrieved successfully", response)
}

func (s *Server) updateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if id != 0 {
		utils.BadRequestResponse(c, "Invalid category id", err)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	productService := services.NewProductServices(s.db)
	response, err := productService.UpdateProduct(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update product", err)
		return
	}

	utils.SuccessResponse(c, "Product has been updated successfully", response)
}

func (s *Server) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid category id", err)
		return
	}

	productService := services.NewProductServices(s.db)
	if err := productService.DeleteCategory(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "Faled to delete product", err)
		return
	}

	utils.SuccessResponse(c, "Product has been deleted successfully", nil)

}
