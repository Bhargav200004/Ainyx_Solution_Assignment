package handler

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"go-backend/internal/models"
	"go-backend/internal/repository"
	"go-backend/internal/service"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	service *service.UserService
	logger  *zap.Logger
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid request body",
		})
	}

	resp, err := h.service.CreateUser(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// GetUserByID handles GET /users/:id
func (h *UserHandler) GetUserByID(c fiber.Ctx) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid user id: must be a positive integer",
		})
	}

	resp, err := h.service.GetUserByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 10
	}

	resp, err := h.service.ListUsers(c.Context(), pageInt, limitInt)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid user id: must be a positive integer",
		})
	}

	var req models.UpdateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid request body",
		})
	}

	resp, err := h.service.UpdateUser(c.Context(), id, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid user id: must be a positive integer",
		})
	}

	if err := h.service.DeleteUser(c.Context(), id); err != nil {
		return h.handleError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// parseID extracts and validates the :id route parameter.
func (h *UserHandler) parseID(c fiber.Ctx) (int32, error) {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id")
	}
	return int32(id), nil
}

// handleError maps domain errors to appropriate HTTP responses.
func (h *UserHandler) handleError(c fiber.Ctx, err error) error {
	// Validation errors → 400
	if _, ok := err.(validator.ValidationErrors); ok {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation failed",
			Details: models.FormatValidationErrors(err),
		})
	}

	// Not found → 404
	if errors.Is(err, repository.ErrUserNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error: "user not found",
		})
	}

	// Everything else → 500
	h.logger.Error("internal server error", zap.Error(err))
	return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
		Error: "internal server error",
	})
}
