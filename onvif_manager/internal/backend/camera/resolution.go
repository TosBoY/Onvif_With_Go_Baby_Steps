package camera

import (
	"math"
	"onvif_manager/pkg/models"
)

func FindClosestResolution(target models.Resolution, available []models.Resolution) models.Resolution {
	if len(available) == 0 {
		return models.Resolution{}
	}

	targetArea := float64(target.Width * target.Height)
	targetRatio := float64(target.Width) / float64(target.Height)

	type candidate struct {
		resolution models.Resolution
		areaDiff   float64
		ratioDiff  float64
	}

	const ratioThreshold = 0.5

	var candidatesWithinRatio []candidate
	var allCandidates []candidate

	for _, res := range available {
		resArea := float64(res.Width * res.Height)
		resRatio := float64(res.Width) / float64(res.Height)

		areaDiff := math.Abs(targetArea - resArea)
		ratioDiff := math.Abs(targetRatio - resRatio)

		c := candidate{
			resolution: res,
			areaDiff:   areaDiff,
			ratioDiff:  ratioDiff,
		}

		allCandidates = append(allCandidates, c)

		if ratioDiff <= ratioThreshold {
			candidatesWithinRatio = append(candidatesWithinRatio, c)
		}
	}

	// Prefer closest area from candidates with acceptable aspect ratio
	if len(candidatesWithinRatio) > 0 {
		best := candidatesWithinRatio[0]
		for _, c := range candidatesWithinRatio[1:] {
			if c.areaDiff < best.areaDiff {
				best = c
			}
		}
		return best.resolution
	}

	// If no acceptable aspect ratios, pick the one with closest ratio,
	// and then closest area among those with similar ratio
	best := allCandidates[0]
	for _, c := range allCandidates[1:] {
		if c.ratioDiff < best.ratioDiff || (c.ratioDiff == best.ratioDiff && c.areaDiff < best.areaDiff) {
			best = c
		}
	}
	return best.resolution
}
