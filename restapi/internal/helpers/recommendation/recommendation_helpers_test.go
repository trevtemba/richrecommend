package recommendationhelpers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/trevtemba/richrecommend/internal/models"
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
				"required":             []string{"conditioners", "shampoos"},
				"additionalProperties": false,
				"DoNotReference":       true,
			},
		},
		{
			name:       "empty categories",
			categories: []string{},
			want: map[string]any{
				"type":                 "object",
				"properties":           map[string]any{},
				"required":             []string{},
				"additionalProperties": false,
				"DoNotReference":       true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSchema(tt.categories)
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
		want    models.RecommendationResponse
		wantErr bool
	}{
		{
			name: "valid JSON",
			params: testParams{
				content:    `{"conditioners": ["A", "B"], "shampoos": ["C"]}`,
				categories: categories,
			},
			want: models.RecommendationResponse{
				Recommendation: map[string][]string{
					"conditioners": {"A", "B"},
					"shampoos":     {"C"},
				},
				ItemCount: 3,
			},
			wantErr: false,
		},
		{
			name: "missing category",
			params: testParams{
				content:    `{"conditioners": ["A", "B"]}`,
				categories: categories,
			},
			want: models.RecommendationResponse{
				Recommendation: map[string][]string{
					"conditioners": {"A", "B"},
					"shampoos":     nil,
				},
				ItemCount: 2,
			},
			wantErr: false,
		},
		{
			name: "malformed JSON",
			params: testParams{
				content:    `{"conditioners": [}`,
				categories: categories,
			},
			want:    models.RecommendationResponse{},
			wantErr: true,
		},
		{
			name: "wrong type values",
			params: testParams{
				content:    `{"conditioners": [1, 2], "shampoos": [3]}`,
				categories: categories,
			},
			want: models.RecommendationResponse{
				Recommendation: map[string][]string{
					"conditioners": nil,
					"shampoos":     nil,
				},
				ItemCount: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseChatResponse(tt.params.content, tt.params.categories)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseChatResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ParseChatResponse() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestGenerateSystemMessage(t *testing.T) {
	type testParams struct {
		systemPrompt      models.SystemPrompt
		contextSchemaName string
		categories        []string
		recsPerCategory   int
	}
	tests := []struct {
		name   string
		params testParams
		want   string
	}{
		{
			name: "basic system prompt with domain",
			params: testParams{
				systemPrompt: models.SystemPrompt{
					Role:      "professional hair care expert",
					Clientele: "black clients with type 4 hair",
					Domain:    "hair care",
				},
				contextSchemaName: "hair_profile",
				categories:        []string{"conditioners", "shampoos"},
				recsPerCategory:   3,
			},

			want: "You are a professional hair care expert for black clients with type 4 hair. " +
				"in hair care Given a hair_profile, recommend up to 3 products per category (conditioners, shampoos). " +
				"Return the response in structured JSON format.",
		},
		{
			name: "system prompt without domain",
			params: testParams{
				systemPrompt: models.SystemPrompt{
					Role:      "nutrition advisor",
					Clientele: "athletes",
					Domain:    "",
				},
				contextSchemaName: "nutrition_profile",
				categories:        []string{"supplements"},
				recsPerCategory:   2,
			},
			want: "You are a nutrition advisor for athletes. " +
				"Given a nutrition_profile, recommend up to 2 products per category (supplements). " +
				"Return the response in structured JSON format.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSystemMessage(tt.params.systemPrompt, tt.params.contextSchemaName, tt.params.categories, tt.params.recsPerCategory)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GenerateSystemMessage mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGenerateUserMessage(t *testing.T) {
	tests := []struct {
		name          string
		contextSchema models.ContextSchema
		want          string
		wantErr       bool
	}{
		{
			name: "valid schema",
			contextSchema: models.ContextSchema{
				Name: "hair_profile",
				Content: map[string]any{
					"curl_type":       "4a",
					"desired_outcome": "moisturize",
					"porosity":        "high",
				},
			},
			want:    "hair_profile: {\ncurl_type: 4a,\ndesired_outcome: moisturize,\nporosity: high\n}\nRecommend products.",
			wantErr: false,
		},
		{
			name: "invalid non-string value",
			contextSchema: models.ContextSchema{
				Name: "hair_profile",
				Content: map[string]any{
					"porosity": 123,
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty schema",
			contextSchema: models.ContextSchema{
				Name:    "empty_profile",
				Content: map[string]any{},
			},
			want:    "empty_profile: {\n}\nRecommend products.",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateUserMessage(tt.contextSchema)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateUserMessage error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GenerateUserMessage mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
