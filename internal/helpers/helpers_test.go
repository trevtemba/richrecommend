package helpers_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/trevtemba/richrecommend/internal/helpers"
)

func TestGenerateSchema(t *testing.T) {
	tests := []struct {
		name       string
		categories []string
		want       map[string]any
	}{
		{
			name:       "with categories",
			categories: []string{"conditioners", "shampoos"},
			want: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"conditioners": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "list of recommended conditioners",
					},
					"shampoos": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "list of recommended shampoos",
					},
				},
				"required": []string{"conditioners", "shampoos"},
			},
		},
		{
			name:       "empty categories",
			categories: []string{},
			want: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
				"required":   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.GenerateSchema(tt.categories)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GenerateSchema mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseChatResponse(t *testing.T) {
	categories := []string{"conditioners", "shampoos"}

	type testParams struct {
		content    string
		categories []string
	}
	tests := []struct {
		name    string
		params  testParams
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "valid JSON",
			params: testParams{
				content:    `{"conditioners": ["A", "B"], "shampoos": ["C"]}`,
				categories: []string{"conditioners, shampoos"},
			},
			want: map[string][]string{
				"conditioners": {"A", "B"},
				"shampoos":     {"C"},
			},
			wantErr: false,
		},
		{
			name: "missing category",
			params: testParams{
				content:    `{"conditioners": ["A", "B"]}`,
				categories: []string{"conditioners, shampoos"},
			},
			want: map[string][]string{
				"conditioners": {"A", "B"},
				"shampoos":     nil,
			},
			wantErr: false,
		},
		{
			name: "malformed JSON",
			params: testParams{
				content:    `{"conditioners": [}`,
				categories: []string{"conditioners, shampoos"},
			},
			want:    map[string][]string{},
			wantErr: true,
		},
		{
			name: "wrong type values",
			params: testParams{
				content:    `{"conditioners": [1, 2], "shampoos": [3]}`,
				categories: []string{"conditioners, shampoos"},
			},
			want: map[string][]string{
				"conditioners": nil,
				"shampoos":     nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := helpers.ParseChatResponse(tt.params.content, categories)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseChatResponse error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ParseChatResponse mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
