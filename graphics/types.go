package graphics

import (
	"fmt"
	"image"
	"image/draw"
	"reflect"
)

type Image struct {
	Imager

	// Pix holds the image's pixels, with the order depending on the underlying image type. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle

	// Bytes per pixel. Calculated during NewImage, and used to ensure provided pixelBytes parameters in received
	// methods are <= the number of bytes used for each pixel in the image format.
	bpp int
}

// NewImage is a factory method to create an Image from an Imager.
// Modified from: https://stackoverflow.com/a/52164510/5061881.
//
// License(s):
// https://creativecommons.org/licenses/by-sa/4.0/
func NewImage(imgr Imager) (*Image, error) {
	img := &Image{
		Imager: imgr,
	}

	v := reflect.ValueOf(imgr)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		pv := v.FieldByName("Pix")
		ps := v.FieldByName("Stride")
		pr := v.FieldByName("Rect")
		if pv.IsValid() && ps.IsValid() && pr.IsValid() {
			pix, ok1 := pv.Interface().([]uint8)
			stride, ok2 := ps.Interface().(int)
			rect, ok3 := pr.Interface().(image.Rectangle)
			if ok1 && ok2 && ok3 {
				img.Pix = pix
				img.Stride = stride
				img.Rect = rect
				img.bpp = len(pix) / (rect.Dx() * rect.Dy())
				return img, nil
			}
		}
	}

	return nil, fmt.Errorf("unknown image type %T", imgr)
}

// You probably don't want to use this. Create a graphics.Image instead using the NewImage factory. This will allow
// the use of methods such as DrawFilledCircle
type Imager interface {
	draw.Image
	PixOffset(x, y int) int
}

type SubImager interface {
	Imager
	SubImage(r image.Rectangle) image.Image
}
