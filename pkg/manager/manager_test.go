package manager

import (
	"reflect"
	"testing"
)

func TestNewModel_Valid(t *testing.T) {
	testCases := []struct {
		name     string
		modelStr string
		base     string
		expected *Model
	}{
		{
			name:     "model without tag",
			modelStr: "deepseek-r1",
			base:     "base",
			expected: &Model{
				Name: "deepseek-r1",
				Tag:  "latest",
				Base: "base",
			},
		},
		{
			name:     "model with tag",
			modelStr: "deepseek-r1:14b",
			base:     "base",
			expected: &Model{
				Name: "deepseek-r1",
				Tag:  "14b",
				Base: "base",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewModel(tc.modelStr, tc.base, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestNewModel_Invalid(t *testing.T) {
	invalidCases := []struct {
		name     string
		modelStr string
		base     string
	}{
		{
			name:     "empty model string",
			modelStr: "",
			base:     "base",
		},
		{
			name:     "colon at the end",
			modelStr: "deepseek-r1:",
			base:     "base",
		},
		{
			name:     "colon at the beginning",
			modelStr: ":14b",
			base:     "base",
		},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewModel(tc.modelStr, tc.base, nil)
			if err == nil {
				t.Errorf("expected error for invalid model string %q, but got none", tc.modelStr)
			}
		})
	}
}
