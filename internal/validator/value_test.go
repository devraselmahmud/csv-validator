package validator

import (
	"testing"
)

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "empty", value: "", want: true},
		{name: "spaces and tab", value: " \t ", want: true},
		{name: "newline", value: "\n", want: true},
		{name: "unicode whitespace", value: "\u00a0", want: true},
		{name: "text", value: "hello", want: false},
		{name: "padded text", value: " hello ", want: false},
		{name: "zero", value: "0", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBlank(tt.value)

			if got != tt.want {
				t.Fatalf("isBlank(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})

	}
}
