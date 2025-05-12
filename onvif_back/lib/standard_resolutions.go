package onvif_test

import "math"

type StandardResolution struct {
	Label  string
	Width  int
	Height int
}

var StandardResolutions = []StandardResolution{
	{Label: "480p", Width: 640, Height: 480},    // SD
	{Label: "720p", Width: 1280, Height: 720},   // HD
	{Label: "1080p", Width: 1920, Height: 1080}, // Full HD
	{Label: "2K", Width: 2048, Height: 1080},    // 2K DCI
	{Label: "1440p", Width: 2560, Height: 1440}, // QHD
	{Label: "4K", Width: 3840, Height: 2160},    // 4K UHD
	{Label: "5K", Width: 5120, Height: 2880},    // 5K
	{Label: "8K", Width: 7680, Height: 4320},    // 8K UHD
}

func FindClosestStandardResolution(width, height int) StandardResolution {
	closestRes := StandardResolution{Label: "Custom", Width: width, Height: height}
	minDiff := float64(width * height) // Initialize with current resolution difference

	targetArea := float64(width * height)
	for _, std := range StandardResolutions {
		stdArea := float64(std.Width * std.Height)
		diff := math.Abs(targetArea - stdArea)

		if diff < minDiff {
			// Check aspect ratio similarity as well
			targetRatio := float64(width) / float64(height)
			stdRatio := float64(std.Width) / float64(std.Height)
			ratioDiff := math.Abs(targetRatio - stdRatio)

			// Only update if the aspect ratio is similar enough (within 10%)
			if ratioDiff < 0.1 {
				minDiff = diff
				closestRes = std
			}
		}
	}

	return closestRes
}
