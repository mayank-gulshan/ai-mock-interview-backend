package dto

// Mirrors Dto/GeminiResponse.kt exactly, including field names/JSON shape
// expected/produced by the Gemini generateContent API.

type GeminiRequest struct {
	Contents         []GeminiContent   `json:"contents"`
	GenerationConfig *GenerationConfig `json:"generationConfig,omitempty"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GenerationConfig struct {
	ResponseMimeType string         `json:"responseMimeType"`
	ResponseSchema   ResponseSchema `json:"responseSchema"`
}

type ResponseSchema struct {
	Type  string     `json:"type"`
	Items SchemaItem `json:"items"`
}

type SchemaItem struct {
	Type       string                    `json:"type"`
	Properties map[string]SchemaProperty `json:"properties"`
	Required   []string                  `json:"required"`
}

type SchemaProperty struct {
	Type string `json:"type"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content ContentResponse `json:"content"`
}

type ContentResponse struct {
	Parts []PartResponse `json:"parts"`
}

type PartResponse struct {
	Text string `json:"text"`
}

type QuestionAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// NewGenerationConfig builds the same defaults as the Kotlin data class
// GenerationConfig(responseMimeType: String = "application/json", ...).
func NewGenerationConfig(schema ResponseSchema) *GenerationConfig {
	return &GenerationConfig{
		ResponseMimeType: "application/json",
		ResponseSchema:   schema,
	}
}

// NewResponseSchema mirrors ResponseSchema(type: String = "ARRAY", items).
func NewResponseSchema(items SchemaItem) ResponseSchema {
	return ResponseSchema{Type: "ARRAY", Items: items}
}

// NewSchemaItem mirrors SchemaItem(type: String = "OBJECT", properties, required).
func NewSchemaItem(properties map[string]SchemaProperty, required []string) SchemaItem {
	return SchemaItem{Type: "OBJECT", Properties: properties, Required: required}
}
