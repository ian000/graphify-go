package graph_test

import (
	"testing"

	"github.com/kings2017/graphify-go/internal/graph"
)

func TestNodeIDGeneration(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{"Single simple part", []string{"class", "MyClass"}, "class_myclass"},
		{"With extension", []string{"file", "index.ts"}, "file_index_ts"},
		{"With slashes", []string{"file", "src/components/button.tsx"}, "file_src_components_button_tsx"},
		{"With dots in method", []string{"entity", ".length()"}, "entity_length"},
		{"With underscores", []string{"func", "_internal_helper"}, "func_internal_helper"},
		{"Multiple dirty parts", []string{"_file_", ".test.js", "__magic__"}, "file_test_js_magic"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := graph.GenerateNodeID(tt.parts...)
			if got != tt.expected {
				t.Errorf("GenerateNodeID() = %v, want %v", got, tt.expected)
			}
		})
	}
}
