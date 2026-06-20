package processor

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"testing"
)

func TestMatchesFilter(t *testing.T) {
	rec := Record{
		ID:     78,
		Name:   "Alice",
		Role:   "admin",
		Salary: 100000.50,
	}

	tests := []struct {
		name     string
		field    string
		expected string
		want     bool
	}{
		// --- id ---
		{"id match", "id", "78", true},
		{"id mismatch", "id", "92", false},
		{"id invalid format", "id", "smth", false},

		// --- name ---
		{"name match", "name", "Alice", true},
		{"name mismatch", "name", "Alex", false},
		{"name case sensitivity", "name", "alice", false},

		// --- role ---
		{"role match", "role", "admin", true},
		{"role mismatch", "role", "user", false},

		// --- salary ---
		{"salary match", "salary", "100000.50", true},
		{"salary integer", "salary", "100000", false},
		{"salary invalid format", "salary", "smth", false},
		{"salary precision", "salary", "100000.5000000001", true},

		// --- unknown ---
		{"unknown field", "email", "x", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesFilter(rec, tt.field, tt.expected)

			if result != tt.want {
				t.Errorf("matchesFilter(field=%q, expected=%q) = %v, want %v",
					tt.field, tt.expected, result, tt.want)
			}
		})
	}
}

func TestMatchesFilter_FloatPrecision(t *testing.T) {
	rec := Record{Salary: 0.1 + 0.2}

	if !matchesFilter(rec, "salary", "0.3") {
		t.Error("expected 0.1+0.2 to match 0.3 with epsilon tolerance")
	}
}

func TestProcess(t *testing.T) {
	input := `{"id":1,"name":"Alice","role":"admin","salary":100}
{"id":2,"name":"Bob","role":"user","salary":50}
{"id":3,"name":"Carol","role":"admin","salary":120}
`

	tests := []struct {
		name        string
		field       string
		value       string
		maxRecords  int
		wantMatched int
		wantSkipped int
	}{
		{
			name:        "no filter — all matched",
			field:       "",
			value:       "",
			wantMatched: 3,
			wantSkipped: 0,
		},
		{
			name:        "filter by role=admin",
			field:       "role",
			value:       "admin",
			wantMatched: 2,
			wantSkipped: 1,
		},
		{
			name:        "max records limit",
			field:       "",
			value:       "",
			maxRecords:  2,
			wantMatched: 2,
			wantSkipped: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(input)
			out := &bytes.Buffer{}

			stats, err := Process(in, out, tt.field, tt.value, tt.maxRecords)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if stats.Matched != tt.wantMatched {
				t.Errorf("matched = %d, want %d", stats.Matched, tt.wantMatched)
			}
			if stats.Skipped != tt.wantSkipped {
				t.Errorf("skipped = %d, want %d", stats.Skipped, tt.wantSkipped)
			}

			if out.Len() > 0 {
				for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
					var rec Record
					if err = json.Unmarshal([]byte(line), &rec); err != nil {
						t.Errorf("output is not valid JSON Lines: %v", err)
					}
				}
			}
		})
	}
}

func BenchmarkProcess(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < 10000; i++ {
		sb.WriteString(`{"id":`)
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString(`,"name":"user","role":"admin","salary":100}`)
		sb.WriteByte('\n')
	}
	input := sb.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in := strings.NewReader(input)
		out := io.Discard
		_, _ = Process(in, out, "role", "admin", 0)
	}
}

func FuzzMatchesFilter(f *testing.F) {
	f.Add("id", "42")
	f.Add("name", "Alice")
	f.Add("salary", "100.50")
	f.Add("unknown", "value")

	rec := Record{ID: 42, Name: "Alice", Role: "admin", Salary: 100.50}

	f.Fuzz(func(t *testing.T, field, expected string) {
		_ = matchesFilter(rec, field, expected)
	})
}
