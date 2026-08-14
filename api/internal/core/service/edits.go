package service

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
)

// EditParams is an alias for the domain edit parameters.
type EditParams = domain.EditParams

// editedFileName returns the on-disk path for a photo's edited copy.
func editedFileName(configPhotoDir, photoId string) string {
	safeId := strings.ReplaceAll(photoId, "/", "_")
	return filepath.Join(configPhotoDir, ".edited", safeId+".jpg")
}

// ApplyEditsForPhoto writes the edited copy for a specific photo and
// returns its path.
func ApplyEditsForPhoto(photoDir, photoId, originalPath string, params EditParams) (string, error) {
	src, err := imaging.Open(originalPath)
	if err != nil {
		return "", fmt.Errorf("unable to open photo for edit: %w", err)
	}

	img := applyEditPipeline(src, params)

	dir := filepath.Join(photoDir, ".edited")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("unable to create edited dir: %w", err)
	}

	outPath := editedFileName(photoDir, photoId)
	f, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("unable to create edited copy: %w", err)
	}
	defer func() { _ = f.Close() }()

	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 92}); err != nil {
		return "", fmt.Errorf("unable to encode edited copy: %w", err)
	}
	slog.Info("Wrote edited copy", "photo", photoId, "path", outPath)
	return outPath, nil
}

// applyEditPipeline applies the full non-destructive pipeline to an image.
func applyEditPipeline(img image.Image, params EditParams) *image.NRGBA {
	out := imaging.Clone(img)

	if params.Rotate != nil {
		switch *params.Rotate % 360 {
		case 90:
			out = imaging.Rotate90(out)
		case 180:
			out = imaging.Rotate180(out)
		case 270:
			out = imaging.Rotate270(out)
		}
	}

	if params.Crop != nil && params.Crop.Width > 0 && params.Crop.Height > 0 {
		b := out.Bounds()
		out = imaging.Crop(out, image.Rect(
			int(params.Crop.X*float64(b.Dx())),
			int(params.Crop.Y*float64(b.Dy())),
			int((params.Crop.X+params.Crop.Width)*float64(b.Dx())),
			int((params.Crop.Y+params.Crop.Height)*float64(b.Dy())),
		))
	}

	if params.Brightness != nil {
		out = imaging.AdjustBrightness(out, *params.Brightness)
	}
	if params.Contrast != nil {
		out = imaging.AdjustContrast(out, *params.Contrast)
	}
	if params.Saturation != nil {
		out = imaging.AdjustSaturation(out, *params.Saturation)
	}
	if params.AutoEnhance != nil && *params.AutoEnhance {
		out = imaging.AdjustGamma(out, 0.85)
	}

	return out
}

// ParseEditParams unmarshals stored edit params, tolerating empty/null.
func ParseEditParams(raw []byte) (*EditParams, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return &EditParams{}, nil
	}
	var p EditParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// MarshalEditParams stores edit params as JSON.
func MarshalEditParams(p *EditParams) ([]byte, error) {
	return json.Marshal(p)
}

// UnescapePath is a small helper for callers.
func UnescapePath(p string) string {
	return utils.UnescapeInvalidCharacters(p)
}
