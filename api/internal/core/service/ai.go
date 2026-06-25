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

	prompt := `Look at this image carefully. What do you see? Write a JSON object with:
- "caption": describe what is shown in one sentence
- "tags": list relevant search keywords (specific, not generic)
- "objects": list visible things by name (text strings, not numbers)
- "is_nsfw": true if inappropriate, false otherwise
- "is_portrait": true if a face is the main subject, false otherwise`

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

	var resp *http.Response
	var lastErr error
	for attempt := 1; attempt <= 2; attempt++ {
		startTime := time.Now()
		resp, lastErr = o.client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
		elapsed := time.Since(startTime)
		if lastErr == nil {
			slog.Info("Ollama response received", "status", resp.StatusCode, "elapsed", elapsed, "attempt", attempt)
			break
		}
		slog.Warn("Ollama request failed, retrying", "elapsed", elapsed, "error", lastErr, "attempt", attempt)
		time.Sleep(5 * time.Second)
	}
	if lastErr != nil {
		slog.Error("Ollama request failed after retries", "error", lastErr)
		return nil, fmt.Errorf("failed to call ollama: %w", lastErr)
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

	// Use a flexible intermediate struct to handle moondream2 quirks:
	// - objects may be bounding box coords (float arrays) or string lists
	// - is_nsfw / is_portrait may be 0/1 (int) instead of false/true (bool)
	// - tags may be missing entirely
	type rawAnalysis struct {
		Caption    string        `json:"caption"`
		Tags       []string      `json:"tags"`
		Objects    []interface{} `json:"objects"`
		IsNSFW     interface{}   `json:"is_nsfw"`
		IsPortrait interface{}   `json:"is_portrait"`
	}
	var raw rawAnalysis
	if err := json.Unmarshal([]byte(cleanResponse), &raw); err != nil {
		slog.Warn("Failed to parse AI analysis JSON", "response", cleanResponse, "error", err)
		return nil, fmt.Errorf("failed to parse AI response as JSON: %w", err)
	}

	// Convert objects to strings, dropping any non-string values
	objects := make([]string, 0, len(raw.Objects))
	for _, obj := range raw.Objects {
		if s, ok := obj.(string); ok {
			objects = append(objects, s)
		}
	}

	// Convert bool-ish values to actual booleans (handles 0/1, false/true, strings)
	toBool := func(v interface{}) bool {
		switch val := v.(type) {
		case bool:
			return val
		case float64:
			return val != 0
		case string:
			return val == "true" || val == "1"
		default:
			return false
		}
	}

	return &port.ImageAnalysis{
		Caption:    raw.Caption,
		Tags:       raw.Tags,
		Objects:    objects,
		IsNSFW:     toBool(raw.IsNSFW),
		IsPortrait: toBool(raw.IsPortrait),
	}, nil
}
