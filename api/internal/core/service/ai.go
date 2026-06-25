package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	json "github.com/goccy/go-json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/disintegration/imaging"
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
			Timeout: 600 * time.Second,
		},
	}
}

func (o *OllamaClient) AnalyzeImage(imagePath string) (*port.ImageAnalysis, error) {
	// Open and resize image to reduce payload size and inference time
	src, err := imaging.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}

	// Resize to max 512px on longest side for faster inference
	resized := imaging.Fit(src, 512, 512, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, resized, imaging.JPEG); err != nil {
		return nil, fmt.Errorf("failed to encode resized image: %w", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())

	prompt := `Describe this image. Output a JSON object with these fields: 
{"caption": "concise 1-sentence description", "tags": ["keyword1", "keyword2"], "objects": ["visible", "physical", "objects"], "is_nsfw": false, "is_portrait": false}
- tags: searchable keywords, avoid generic words like "image" or "photo"
- objects: distinct physical items visible
- is_nsfw: true only for nudity/sexual content
- is_portrait: true if a person's face is the main subject`

	reqBody := ollamaGenerateRequest{
		Model:  o.model,
		Prompt: prompt,
		Images: []string{base64Image},
		Stream: false,
		Format: "json",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", o.host)
	slog.Info("Calling Ollama for image analysis", "model", o.model, "imageSize", len(base64Image))

	startTime := time.Now()
	resp, err := o.client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	elapsed := time.Since(startTime)
	if err != nil {
		slog.Error("Ollama request failed", "elapsed", elapsed, "error", err)
		return nil, fmt.Errorf("failed to call ollama: %w", err)
	}
	defer resp.Body.Close()

	slog.Info("Ollama response received", "status", resp.StatusCode, "elapsed", elapsed)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	slog.Info("Ollama raw response", "body", string(body))

	var ollamaResp ollamaGenerateResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ollama response: %w", err)
	}

	slog.Info("Ollama extracted response text", "text", ollamaResp.Response)

	// Clean up the response: extract JSON if it's wrapped in markdown code blocks
	cleanResponse := strings.TrimSpace(ollamaResp.Response)
	cleanResponse = strings.TrimPrefix(cleanResponse, "```json")
	cleanResponse = strings.TrimPrefix(cleanResponse, "```")
	cleanResponse = strings.TrimSuffix(cleanResponse, "```")
	cleanResponse = strings.TrimSpace(cleanResponse)

	var analysis port.ImageAnalysis
	if err := json.Unmarshal([]byte(cleanResponse), &analysis); err != nil {
		slog.Warn("Failed to parse AI analysis JSON", "response", cleanResponse, "error", err)
		return nil, fmt.Errorf("failed to parse AI response as JSON: %w", err)
	}

	return &analysis, nil
}
