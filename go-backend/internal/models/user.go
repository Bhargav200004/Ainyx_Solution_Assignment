package models

import (
	"math"
	"time"
)

// CreateUserRequest represents the JSON body for creating a user.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob"  validate:"required,dateformat"`
}

// UpdateUserRequest represents the JSON body for updating a user.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob"  validate:"required,dateformat"`
}

// UserResponse represents the JSON response for a single user.
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  *int   `json:"age,omitempty"`
}

// PaginatedResponse wraps a list of users with pagination metadata.
type PaginatedResponse struct {
	Data       []UserResponse `json:"data"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int64          `json:"total"`
	TotalPages int            `json:"total_pages"`
}

// ErrorResponse represents a standardized error JSON response.
type ErrorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

// DOBLayout is the expected date format for DOB fields.
const DOBLayout = "2006-01-02"

// CalculateAge computes the age in full years from a date of birth to today.
// It accounts for whether the birthday has occurred yet this year.
func CalculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()

	// Handle leap year edge case: if dob is Feb 29 and current year is not
	// a leap year, YearDay comparison may be off by one after Feb 28.
	// We use a direct month/day comparison as a more robust check.
	dobMonth := dob.Month()
	dobDay := dob.Day()
	nowMonth := now.Month()
	nowDay := now.Day()

	years = now.Year() - dob.Year()
	if nowMonth < dobMonth || (nowMonth == dobMonth && nowDay < dobDay) {
		years--
	}

	return int(math.Max(0, float64(years)))
}
