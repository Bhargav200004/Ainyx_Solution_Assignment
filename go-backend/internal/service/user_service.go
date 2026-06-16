package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"

	"go-backend/internal/models"
	"go-backend/internal/repository"
)

// UserService provides business logic for user operations.
type UserService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

// NewUserService creates a new UserService.
func NewUserService(repo repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

// CreateUser validates input, creates a user, and returns the response.
func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	if err := models.Validate(req); err != nil {
		return nil, err
	}

	dob, err := time.Parse(models.DOBLayout, req.DOB)
	if err != nil {
		return nil, err
	}

	pgDate := pgtype.Date{
		Time:  dob,
		Valid: true,
	}

	user, err := s.repo.Create(ctx, req.Name, pgDate)
	if err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return nil, err
	}

	s.logger.Info("user created",
		zap.Int32("id", user.ID),
		zap.String("name", user.Name),
	)

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Time.Format(models.DOBLayout),
	}, nil
}

// GetUserByID retrieves a user by ID and calculates their age.
func (s *UserService) GetUserByID(ctx context.Context, id int32) (*models.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, err
		}
		s.logger.Error("failed to get user", zap.Int32("id", id), zap.Error(err))
		return nil, err
	}

	age := models.CalculateAge(user.Dob.Time)

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Time.Format(models.DOBLayout),
		Age:  &age,
	}, nil
}

// ListUsers returns a paginated list of users with calculated ages.
func (s *UserService) ListUsers(ctx context.Context, page, limit int) (*models.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	users, err := s.repo.List(ctx, int32(limit), int32(offset))
	if err != nil {
		s.logger.Error("failed to list users", zap.Error(err))
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		s.logger.Error("failed to count users", zap.Error(err))
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	responses := make([]models.UserResponse, 0, len(users))
	for _, u := range users {
		age := models.CalculateAge(u.Dob.Time)
		responses = append(responses, models.UserResponse{
			ID:   u.ID,
			Name: u.Name,
			DOB:  u.Dob.Time.Format(models.DOBLayout),
			Age:  &age,
		})
	}

	return &models.PaginatedResponse{
		Data:       responses,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// UpdateUser validates input, updates a user, and returns the response.
func (s *UserService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error) {
	if err := models.Validate(req); err != nil {
		return nil, err
	}

	dob, err := time.Parse(models.DOBLayout, req.DOB)
	if err != nil {
		return nil, err
	}

	pgDate := pgtype.Date{
		Time:  dob,
		Valid: true,
	}

	user, err := s.repo.Update(ctx, id, req.Name, pgDate)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, err
		}
		s.logger.Error("failed to update user", zap.Int32("id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("user updated",
		zap.Int32("id", user.ID),
		zap.String("name", user.Name),
	)

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Time.Format(models.DOBLayout),
	}, nil
}

// DeleteUser removes a user by ID.
func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	// Check if user exists first.
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return err
		}
		s.logger.Error("failed to check user before delete", zap.Int32("id", id), zap.Error(err))
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete user", zap.Int32("id", id), zap.Error(err))
		return err
	}

	s.logger.Info("user deleted", zap.Int32("id", id))
	return nil
}
