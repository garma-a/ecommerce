package products

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductsHandler struct {
	service  IService
	validate *validator.Validate
}

// NewHandler allows dependency injection for easier testing.
func NewHandler(service IService) *ProductsHandler {
	return &ProductsHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (ph *ProductsHandler) GetProducts(c *gin.Context) {
	products, err := ph.service.GetProducts(c.Request.Context())
	if err != nil {
		slog.Error("Error getting products", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (ph *ProductsHandler) CreateProduct(c *gin.Context) {
	type CreateProductRequest struct {
		Name  string  `json:"name" validate:"required,min=2,max=120"`
		Price float64 `json:"price" validate:"required,gt=0,lte=1000000"`
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	// Validate request with validator tags
	if err := ph.validate.Struct(req); err != nil {
		validationErrs, ok := err.(validator.ValidationErrors)
		if ok {
			// Return field-level validation errors
			fieldErrors := make(map[string]string)
			for _, fieldErr := range validationErrs {
				fieldErrors[fieldErr.Field()] = formatValidationError(fieldErr)
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "validation failed",
				"fields": fieldErrors,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation error",
		})
		return
	}

	product, err := ph.service.CreateProduct(c.Request.Context(), req.Name, req.Price)
	if err != nil {
		// Check if it's a validation error from service layer
		var validationErr *ValidationError
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		// Internal/database errors
		slog.Error("Error creating product", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (ph *ProductsHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product ID",
		})
		return
	}

	product, err := ph.service.GetProductByID(c.Request.Context(), int32(id))
	if err != nil {
		// Check for "not found" error (pgx returns specific error)
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "product not found",
			})
			return
		}

		slog.Error("Error getting product by ID", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (ph *ProductsHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product ID",
		})
		return
	}

	err = ph.service.DeleteProduct(c.Request.Context(), int32(id))
	if err != nil {
		slog.Error("Error deleting product", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// formatValidationError formats validator errors into human-readable messages
func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "field is required"
	case "min":
		return "field must be at least " + err.Param() + " characters"
	case "max":
		return "field cannot exceed " + err.Param() + " characters"
	case "gt":
		return "field must be greater than " + err.Param()
	case "lte":
		return "field must be less than or equal to " + err.Param()
	default:
		return "field is invalid"
	}
}
