package validator

import "slices"

func missingColumns(header, required []string) []string {
	var missing []string

	for _, column := range required {
		if !slices.Contains(header, column) {
			missing = append(missing, column)

		}

	}
	return missing
}
