package graphics

import (
	"image"
	"math"

	"github.com/nfnt/resize"
)

// ResizeMaintain resizes img while maintaining its aspect ratio and ensuring that the new image fills the target size.
// The source image will be cropped in the smaller target dimension if they are not the same aspect ratio.
// The upper-left corners are aligned.
// This method uses the Nearest Neighbor interpolation algorithm. For other algorithms, use ResizeMaintainWithInterp.
func ResizeMaintain(img SubImager, targetWidth, targetHeight uint) image.Image {
	return ResizeMaintainWithInterp(img, targetWidth, targetHeight, resize.NearestNeighbor)
}

// ResizeMaintainWithInterp resizes img while maintaining its aspect ratio and ensuring that the new image fills the target size.
// The source image will be cropped in the smaller target dimension if they are not the same aspect ratio.
// The upper-left corners are aligned.
// The resize is performed using the interpolation algorithm provided by function. Note that resize.NearestNeighbor is
// the fastest available algorithm, but will not always produce clean results.
func ResizeMaintainWithInterp(img SubImager, targetWidth, targetHeight uint, function resize.InterpolationFunction) image.Image {
	// Todo: the logic for which dimension to preserve is fucky. E.g. if the second condition of the second if is changed to float division, it breaks on BlueMoon
	var r image.Image
	var w, h uint
	if targetWidth >= targetHeight && img.Bounds().Dx() >= img.Bounds().Dy() {
		if float32(img.Bounds().Dy())/float32(targetHeight) <= float32(img.Bounds().Dx())/float32(targetWidth) &&
			img.Bounds().Dx()/img.Bounds().Dy() <= int(targetWidth/targetHeight) {
			w = 0
			h = targetHeight
		} else {
			w = targetWidth
			h = 0
		}
	} else {
		if float32(img.Bounds().Dx())/float32(targetWidth) <= float32(img.Bounds().Dy())/float32(targetHeight) &&
			img.Bounds().Dy()/img.Bounds().Dx() <= int(targetHeight/targetWidth) {
			w = targetWidth
			h = 0
		} else {
			w = 0
			h = targetHeight
		}
	}
	r = resize.Resize(w, h, img, function)
	return r.(SubImager).SubImage(image.Rectangle{Max: image.Point{X: int(targetWidth), Y: int(targetHeight)}})
}

// MultiplyRect resizes (multiplies) r by factor. r.Min will remain the same. The new r.Size() will be r.Size().Mul(factor).
// If factor is < 0, r is returned unchanged. If factor = 0, the zero rectangle is returned.
func MultiplyRect(r image.Rectangle, factor float64) image.Rectangle {
	if factor < 0 || factor == 1 {
		return r
	}
	if factor == 0 {
		return image.Rectangle{}
	}

	n := r.Sub(r.Min)
	n.Max.X = int(math.Round(float64(n.Max.X) * factor))
	n.Max.Y = int(math.Round(float64(n.Max.Y) * factor))

	return n.Add(r.Min)
}

// CloneFrom clones the pixel data from src into img.
// There will be unexpected results if the fields of the Images don't match or the underlying image.Image Color Models
// don't match (in that case, use draw.Draw).
func (img *Image) CloneFrom(src *Image) {
	l := copy(img.Pix, src.Pix)
	img.Pix = img.Pix[0:l:l]
}

// CloneFromRange clones the pixel data from src into img, in the range [from,to).
// It will panic if from or to are < 0 or > the length of either Image's Pix slice or from > to.
// There will be unexpected results if the fields of the Images (particularly Bounds().Size() or Stride) don't match
// or the underlying image.Image Color Models don't match (in that case, use draw.Draw).
func (img *Image) CloneFromRange(src *Image, from, to int) {
	copy(img.Pix[from:to], src.Pix[from:to])
}

// CloneFromRows clones the pixel data from src into img, for the rows (y values) in the range [from,to].
// It will panic if from or to are < 0 or >= Bounds().Size().Y or from > to.
// There will be unexpected results if the fields of the Images (particularly Bounds().Size() or Stride) don't match
// or the underlying image.Image Color Models don't match (in that case, use draw.Draw).
func (img *Image) CloneFromRows(src *Image, from, to int) {
	start := src.PixOffset(0, from)
	// The start of the row after to is the (unincluded) end index of the slice to cpy
	end := src.PixOffset(0, to+1)
	copy(img.Pix[start:end], src.Pix[start:end])
}

// CloneFromRect clones the pixel data within rect (pixel coordinates) from src into img.
// It will panic if rect is not fully contained by both img and src (their Bounds().Size(), not their Bounds()).
// The images are assumed to be the same size; the rect is shared between images.
// It is equivalent to draw.Draw(img, rect, src, rect.Min, draw.Src), but much faster thanks to
// specific-case optimization.
// There will be unexpected results if the fields of the Images (particularly Bounds().Size() or Stride) don't match
// or the underlying image.Image Color Models don't match (in that case, use draw.Draw).
func (img *Image) CloneFromRect(src *Image, rect image.Rectangle) {
	var start int
	dx := rect.Dx() * img.bpp
	x0 := rect.Min.X * img.bpp
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		start = y*img.Stride + x0
		copy(img.Pix[start:start+dx], src.Pix[start:start+dx])
	}
}

// PlaceAtPoint copies (all of) src onto img at pt on img. It is equivalent to
// draw.Draw(img, src.Bounds().Sub(src.Bounds().Min), src, image.Point{}, draw.Src), but >50x the speed thanks to
// specific-case optimization.
func (img *Image) PlaceAtPoint(src *image.RGBA, pt image.Point) {
	// Since we know we'll be repeating for each row, calculating the indexes manually with values that don't change
	// computed only once is about 2x the speed of using PixOffset.
	var imgStart, srcStart int
	dx := src.Bounds().Dx() * 4
	x0 := pt.X * img.bpp
	for y := 0; y < src.Bounds().Dy(); y++ {
		imgStart = (pt.Y+y)*img.Stride + x0
		srcStart = y * src.Stride
		copy(img.Pix[imgStart:imgStart+dx], src.Pix[srcStart:srcStart+dx])
	}
}
