package camera

import (
	"main_back/pkg/models"
	"math"
)

func FindClosestResolution(target models.Resolution, available []models.Resolution) models.Resolution {
	if len(available) == 0 {
		return models.Resolution{} // Return zero value if no available resolutions
	}

	closestRes := available[0] // Initialize with the first available resolution
	minDiff := float64(abs(target.Width*target.Height - closestRes.Width*closestRes.Height))

	targetArea := float64(target.Width * target.Height)
	targetRatio := float64(target.Width) / float64(target.Height)

	for _, res := range available {
		stdArea := float64(res.Width * res.Height)
		diff := math.Abs(targetArea - stdArea)

		if diff < minDiff {
			// Check aspect ratio similarity as well
			stdRatio := float64(res.Width) / float64(res.Height)
			ratioDiff := math.Abs(targetRatio - stdRatio)

			// Only update if the aspect ratio is similar enough (within 10%)
			if ratioDiff < 0.1 || targetRatio == 0 || stdRatio == 0 {
				minDiff = diff
				closestRes = res
			}
		}
	}

	return closestRes
}

// Helper function for absolute value of int
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
