package service

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

func TestApplyEditPipeline(t *testing.T) {
	// 100x60 image with a solid left half and white right half.
	base := image.NewNRGBA(image.Rect(0, 0, 100, 60))
	draw.Draw(base, image.Rect(0, 0, 50, 60), &image.Uniform{C: color.RGBA{R: 100, G: 100, B: 100, A: 255}}, image.Point{}, draw.Src)
	draw.Draw(base, image.Rect(50, 0, 100, 60), &image.Uniform{C: color.RGBA{R: 255, G: 255, B: 255, A: 255}}, image.Point{}, draw.Src)

	t.Run("rotate 90 swaps dimensions", func(t *testing.T) {
		rot := 90
		out := applyEditPipeline(base, EditParams{Rotate: &rot})
		b := out.Bounds()
		assert.Equal(t, 60, b.Dx(), "after 90deg rotate width should be original height")
		assert.Equal(t, 100, b.Dy(), "after 90deg rotate height should be original width")
	})

	t.Run("crop keeps only the crop region", func(t *testing.T) {
		crop := domain.CropParams{X: 0, Y: 0, Width: 0.5, Height: 1.0}
		out := applyEditPipeline(base, EditParams{Crop: &crop})
		b := out.Bounds()
		assert.Equal(t, 50, b.Dx(), "half-width crop")
		assert.Equal(t, 60, b.Dy(), "full-height crop")
	})

	t.Run("rotate then crop compose", func(t *testing.T) {
		rot := 90
		crop := domain.CropParams{X: 0, Y: 0, Width: 1.0, Height: 0.5}
		out := applyEditPipeline(base, EditParams{Rotate: &rot, Crop: &crop})
		b := out.Bounds()
		assert.Equal(t, 60, b.Dx())
		assert.Equal(t, 50, b.Dy(), "half of rotated height")
	})

	t.Run("brightness changes pixel value", func(t *testing.T) {
		br := 50.0
		out := applyEditPipeline(base, EditParams{Brightness: &br})
		// The dark region should get brighter.
		r, _, _, _ := out.At(10, 30).RGBA()
		r0, _, _, _ := base.At(10, 30).RGBA()
		assert.Greater(t, r, r0, "brightness+50 should brighten the dark region")
	})

	t.Run("no params returns clone with same dimensions", func(t *testing.T) {
		out := applyEditPipeline(base, EditParams{})
		assert.Equal(t, base.Bounds().Dx(), out.Bounds().Dx())
		assert.Equal(t, base.Bounds().Dy(), out.Bounds().Dy())
	})
}

func TestApplyEditsForPhoto(t *testing.T) {
	dir := t.TempDir()

	// Create a source image file
	src := image.NewNRGBA(image.Rect(0, 0, 80, 40))
	draw.Draw(src, src.Bounds(), &image.Uniform{C: color.RGBA{R: 50, G: 50, B: 50, A: 255}}, image.Point{}, draw.Src)
	srcPath := filepath.Join(dir, "orig.png")
	f, _ := os.Create(srcPath)
	_ = png.Encode(f, src)
	_ = f.Close()

	t.Run("writes edited copy", func(t *testing.T) {
		rot := 90
		outPath, err := ApplyEditsForPhoto(dir, "photo-abc", srcPath, EditParams{Rotate: &rot})
		assert.NoError(t, err)
		assert.NotEmpty(t, outPath)
		assert.FileExists(t, outPath)

		img, err := imaging.Open(outPath)
		assert.NoError(t, err)
		assert.Equal(t, 40, img.Bounds().Dx(), "rotated copy should be 40x80")
		assert.Equal(t, 80, img.Bounds().Dy())
	})

	t.Run("missing source errors", func(t *testing.T) {
		_, err := ApplyEditsForPhoto(dir, "photo-xyz", filepath.Join(dir, "missing.jpg"), EditParams{})
		assert.Error(t, err)
	})
}
