package service

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// writeTestPNG writes a solid-colour PNG of the given size.
func writeTestPNG(t *testing.T, dir, name string, w, h int, fill color.RGBA) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: fill}, image.Point{}, draw.Src)
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	assert.NoError(t, err)
	defer func() { _ = f.Close() }()
	assert.NoError(t, png.Encode(f, img))
	return path
}

// writeNoisyPNG writes an image filled with random noise (high Laplacian
// variance => "sharp").
func writeNoisyPNG(t *testing.T, dir, name string, w, h int) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	seed := uint32(42)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			seed = seed*1664525 + 1013904223
			v := uint8(seed >> 24)
			img.SetRGBA(x, y, color.RGBA{R: v, G: v, B: v, A: 255})
		}
	}
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	assert.NoError(t, err)
	defer func() { _ = f.Close() }()
	assert.NoError(t, png.Encode(f, img))
	return path
}

// writeGradientPNG writes a smooth vertical gradient (low variance, but not
// solid — represents a well-exposed, moderately detailed photo).
func writeGradientPNG(t *testing.T, dir, name string, w, h int) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		v := uint8(30 + (y*180)/h) // 30..210, well-exposed
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, color.RGBA{R: v, G: v, B: v, A: 255})
		}
	}
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	assert.NoError(t, err)
	defer func() { _ = f.Close() }()
	assert.NoError(t, png.Encode(f, img))
	return path
}

func TestAnalyzeQuality(t *testing.T) {
	dir := t.TempDir()

	t.Run("nearly-solid image is low quality", func(t *testing.T) {
		path := writeTestPNG(t, dir, "solid.png", 200, 200, color.RGBA{R: 10, G: 10, B: 10, A: 255})
		q, err := AnalyzeQuality(path)
		assert.NoError(t, err)
		assert.NotNil(t, q)
		assert.True(t, q.IsNearlySolid, "solid image should be detected as nearly-solid")
		assert.True(t, q.IsLowQuality)
		assert.Less(t, q.Overall, 40)
	})

	t.Run("noisy image is sharp", func(t *testing.T) {
		path := writeNoisyPNG(t, dir, "noisy.png", 200, 200)
		q, err := AnalyzeQuality(path)
		assert.NoError(t, err)
		assert.NotNil(t, q)
		assert.False(t, q.IsNearlySolid)
		assert.Greater(t, q.BlurScore, float64(100), "noise should produce high Laplacian variance")
		assert.False(t, q.IsLowQuality)
	})

	t.Run("gradient image is acceptable and not solid", func(t *testing.T) {
		path := writeGradientPNG(t, dir, "gradient.png", 200, 200)
		q, err := AnalyzeQuality(path)
		assert.NoError(t, err)
		assert.NotNil(t, q)
		assert.False(t, q.IsNearlySolid)
		assert.GreaterOrEqual(t, q.ExposureScore, float64(30))
	})

	t.Run("unreadable file returns nil without error", func(t *testing.T) {
		path := filepath.Join(dir, "fake.jpg")
		assert.NoError(t, os.WriteFile(path, []byte("not an image"), 0644))
		q, err := AnalyzeQuality(path)
		assert.NoError(t, err)
		assert.Nil(t, q)
	})
}

func TestQualityLabel(t *testing.T) {
	assert.Equal(t, "Unknown", (*QualityAnalysis)(nil).qualityLabel())
	assert.Equal(t, "Empty", (&QualityAnalysis{IsNearlySolid: true}).qualityLabel())
	assert.Equal(t, "Low quality", (&QualityAnalysis{IsLowQuality: true}).qualityLabel())
	assert.Equal(t, "Good", (&QualityAnalysis{Overall: 90}).qualityLabel())
	assert.Equal(t, "Acceptable", (&QualityAnalysis{Overall: 60}).qualityLabel())
}
