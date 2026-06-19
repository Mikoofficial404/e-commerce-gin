package handler

import (
	"ecommerce-gin/internal/service"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	MaxUploadSize = 1 << 20 // 1 MB
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: *svc}
}

// @Summary Create a new product
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Product Name"
// @Param description formData string true "Product Description"
// @Param price formData number true "Price"
// @Param stock formData integer true "Stock"
// @Param file formData file true "Product Image"
// @Success 200 {object} map[string]interface{} "Upload successful"
// @Router /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)
	if err := c.Request.ParseMultipartForm(MaxUploadSize); err != nil {
		if _, ok := err.(*http.MaxBytesError); ok {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": fmt.Sprintf("file too large (max: %d bytes)", MaxUploadSize),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file form required"})
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")

	price, err := strconv.ParseFloat(c.PostForm("price"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}

	stock, err := strconv.Atoi(c.PostForm("stock"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stock"})
		return
	}

	dst := filepath.Join("./public/uploads/", filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	product, err := h.productService.CreateProduct(name, description, price, stock, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "upload successful",
		"filename": file.Filename,
		"product":  product,
	})
}

// @Summary Get All Product with caching
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search query"
// @Success 200 {object} map[string]interface{} "Success"
// @Router /products [get]
func (h *ProductHandler) FindAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	isResult, total, err := h.productService.FindAllProduct(page, limit, search)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   isResult,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
