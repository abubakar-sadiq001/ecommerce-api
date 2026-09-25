package server

import (
	"strconv"

	"github.com/abubakar-sadiq001/learning-go-shop/internal/dto"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/services"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) getCategories(c *gin.Context) {
	response, err := s.productService.GetCategories()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve categories", err)
		return
	}

	utils.SuccessResponse(c, "Categories", response)
}

func (s *Server) updateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid category id", err)
		return
	}

	var req dto.UpdatedCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	response, err := s.productService.UpdateCategory(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update category", err)
		return
	}

	utils.SuccessResponse(c, "Category has been updated successfully", response)
}

func (s *Server) deleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid category id", err)
		return
	}

	err = s.productService.DeleteCategory(uint(id))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete category", err)
		return
	}

	utils.SuccessResponse(c, "Category has been deleted successfully", nil)
}

func (s *Server) createCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	response, err := s.productService.CreateCategory(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create category", err)
		return
	}

	utils.CreatedResponse(c, "Category created successfully", response)

}

func (s *Server) createProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	response, err := s.productService.CreateProduct(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create product", err)
		return
	}

	utils.CreatedResponse(c, "Product has been created successfully", response)
}

func (s *Server) getProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	response, meta, err := s.productService.GetProducts(page, pageLimit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve products", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Product retrieved successfully", response, *meta)
}

func (s *Server) getProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid category id", err)
		return
	}

	response, err := s.productService.GetProduct(uint(id))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve product", err)
		return
	}

	utils.SuccessResponse(c, "Product retrieved successfully", response)
}

func (s *Server) updateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid category id", err)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	response, err := s.productService.UpdateProduct(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update product", err)
		return
	}

	utils.SuccessResponse(c, "Product has been updated successfully", response)
}

func (s *Server) deleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product id", err)
		return
	}

	if err := s.productService.DeleteProduct(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "Faled to delete product", err)
		return
	}

	utils.SuccessResponse(c, "Product has been deleted successfully", nil)

}

func (s *Server) uploadProductImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product id", err)
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		utils.BadRequestResponse(c, "No file is uploaded", err)
		return
	}

	url, err := s.uploadService.UploadImage(uint(id), file)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to upload file", err)
		return
	}

	productService := services.NewProductService(s.db)
	if err := productService.AddProductImage(uint(id), url, file.Filename); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to save image record", err)
		return
	}

	utils.SuccessResponse(c, "Image uploaded successfully", url)

}
