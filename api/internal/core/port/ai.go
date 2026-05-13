package port

// ImageAnalysis represents the result of analyzing an image with AI
type ImageAnalysis struct {
	Caption    string   `json:"caption"`
	Tags       []string `json:"tags"`
	Objects    []string `json:"objects"`
	IsNSFW     bool     `json:"is_nsfw"`
	IsPortrait bool     `json:"is_portrait"`
}

// AIService is an interface for AI-powered image analysis
type AIService interface {
	// AnalyzeImage analyzes an image at the given path and returns structured metadata
	AnalyzeImage(imagePath string) (*ImageAnalysis, error)
}
