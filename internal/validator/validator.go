package validator

import (
	"encoding/csv"
	"fmt"
	"io"
	"slices"
)

type Options struct {
	RequiredColumns []string
	IDColumn        string
}

type Finding struct {
	Record  int
	Column  string
	Message string
}

func Validate(input io.Reader, options Options) ([]Finding, error) {
	if isBlank(options.IDColumn) {
		return nil, fmt.Errorf("ID column must not be blank")
	}

	for _, column := range options.RequiredColumns {
		if isBlank(column) {
			return nil, fmt.Errorf("required column name must not be blank")
		}
	}

	if !slices.Contains(options.RequiredColumns, options.IDColumn) {
		return nil, fmt.Errorf(
			"ID column %q must belong to required columns",
			options.IDColumn,
		)
	}

	reader := csv.NewReader(input)

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read CSV header: %w", err)
	}

	var findings []Finding

	seenHeaders := make(map[string]struct{})

	for _, column := range headers {
		if isBlank(column) {
			findings = append(findings, Finding{
				Record:  1,
				Column:  column,
				Message: "column name blank",
			})
			continue
		}

		if _, exists := seenHeaders[column]; exists {
			findings = append(findings, Finding{
				Record:  1,
				Column:  column,
				Message: "duplicate column name",
			})
		}

		seenHeaders[column] = struct{}{}
	}

	for _, column := range missingColumns(headers, options.RequiredColumns) {
		findings = append(findings, Finding{
			Record:  1,
			Column:  column,
			Message: "required column missing",
		})
	}

	if len(findings) > 0 {
		return findings, nil
	}

	indices := make(map[string]int, len(headers))
	for index, column := range headers {
		indices[column] = index
	}

	tracker := idTracker{}
	idIndex := indices[options.IDColumn]

	for recordNumber := 2; ; recordNumber++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return findings, fmt.Errorf("read CSV record %d: %w", recordNumber, err)
		}

		for _, column := range options.RequiredColumns {
			if isBlank(record[indices[column]]) {
				findings = append(findings, Finding{
					Record:  recordNumber,
					Column:  column,
					Message: "required value missing",
				})
			}
		}

		if tracker.isDuplicate(record[idIndex]) {
			findings = append(findings, Finding{
				Record:  recordNumber,
				Column:  options.IDColumn,
				Message: "duplicate ID",
			})
		}
	}

	return findings, nil
}
