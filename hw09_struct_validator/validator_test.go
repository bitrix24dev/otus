package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

type User struct {
	ID     string `json:"id" validate:"len:36"`
	Name   string
	Age    int             `validate:"min:18|max:50"`
	Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	Role   UserRole        `validate:"in:admin,stuff"`
	Phones []string        `validate:"len:11"`
	meta   json.RawMessage //nolint:unused
}

type App struct {
	Version string `validate:"len:5"`
}

type Response struct {
	Code int    `validate:"in:200,404,500"`
	Body string `json:"omitempty"`
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John",
				Age:    25,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "123",
				Name:   "John",
				Age:    17,
				Email:  "johnexample.com",
				Role:   "user",
				Phones: []string{"123"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: errors.New("length of ID must be 36")},
				{Field: "Age", Err: errors.New("Age must be at least 18")},
				{Field: "Email", Err: errors.New("Email does not match regexp ^\\w+@\\w+\\.\\w+$")},
				{Field: "Role", Err: errors.New("Role is not in [admin stuff]")},
				{Field: "Phones", Err: errors.New("length of Phones must be 11")},
			},
		},
		{
			in: App{
				Version: "1.0", // Длина 4, ожидается ошибка, так как длина должна быть 5
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: errors.New("length of Version must be 5")},
			},
		},
		{
			in: Response{
				Code: 201,
				Body: "OK",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: errors.New("Code is not in [200 404 500]")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			if tt.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
