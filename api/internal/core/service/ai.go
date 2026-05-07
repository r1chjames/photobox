package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Images []string `json:"images,omitempty"`
	Stream bool   `json:"stream"`
	Format string `json:"format,omitempty"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type OllamaClient struct {
	host   string
	model  string
	client *http.Client
}

// NewOllamaClient creates a new Ollama AI service client
func NewOllamaClient(config appconfig.AppConfig) port.AIService {
	return &OllamaClient{
		host:  strings.TrimRight(config.OllamaHost, "/"),
		model: config.OllamaModel,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (o *OllamaClient) AnalyzeImage(imagePath string) (*port.ImageAnalysis, error) {
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	prompt := `Analyze this image. Return ONLY valid JSON in this format:
{"caption": "...", "tags": ["..."], "objects": ["..."], "is_nsfw": false, "is_portrait": false}`

	reqBody := ollamaGenerateRequest{
		Model:  o.model,
		Prompt: prompt,
		Images: []string{base64Image},
		Stream: false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", o.host)
	resp, err := o.client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp ollamaGenerateResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ollama response: %w", err)
	}

	// Clean up the response: extract JSON if it's wrapped in markdown code blocks
	cleanResponse := strings.TrimSpace(ollamaResp.Response)
	cleanResponse = strings.TrimPrefix(cleanResponse, "```json")
	cleanResponse = strings.TrimPrefix(cleanResponse, "```")
	cleanResponse = strings.TrimSuffix(cleanResponse, "```")
	cleanResponse = strings.TrimSpace(cleanResponse)

	var analysis port.ImageAnalysis
	if err := json.Unmarshal([]byte(cleanResponse), &analysis); err != nil {
		slog.Warn("Failed to parse AI analysis JSON", "response", cleanResponse, "error", err)
		// Return a best-effort result with empty fields
		return &port.ImageAnalysis{
			Caption: "",
			Tags:    []string{},
			Objects: []string{},
		}, nil
	}

	return &analysis, nil
}
