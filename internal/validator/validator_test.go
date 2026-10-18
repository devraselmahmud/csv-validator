package validator

import (
	"encoding/csv"
	"errors"
	"io"
	"slices"
	"strings"
	"testing"
)

func TestValidateReportsMissingRequiredColumn(t *testing.T) {
	input := strings.NewReader("id,name\n42,Ada\n")
	options := Options{
		RequiredColumns: []string{"id", "name", "email"},
		IDColumn:        "id",
	}

	findings, err := Validate(input, options)
	if err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("Validate() returned %d findings, want 1", len(findings))
	}

	got := findings[0]
	if got.Record != 1 || got.Column != "email" || got.Message != "required column missing" {
		t.Fatalf("Validate() finding = %+v, want missing email column at record 1", got)
	}
}

func TestValidateReportsBlankRequiredValue(t *testing.T) {
	input := strings.NewReader("id,name,notes\n42, \t ,\n")
	options := Options{
		RequiredColumns: []string{"id", "name"},
		IDColumn:        "id",
	}

	findings, err := Validate(input, options)
	if err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("Validate() returned %d findings, want 1", len(findings))
	}

	got := findings[0]
	if got.Record != 2 || got.Column != "name" || got.Message != "required value missing" {
		t.Fatalf("Validate() finding = %+v, want missing name value at record 2", got)
	}
}

func TestValidateReportsDuplicateID(t *testing.T) {
	input := strings.NewReader("id,name\n42,Ada\n 42 ,Grace\n")
	options := Options{
		RequiredColumns: []string{"id", "name"},
		IDColumn:        "id",
	}

	findings, err := Validate(input, options)
	if err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("Validate() returned %d findings, want 1", len(findings))
	}

	got := findings[0]
	if got.Record != 3 || got.Column != "id" || got.Message != "duplicate ID" {
		t.Fatalf("Validate() finding = %+v, want duplicate ID at record 3", got)
	}
}

func TestValidateRejectsIDColumnNotInRequiredColumns(t *testing.T) {
	input := strings.NewReader("id,name\n42,Ada\n")
	options := Options{
		RequiredColumns: []string{"name"},
		IDColumn:        "id",
	}

	findings, err := Validate(input, options)
	if err == nil {
		t.Fatal("Validate() returned no error for ID column outside required columns")
	}

	if len(findings) != 0 {
		t.Fatalf("Validate() returned findings for invalid options: %+v", findings)
	}
}

func TestValidateRejectsBlankRequiredColumn(t *testing.T) {
	input := strings.NewReader("id, \n42,Ada\n")
	options := Options{
		RequiredColumns: []string{"id", " "},
		IDColumn:        "id",
	}

	_, err := Validate(input, options)
	if err == nil {
		t.Fatal("Validate() returned no error for blank required column name")
	}
}

func TestValidateRejectsBlankIDColumn(t *testing.T) {
	input := strings.NewReader(" ,name\n42,Ada\n")
	options := Options{
		RequiredColumns: []string{" ", "name"},
		IDColumn:        " ",
	}

	_, err := Validate(input, options)
	if err == nil {
		t.Fatal("Validate() returned no error for blank ID column name")
	}
}

func TestValidateReportsBlankHeader(t *testing.T) {
	input := strings.NewReader("id, \n42,Ada\n")
	options := Options{
		RequiredColumns: []string{"id"},
		IDColumn:        "id",
	}

	findings, err := Validate(input, options)
	if err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("Validate() returned %d findings, want 1", len(findings))
	}

	got := findings[0]
	if got.Record != 1 || got.Column != " " || got.Message != "column name blank" {
		t.Fatalf("Validate() finding = %+v, want blank header at record 1", got)
	}
}

func TestValidateReportsDuplicateHeader(t *testing.T) {
	input := strings.NewReader("id,name,name\n42,Ada,Grace\n")
	options := Options{
		RequiredColumns: []string{"id", "name"},
		IDColumn:        "id",
	}

	findings, err := Validate(input, options)
	if err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("Validate() returned %d findings, want 1", len(findings))
	}

	got := findings[0]
	if got.Record != 1 || got.Column != "name" || got.Message != "duplicate column name" {
		t.Fatalf("Validate() finding = %+v, want duplicate name header at record 1", got)
	}
}

func TestValidateReturnsInputErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{name: "empty input", input: "", want: io.EOF},
		{
			name:  "unclosed quoted field",
			input: "id,name\n42,\"Ada\n",
			want:  csv.ErrQuote,
		},
		{
			name:  "too few fields",
			input: "id,name\n42\n",
			want:  csv.ErrFieldCount,
		},
		{
			name:  "too many fields",
			input: "id,name\n42,Ada,extra\n",
			want:  csv.ErrFieldCount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := Options{
				RequiredColumns: []string{"id", "name"},
				IDColumn:        "id",
			}

			_, err := Validate(strings.NewReader(tt.input), options)
			if !errors.Is(err, tt.want) {
				t.Fatalf("Validate() error = %v, want wrapped %v", err, tt.want)
			}
		})
	}
}

func TestValidateAcceptsValidCSV(t *testing.T) {
	input := strings.NewReader(
		"name,notes,id\nAda,,42\nGrace,\"hello, world\",43\n",
	)
	options := Options{
		RequiredColumns: []string{"id", "name"},
		IDColumn:        "id",
	}

	findings, err := Validate(input, options)
	if err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	if len(findings) != 0 {
		t.Fatalf("Validate() returned findings for valid CSV: %+v", findings)
	}
}

func TestValidateCountsRecordsNotPhysicalLines(t *testing.T) {
	input := strings.NewReader(
		"id,name\n42,\"Ada\nLovelace\"\n43, \n",
	)
	options := Options{
		RequiredColumns: []string{"id", "name"},
		IDColumn:        "id",
	}

	findings, err := Validate(input, options)
	if err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("Validate() returned %d findings, want 1", len(findings))
	}

	want := Finding{
		Record:  3,
		Column:  "name",
		Message: "required value missing",
	}
	if findings[0] != want {
		t.Fatalf("Validate() finding = %+v, want %+v", findings[0], want)
	}
}

func TestValidateCollectsMultipleFindings(t *testing.T) {
	input := strings.NewReader(
		"id,name\n42,Ada\n 42 , \n,Grace\n \t ,Lin\n",
	)
	options := Options{
		RequiredColumns: []string{"id", "name"},
		IDColumn:        "id",
	}

	got, err := Validate(input, options)
	if err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	want := []Finding{
		{Record: 3, Column: "name", Message: "required value missing"},
		{Record: 3, Column: "id", Message: "duplicate ID"},
		{Record: 4, Column: "id", Message: "required value missing"},
		{Record: 5, Column: "id", Message: "required value missing"},
	}

	if !slices.Equal(got, want) {
		t.Fatalf("Validate() findings = %+v, want %+v", got, want)
	}
}
