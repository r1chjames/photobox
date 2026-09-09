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
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

type faceDetectRequest struct {
	Image    string  `json:"image"`
	MinScore float64 `json:"min_score"`
}

type faceEngineResponse struct {
	Faces []struct {
		Box       [4]float64 `json:"box"`
		Score     float64    `json:"score"`
		Embedding []float32  `json:"embedding"`
	} `json:"faces"`
}

// FaceEngineClient is an HTTP client for the local face-engine microservice.
type FaceEngineClient struct {
	baseURL string
	client  *http.Client
}

// NewFaceEngineClient builds a client for the face engine at config.FaceEngineURL.
func NewFaceEngineClient(config appconfig.AppConfig) port.FaceEngine {
	return &FaceEngineClient{
		baseURL: strings.TrimRight(config.FaceEngineURL, "/"),
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *FaceEngineClient) Detect(imagePath string) ([]port.FaceEngineResult, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	reqBody, err := json.Marshal(faceDetectRequest{
		Image:    base64.StdEncoding.EncodeToString(data),
		MinScore: 0.5,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/detect", c.baseURL)
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		slog.Warn("Face engine request failed", "url", url, "error", err)
		return nil, fmt.Errorf("%w: %v", domain.ErrFaceEngineUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Warn("Face engine returned non-200", "status", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("%w: engine status %d", domain.ErrFaceEngineUnavailable, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read face engine response: %w", err)
	}

	var parsed faceEngineResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse face engine response: %w", err)
	}

	results := make([]port.FaceEngineResult, 0, len(parsed.Faces))
	for _, f := range parsed.Faces {
		if len(f.Embedding) != 128 {
			slog.Warn("Face engine returned unexpected embedding dim", "dim", len(f.Embedding))
			continue
		}
		var emb [128]float32
		copy(emb[:], f.Embedding)
		results = append(results, port.FaceEngineResult{
			Box:       f.Box,
			Score:     f.Score,
			Embedding: emb,
		})
	}
	return results, nil
}
