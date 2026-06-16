package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	sqlcdb "go-backend/db/sqlc/generated"
)

// ErrUserNotFound is returned when a user is not found in the database.
var ErrUserNotFound = errors.New("user not found")

// UserRepository defines the contract for user data access.
type UserRepository interface {
	Create(ctx context.Context, name string, dob pgtype.Date) (sqlcdb.User, error)
	GetByID(ctx context.Context, id int32) (sqlcdb.User, error)
	List(ctx context.Context, limit, offset int32) ([]sqlcdb.User, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, id int32, name string, dob pgtype.Date) (sqlcdb.User, error)
	Delete(ctx context.Context, id int32) error
}

// userRepository implements UserRepository using SQLC-generated queries.
type userRepository struct {
	queries *sqlcdb.Queries
}

// NewUserRepository creates a new UserRepository backed by SQLC queries.
func NewUserRepository(queries *sqlcdb.Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) Create(ctx context.Context, name string, dob pgtype.Date) (sqlcdb.User, error) {
	return r.queries.CreateUser(ctx, sqlcdb.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (sqlcdb.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlcdb.User{}, ErrUserNotFound
		}
		return sqlcdb.User{}, err
	}
	return user, nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int32) ([]sqlcdb.User, error) {
	return r.queries.ListUsers(ctx, sqlcdb.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}

func (r *userRepository) Update(ctx context.Context, id int32, name string, dob pgtype.Date) (sqlcdb.User, error) {
	user, err := r.queries.UpdateUser(ctx, sqlcdb.UpdateUserParams{
		ID:   id,
		Name: name,
		Dob:  dob,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlcdb.User{}, ErrUserNotFound
		}
		return sqlcdb.User{}, err
	}
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}
