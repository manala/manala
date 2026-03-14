package scale

import (
	"image"

	"golang.org/x/image/draw"
)

// Scale resizes the source image to fit within the specified width and height while preserving aspect ratio.
// If the source is larger than the target dimensions, it scales down using nearest neighbor interpolation.
// If the source is smaller, it centers the image without scaling.
// The result is returned on a transparent canvas of the specified dimensions.
func Scale(src *image.NRGBA, width, height int) *image.NRGBA {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	// Create transparent canvas
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))

	// Center offset
	offsetX := (width - srcWidth) / 2
	offsetY := (height - srcHeight) / 2

	if srcWidth > width || srcHeight > height {
		// Scale down, preserve aspect ratio
		scaleX := float64(width) / float64(srcWidth)
		scaleY := float64(height) / float64(srcHeight)
		scale := scaleX
		if scaleY < scaleX {
			scale = scaleY
		}
		dstWidth := int(float64(srcWidth) * scale)
		dstHeight := int(float64(srcHeight) * scale)
		offsetX = (width - dstWidth) / 2
		offsetY = (height - dstHeight) / 2
		rect := image.Rect(offsetX, offsetY, offsetX+dstWidth, offsetY+dstHeight)
		draw.NearestNeighbor.Scale(dst, rect, src, srcBounds, draw.Over, nil)
	} else {
		// No scaling needed, direct copy
		rect := image.Rect(offsetX, offsetY, offsetX+srcWidth, offsetY+srcHeight)
		draw.Draw(dst, rect, src, srcBounds.Min, draw.Over)
	}

	return dst
}
