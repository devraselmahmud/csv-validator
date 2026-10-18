package validator

import (
	"slices"
	"testing"
)

func TestMissingColumnsReportsAbsentColumn(t *testing.T) {
	headers := []string{"id", "name"}
	required := []string{"id", "name", "email"}

	got := missingColumns(headers, required)
	want := []string{"email"}

	if !slices.Equal(got, want) {
		t.Fatalf("missingColumns() = %v, want %v", got, want)
	}

}

func TestMissingColumnsReturnsEmptyWhenAllRequiredColumnsExist(t *testing.T) {
	headers := []string{"email", "id", "name", "notes"}
	required := []string{"id", "name", "email"}

	got := missingColumns(headers, required)

	if len(got) != 0 {
		t.Fatalf("missingColumns() = %v, want no missing columns", got)
	}
}

func TestMissingColumnsDoesNotTrimHeaders(t *testing.T) {
	headers := []string{"id", " email "}
	required := []string{"id", "email"}

	got := missingColumns(headers, required)
	want := []string{"email"}

	if !slices.Equal(got, want) {
		t.Fatalf("missingColumns() = %v, want %v", got, want)
	}
}
