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
	Model   string         `json:"model"`
	Prompt  string         `json:"prompt"`
	Images  []string       `json:"images,omitempty"`
	Stream  bool           `json:"stream"`
	Format  string         `json:"format,omitempty"`
	Options map[string]any `json:"options,omitempty"`
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

	// Resize to max 320px — fewer image tokens (400 vs 729 at 512px) = faster inference
	resized := imaging.Fit(src, 320, 320, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, resized, imaging.JPEG); err != nil {
		return nil, fmt.Errorf("failed to encode resized image: %w", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Simple natural-language prompt — moondream2 works best with Q&A, not JSON
	prompt := `Describe this image in one sentence.`

	reqBody := ollamaGenerateRequest{
		Model:  o.model,
		Prompt: prompt,
		Images: []string{base64Image},
		Stream: false,
		Options: map[string]any{
			"num_ctx":    2048,
			"num_predict": 64, // one sentence is ~20 tokens
		},
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

	// Use the raw text directly as the caption (moondream2 returns natural language)
	return &port.ImageAnalysis{
		Caption:    cleanResponse,
		Tags:       nil,
		Objects:    nil,
		IsNSFW:     false,
		IsPortrait: false,
	}, nil
}
