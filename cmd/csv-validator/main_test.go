package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name, csv, required, id string
		want                    int
		message                 string
	}{
		{"valid", "id,name,notes\n42,Ada,\n", "id,name", "id", 0, "CSV valid"},
		{"duplicate", "id,name\n42,Ada\n 42 ,Grace\n", "id,name", "id", 1, "duplicate ID"},
		{"blank", "id,name\n42, \n", "id,name", "id", 1, "required value missing"},
		{"missing header", "id\n42\n", "id,name", "id", 1, "required column missing"},
		{"malformed", "id,name\n42\n", "id,name", "id", 2, "read CSV record"},
		{"invalid id", "id,name\n42,Ada\n", "name", "id", 2, "must belong"},
		{"duplicate required", "id,name\n42,Ada\n", "id,id", "id", 2, "duplicate required column"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "input.csv")
			if err := os.WriteFile(path, []byte(tt.csv), 0600); err != nil {
				t.Fatal(err)
			}
			var out, stderr bytes.Buffer
			got := run([]string{"--file", path, "--required", tt.required, "--id", tt.id}, &out, &stderr)
			if got != tt.want {
				t.Fatalf("code = %d, want %d; stderr: %s", got, tt.want, &stderr)
			}
			if !strings.Contains(out.String()+stderr.String(), tt.message) {
				t.Fatalf("missing %q; output: %s%s", tt.message, &out, &stderr)
			}
		})
	}
}

func TestRunArguments(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		want int
	}{
		{"missing", nil, 2}, {"help", []string{"--help"}, 0},
		{"unknown", []string{"--unknown"}, 2},
		{"positional", []string{"extra"}, 2},
		{"missing file", []string{"--file", filepath.Join(t.TempDir(), "absent.csv"), "--required", "id", "--id", "id"}, 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var out, stderr bytes.Buffer
			if got := run(tt.args, &out, &stderr); got != tt.want {
				t.Fatalf("code = %d, want %d", got, tt.want)
			}
		})
	}
}
