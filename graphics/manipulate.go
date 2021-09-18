package graphics

import (
	"image"

	"github.com/nfnt/resize"
)

// ResizeMaintain resizes img while maintaining its aspect ratio and ensuring that the new image fills the target size.
// The source image will be cut off in the smaller target dimension if they are not the same aspect ratio.
// The upper-left corners are aligned.
// This method uses the Nearest Neighbor interpolation algorithm. For other algorithms, use ResizeMaintainWithInterp.
func ResizeMaintain(img image.Image, targetWidth, targetHeight uint) image.Image {
	return ResizeMaintainWithInterp(img, targetWidth, targetHeight, resize.NearestNeighbor)
}

// ResizeMaintainWithInterp resizes img while maintaining its aspect ratio and ensuring that the new image fills the target size.
// The source image will be cut off in the smaller target dimension if they are not the same aspect ratio.
// The upper-left corners are aligned.
// The resize is performed using the interpolation algorithm provided by function. Note that resize.NearestNeighbor is
// the fastest available algorithm, but will not always produce clean results.
func ResizeMaintainWithInterp(img image.Image, targetWidth, targetHeight uint, function resize.InterpolationFunction) image.Image {
	// Todo: I think this check is only sufficient if img height is larger than target. Does it just need to be reversed if not, or is it more complex than that (dependent on which width is larger as well)?
	if targetWidth >= targetHeight && img.Bounds().Dx() >= img.Bounds().Dy() {
		return resize.Resize(0, uint(targetHeight), img, function)
	} else {
		return resize.Resize(uint(targetWidth), 0, img, function)
	}
}
