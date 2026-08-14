package service

import (
	"fmt"
	"image"
	"log/slog"
	"math"

	"github.com/disintegration/imaging"
)

// QualityAnalysis holds deterministic, technical photo-quality metrics.
// These are computed locally (no AI) so they are fast and reproducible.
type QualityAnalysis struct {
	// BlurScore is the Laplacian variance (higher = sharper). Threshold
	// values are tuned for typical consumer photos.
	BlurScore float64
	// ExposureScore ranges 0-100; 100 is a well-exposed image, lower
	// values indicate over/under-exposure.
	ExposureScore float64
	// Overall is a composite 0-100 score.
	Overall int
	// IsLowQuality is true when the photo is likely unusable (blurry,
	// badly exposed, or a near-solid image such as a lens cap).
	IsLowQuality bool
	// IsNearlySolid is true when the image is mostly one colour (lens
	// cap, black frame, etc.).
	IsNearlySolid bool
}

// qualityThresholds controls classification.
var qualityThresholds = struct {
	blurLow       float64 // below this => blurry
	exposureLow   float64 // below this => badly exposed
	solidRatio    float64 // fraction of dominant luminance bin => nearly solid
	nearSolidVar  float64 // Laplacian variance below this + solid => useless
	overallLow    int     // composite below this => low quality
}{
	blurLow:      30,
	exposureLow:  30,
	solidRatio:   0.85,
	nearSolidVar: 20,
	overallLow:   40,
}

// Laplacian kernel emphasises edges; a blurry image has low edge energy.
var laplacianKernel = [9]float64{
	0, 1, 0,
	1, -4, 1,
	0, 1, 0,
}

// AnalyzeQuality computes quality metrics for an image file.
// Returns (nil, nil) when the image cannot be decoded (not an error — some
// files in the library may be unreadable by design).
func AnalyzeQuality(path string) (*QualityAnalysis, error) {
	src, err := imaging.Open(path)
	if err != nil {
		slog.Debug("quality: unable to open image", "path", path, "error", err)
		return nil, nil
	}

	// Downscale to a bounded size so large originals don't cost much.
	src = imaging.Fit(src, 480, 480, imaging.Lanczos)

	gray := imaging.Grayscale(src)
	lap := imaging.Convolve3x3(gray, laplacianKernel, nil)

	// Blur uses the Laplacian edge response; exposure + solidity are
	// computed from the original grayscale, not the (mostly-black) edge map.
	blurScore := laplacianVariance(lap)
	exposureScore := exposureScoreFromHistogram(gray)
	_, isNearlySolid := solidRatioOf(gray)

	composite := compositeScore(blurScore, exposureScore, isNearlySolid)
	overall := int(math.Round(composite))

	isLowQuality := overall < qualityThresholds.overallLow ||
		(isNearlySolid && blurScore < qualityThresholds.nearSolidVar)

	return &QualityAnalysis{
		BlurScore:     blurScore,
		ExposureScore: exposureScore,
		Overall:       overall,
		IsLowQuality:  isLowQuality,
		IsNearlySolid: isNearlySolid,
	}, nil
}

// laplacianVariance computes the variance of the Laplacian response over the
// whole (grayscale) image — the classic blur metric.
func laplacianVariance(img image.Image) float64 {
	bounds := img.Bounds()
	if bounds.Dx() < 3 || bounds.Dy() < 3 {
		return math.MaxFloat64 // tiny images are uninformative; treat as sharp
	}

	var sum, sumSq float64
	var count int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// luminance-weighted average of RGB (each 0-65535)
			v := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 65535.0 * 255.0
			sum += v
			sumSq += v * v
			count++
		}
	}
	if count == 0 {
		return 0
	}
	mean := sum / float64(count)
	variance := (sumSq / float64(count)) - (mean * mean)
	if variance < 0 {
		variance = 0
	}
	return variance
}

// exposureScoreFromHistogram bins luminance and penalises images dominated
// by the extreme ends (crushed blacks / blown whites).
func exposureScoreFromHistogram(img image.Image) float64 {
	const bins = 16
	var hist [bins]int
	total := 0

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			lum := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 65535.0
			idx := int(lum * bins)
			if idx >= bins {
				idx = bins - 1
			}
			hist[idx]++
			total++
		}
	}
	if total == 0 {
		return 100
	}

	// Middle bins are "well exposed". Weight edges (0 = black, 15 = white)
	// as bad exposure.
	var weighted float64
	for i, count := range hist {
		frac := float64(count) / float64(total)
		if i == 0 || i == bins-1 {
			weighted += frac * 0.2 // heavily penalise crushed/blown extremes
		} else {
			weighted += frac * 1.0
		}
	}
	score := weighted * 100.0
	if score > 100 {
		score = 100
	}
	return score
}

// solidRatioOf returns the fraction of pixels in the most common luminance
// bin, and whether that fraction suggests a nearly-solid image.
func solidRatioOf(img image.Image) (ratio float64, isSolid bool) {
	const bins = 32
	var hist [bins]int
	total := 0

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			lum := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 65535.0
			idx := int(lum * bins)
			if idx >= bins {
				idx = bins - 1
			}
			hist[idx]++
			total++
		}
	}
	if total == 0 {
		return 1, true
	}
	max := 0
	for _, c := range hist {
		if c > max {
			max = c
		}
	}
	ratio = float64(max) / float64(total)
	return ratio, ratio >= qualityThresholds.solidRatio
}

// compositeScore blends blur + exposure into 0-100.
func compositeScore(blur, exposure float64, isNearlySolid bool) float64 {
	// Map blur variance (0-~500 typical) to 0-100; higher variance => sharper.
	blurNorm := math.Min(blur/300.0, 1.0) * 100.0
	score := 0.6*blurNorm + 0.4*exposure
	if isNearlySolid {
		score *= 0.3 // mostly-solid frames are almost always junk
	}
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}
	return score
}

// qualityLabel returns a human label for an analysis.
func (q *QualityAnalysis) qualityLabel() string {
	switch {
	case q == nil:
		return "Unknown"
	case q.IsNearlySolid:
		return "Empty"
	case q.IsLowQuality:
		return "Low quality"
	case q.Overall >= 75:
		return "Good"
	default:
		return "Acceptable"
	}
}

// String is a compact summary for logs.
func (q *QualityAnalysis) String() string {
	if q == nil {
		return "quality=unknown"
	}
	return fmt.Sprintf("quality=%d blur=%.1f exposure=%.1f low=%t solid=%t",
		q.Overall, q.BlurScore, q.ExposureScore, q.IsLowQuality, q.IsNearlySolid)
}
