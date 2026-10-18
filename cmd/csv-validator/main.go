package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/devraselmahmud/csv-validator/internal/validator"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("csv-validator", flag.ContinueOnError)
	flags.SetOutput(stderr)
	file := flags.String("file", "", "CSV file path")
	required := flags.String("required", "", "Comma-separated required column names (exact matching)")
	id := flags.String("id", "", "ID column; must be a required column")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 || strings.TrimSpace(*file) == "" || strings.TrimSpace(*required) == "" || strings.TrimSpace(*id) == "" {
		fmt.Fprintln(stderr, "--file, --required and --id are required; positional arguments are not supported")
		flags.PrintDefaults()
		return 2
	}
	columns := strings.Split(*required, ",")
	seen := make(map[string]struct{})
	for _, column := range columns {
		if strings.TrimSpace(column) == "" {
			fmt.Fprintln(stderr, "required column name must not be blank")
			return 2
		}
		if _, exists := seen[column]; exists {
			fmt.Fprintf(stderr, "duplicate required column %q\n", column)
			return 2
		}
		seen[column] = struct{}{}
	}
	input, err := os.Open(*file)
	if err != nil {
		fmt.Fprintf(stderr, "open CSV: %v\n", err)
		return 2
	}
	defer input.Close()
	findings, err := validator.Validate(input, validator.Options{RequiredColumns: columns, IDColumn: *id})
	for _, finding := range findings {
		if _, writeErr := fmt.Fprintf(stdout, "record %d, column %q: %s\n", finding.Record, finding.Column, finding.Message); writeErr != nil {
			fmt.Fprintf(stderr, "write report: %v\n", writeErr)
			return 2
		}
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if len(findings) > 0 {
		return 1
	}
	if _, err := fmt.Fprintln(stdout, "CSV valid"); err != nil {
		fmt.Fprintf(stderr, "write report: %v\n", err)
		return 2
	}
	return 0
}
