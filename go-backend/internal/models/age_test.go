package models

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "30 years ago, birthday already passed this year",
			dob:      time.Date(now.Year()-30, now.Month()-1, 1, 0, 0, 0, 0, time.UTC),
			expected: 30,
		},
		{
			name:     "25 years ago, birthday is today",
			dob:      time.Date(now.Year()-25, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			expected: 25,
		},
		{
			name:     "20 years ago, birthday hasn't come yet this year",
			dob:      addMonthSafe(now, -20*12+1),
			expected: 19,
		},
		{
			name:     "newborn, born today",
			dob:      time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "1 year old, born exactly 1 year ago",
			dob:      time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "leap year DOB Feb 29 2000",
			dob:      time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: calculateExpectedAge(time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:     "born on Jan 1 of 30 years ago",
			dob:      time.Date(now.Year()-30, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: calculateExpectedAge(time.Date(now.Year()-30, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:     "born on Dec 31 of 30 years ago",
			dob:      time.Date(now.Year()-30, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: calculateExpectedAge(time.Date(now.Year()-30, 12, 31, 0, 0, 0, 0, time.UTC)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAge(tt.dob)
			if got != tt.expected {
				t.Errorf("CalculateAge(%v) = %d, want %d", tt.dob.Format("2006-01-02"), got, tt.expected)
			}
		})
	}
}

// addMonthSafe adds months to a time, handling month-end clamping.
func addMonthSafe(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

// calculateExpectedAge is a reference implementation for validation.
func calculateExpectedAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}
	if years < 0 {
		return 0
	}
	return years
}

func TestCalculateAge_NegativeReturnsZero(t *testing.T) {
	// Future date should return 0.
	future := time.Now().AddDate(1, 0, 0)
	got := CalculateAge(future)
	if got != 0 {
		t.Errorf("CalculateAge(future) = %d, want 0", got)
	}
}
