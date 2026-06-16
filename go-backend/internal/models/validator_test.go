package models

import (
	"testing"
)

func TestValidate_CreateUserRequest_Valid(t *testing.T) {
	req := CreateUserRequest{
		Name: "Alice",
		DOB:  "1990-05-10",
	}
	if err := Validate(req); err != nil {
		t.Errorf("expected no error for valid request, got: %v", err)
	}
}

func TestValidate_CreateUserRequest_MissingName(t *testing.T) {
	req := CreateUserRequest{
		Name: "",
		DOB:  "1990-05-10",
	}
	if err := Validate(req); err == nil {
		t.Error("expected error for missing name, got nil")
	}
}

func TestValidate_CreateUserRequest_MissingDOB(t *testing.T) {
	req := CreateUserRequest{
		Name: "Alice",
		DOB:  "",
	}
	if err := Validate(req); err == nil {
		t.Error("expected error for missing dob, got nil")
	}
}

func TestValidate_CreateUserRequest_InvalidDateFormat(t *testing.T) {
	tests := []struct {
		name string
		dob  string
	}{
		{"wrong format MM/DD/YYYY", "05/10/1990"},
		{"wrong format DD-MM-YYYY", "10-05-1990"},
		{"invalid date", "1990-13-01"},
		{"invalid day", "1990-02-30"},
		{"random string", "not-a-date"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateUserRequest{
				Name: "Alice",
				DOB:  tt.dob,
			}
			if err := Validate(req); err == nil {
				t.Errorf("expected error for dob=%q, got nil", tt.dob)
			}
		})
	}
}

func TestValidate_CreateUserRequest_FutureDOB(t *testing.T) {
	req := CreateUserRequest{
		Name: "Alice",
		DOB:  "2099-01-01",
	}
	if err := Validate(req); err == nil {
		t.Error("expected error for future dob, got nil")
	}
}

func TestValidate_UpdateUserRequest_Valid(t *testing.T) {
	req := UpdateUserRequest{
		Name: "Bob Updated",
		DOB:  "1985-12-25",
	}
	if err := Validate(req); err != nil {
		t.Errorf("expected no error for valid update request, got: %v", err)
	}
}

func TestValidate_NameTooLong(t *testing.T) {
	longName := make([]byte, 256)
	for i := range longName {
		longName[i] = 'a'
	}
	req := CreateUserRequest{
		Name: string(longName),
		DOB:  "1990-01-01",
	}
	if err := Validate(req); err == nil {
		t.Error("expected error for name exceeding 255 chars, got nil")
	}
}

func TestFormatValidationErrors_ContainsFieldNames(t *testing.T) {
	req := CreateUserRequest{
		Name: "",
		DOB:  "",
	}
	err := Validate(req)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	errors := FormatValidationErrors(err)
	if _, ok := errors["name"]; !ok {
		t.Error("expected 'name' field in validation errors")
	}
	if _, ok := errors["dob"]; !ok {
		t.Error("expected 'dob' field in validation errors")
	}
}
