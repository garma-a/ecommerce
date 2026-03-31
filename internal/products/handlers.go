package products

import (
	"ecom/internal/response"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (ph *ProductsHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := ph.service.GetProducts(r.Context())
	if err != nil {
		slog.Error("Error getting products", "error", err)
		response.InternalServerError(w)
		return
	}
	response.OK(w, products)
}

func (ph *ProductsHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	type CreateProductRequest struct {
		Name  string  `json:"name" validate:"required,min=2,max=120"`
		Price float64 `json:"price" validate:"required,gt=0,lte=1000000"`
	}

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request payload")
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
			response.ValidationError(w, "validation failed", fieldErrors)
			return
		}
		response.BadRequest(w, "validation error")
		return
	}

	product, err := ph.service.CreateProduct(r.Context(), req.Name, req.Price)
	if err != nil {
		// Check if it's a validation error from service layer
		var validationErr *ValidationError
		if errors.As(err, &validationErr) {
			response.BadRequest(w, validationErr.Error())
			return
		}

		// Internal/database errors
		slog.Error("Error creating product", "error", err)
		response.InternalServerError(w)
		return
	}

	response.Created(w, product)
}

func (ph *ProductsHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid product ID")
		return
	}

	product, err := ph.service.GetProductByID(r.Context(), int32(id))
	if err != nil {
		// Check for "not found" error (pgx returns specific error)
		if err.Error() == "no rows in result set" {
			response.NotFound(w, "product not found")
			return
		}

		slog.Error("Error getting product by ID", "error", err, "id", id)
		response.InternalServerError(w)
		return
	}

	response.OK(w, product)
}

func (ph *ProductsHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid product ID")
		return
	}

	err = ph.service.DeleteProduct(r.Context(), int32(id))
	if err != nil {
		slog.Error("Error deleting product", "error", err, "id", id)
		response.InternalServerError(w)
		return
	}

	response.NoContent(w)
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
