package parser

import (
	"os"
	"testing"
)

func TestNewMermaidImpl(t *testing.T) {
	tests := []struct {
		name              string
		content           string
		expectedNodes     []Node
		expectedRelations []Relation
		expectError       bool
	}{
		{
			name: "Basic Graph",
			content: `graph TD
    A["Start"]
    B["End"]
    A --> B`,
			expectedNodes: []Node{
				{ID: "A", Name: "\"Start\"", Comment: ""},
				{ID: "B", Name: "\"End\"", Comment: ""},
			},
			expectedRelations: []Relation{
				{From: "A", To: "B"},
			},
			expectError: false,
		},
		{
			name: "Graph with Comments",
			content: `graph TD
    %% comment for A
    A["Node A"]
    B["Node B"]
    A --> B`,
			expectedNodes: []Node{
				{ID: "A", Name: "\"Node A\"", Comment: "comment for A"},
				{ID: "B", Name: "\"Node B\"", Comment: ""},
			},
			expectedRelations: []Relation{
				{From: "A", To: "B"},
			},
			expectError: false,
		},
		{
			name: "Complex Graph with Subgraph",
			content: `graph TD
    subgraph SG1
        A["Node A"]
        B["Node B"]
        A --> B
    end
    C["Node C"]
    B --> C`,
			expectedNodes: []Node{
				{ID: "A", Name: "\"Node A\"", Comment: ""},
				{ID: "B", Name: "\"Node B\"", Comment: ""},
				{ID: "C", Name: "\"Node C\"", Comment: ""},
			},
			expectedRelations: []Relation{
				{From: "A", To: "B"},
				{From: "B", To: "C"},
			},
			expectError: false,
		},
		{
			name:        "Invalid File Content",
			content:     `invalid mermaid content`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tmpFile, err := os.CreateTemp("", "mermaid-*.mmd")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.content); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			if err := tmpFile.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}

			// Run the function under test
			impl, err := NewMermaidImpl(tmpFile.Name())

			// Check error expectations
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("NewMermaidImpl returned unexpected error: %v", err)
			}

			// Verify Nodes
			if len(impl.Nodes) != len(tt.expectedNodes) {
				t.Errorf("Expected %d nodes, got %d", len(tt.expectedNodes), len(impl.Nodes))
			}

			// Helper map for finding expected nodes easily, though since we use structs we could just loop.
			// Let's iterate over found nodes and check if they exist in expected string.
			// Or better, convert both to maps by ID for easier comparison since order is not guaranteed.
			actualNodeMap := make(map[string]Node)
			for _, n := range impl.Nodes {
				actualNodeMap[n.ID] = n
			}

			for _, expected := range tt.expectedNodes {
				actual, exists := actualNodeMap[expected.ID]
				if !exists {
					t.Errorf("Expected node ID %s not found", expected.ID)
					continue
				}
				if actual.Name != expected.Name {
					t.Errorf("Node %s: expected Name '%s', got '%s'", expected.ID, expected.Name, actual.Name)
				}
				if actual.Comment != expected.Comment {
					t.Errorf("Node %s: expected Comment '%s', got '%s'", expected.ID, expected.Comment, actual.Comment)
				}
			}

			// Verify Relations
			if len(impl.Relations) != len(tt.expectedRelations) {
				t.Errorf("Expected %d relations, got %d", len(tt.expectedRelations), len(impl.Relations))
			}

			for _, expectedRel := range tt.expectedRelations {
				found := false
				for _, rel := range impl.Relations {
					if rel.From == expectedRel.From && rel.To == expectedRel.To {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected relation %s -> %s not found", expectedRel.From, expectedRel.To)
				}
			}
		})
	}
}

func TestNewMermaidImpl_FileNotFound(t *testing.T) {
	_, err := NewMermaidImpl("non_existent_file.mmd")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
