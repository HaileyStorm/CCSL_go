package graphics

// Todo: everything assumes pixelBytes, and SetPixel assumes RGBA type byte order

// DrawCircleBorder draws a rasterized circle border (ring 1 pixel wide), centered on (cx, cy) and of the
// color provided by pixelBytes, using the Midpoint Circle algorithm.
// pixelBytes may be the first n bytes of a pixel may be provided instead of all bytes.
func (img *Image) DrawCircleBorder(cx, cy, rad int, pixelBytes ...uint8) {
	// If circle falls entirely outside the image, return
	if (cx+rad < 0 || cx-rad > img.Bounds().Dx()) && (cy+rad < 0 || cy-rad > img.Bounds().Dy()) {
		return
	}

	dx, dy, ex, ey := rad-1, 0, 1, 1
	err := ex - (rad * 2)

	for dx > dy {
		img.SetPixel(cx+dx, cy+dy, pixelBytes...)
		img.SetPixel(cx+dy, cy+dx, pixelBytes...)
		img.SetPixel(cx-dy, cy+dx, pixelBytes...)
		img.SetPixel(cx-dx, cy+dy, pixelBytes...)
		img.SetPixel(cx-dx, cy-dy, pixelBytes...)
		img.SetPixel(cx-dy, cy-dx, pixelBytes...)
		img.SetPixel(cx+dy, cy-dx, pixelBytes...)
		img.SetPixel(cx+dx, cy-dy, pixelBytes...)

		if err <= 0 {
			dy++
			err += ey
			ey += 2
		}
		if err > 0 {
			dx--
			ex += 2
			err += ex - (rad * 2)
		}
	}
}

// DrawFilledCircle draws a filled-in (rasterized) circle, centered on (cx, cy) and of the color provided by pixelBytes,
// using a (heavy) modification to the Midpoint Circle algorithm.
// pixelBytes may be the first n bytes of a pixel may be provided instead of all bytes.
// This method is adapted from https://stackoverflow.com/q/10878209/5061881.
//
// License(s):
// https://creativecommons.org/licenses/by-sa/3.0/
func (img *Image) DrawFilledCircle(cx, cy, rad int, pixelBytes ...uint8) {
	// If circle falls entirely outside the environment, return
	if (cx+rad < 0 || cx-rad > img.Bounds().Dx()) && (cy+rad < 0 || cy-rad > img.Bounds().Dy()) {
		return
	}

	err, x, y := -rad, rad, 0
	var lastY int

	for x >= y {
		lastY = y
		err += y
		y++
		err += y

		img.drawTwoCenteredLines(cx, cy, x, lastY, pixelBytes...)

		if err >= 0 {
			if x != lastY {
				img.drawTwoCenteredLines(cx, cy, lastY, x, pixelBytes...)
			}

			err -= x
			x--
			err -= x
		}
	}
}

// drawTwoCenteredLines draws two lines of length 2*dx+1, centered on (cx,cy) and of the color provided by pixelBytes,
// and with a gap of 2*dx-1 rows/pixels between them (that is, the line at cy and dy-1 lines to either side of it are
// not drawn).
// pixelBytes may be the first n bytes of a pixel may be provided instead of all bytes.
// This is used by DrawFilledCircle. See attribution there.
func (img *Image) drawTwoCenteredLines(cx, cy, dx, dy int, pixelBytes ...uint8) {
	img.DrawHLine(cx-dx, cy+dy, cx+dx, pixelBytes...)
	if dy != 0 {
		img.DrawHLine(cx-dx, cy-dy, cx+dx, pixelBytes...)
	}
}

// DrawHLine draws a horizontal line from (x0,y0) to (x1,y0), of the color provided by pixelBytes.
// pixelBytes may be the first n bytes of a pixel may be provided instead of all bytes.
func (img *Image) DrawHLine(x0, y0, x1 int, pixelBytes ...uint8) {
	for x := x0; x <= x1; x++ {
		img.SetPixel(x, y0, pixelBytes...)
	}
}

// DrawVLine draws a vertical line from (x0,y0) to (x0,y1), of the color provided by pixelBytes.
// pixelBytes may be the first n bytes of a pixel may be provided instead of all bytes.
func (img *Image) DrawVLine(x0, y0, y1 int, pixelBytes ...uint8) {
	for y := y0; y <= y1; y++ {
		img.SetPixel(x0, y, pixelBytes...)
	}
}

// SetPixel sets the pixel at (x,y) on s.img to the color provided by pixelBytes.
// Returning an error comes at a ~1.7-2.7% performance hit. Given the typical use cases for this function, it is probably
// best not to return errors. This puts the onus on the end user to ensure they are providing in-bound pixel coordinates
// and correct pixelBytes.
// pixelBytes may be the first n bytes of a pixel may be provided instead of all bytes.
func (img *Image) SetPixel(x, y int, pixelBytes ...uint8) { //error {
	// This function is ~2.5-5x the speed of img.Set() (because of color conversion etc.)

	// Checking > instead of != allows not providing all the bytes for the pixel
	// (only changing the first len(pixelBytes) bytes of the pixel).
	// Useful particularly for leaving out alpha e.g. when image is always fully opaque.
	if len(pixelBytes) > img.bpp {
		return //fmt.Errorf("length of pixelBytes (%d) cannot exceed the number of pixels/byte for the image format (%d)",
		//len(pixelBytes), img.bpp)
	}

	// This way of checking bounds is ~1.5% faster than the standard library point in rect check (and the offset itself -
	// or rather, the difference in execution time between it and the rect check - will only be wasted when the point
	// is outside the image, which in good programming should be comparatively rare (particularly if there are a lot
	// of calls being made in a short period of time and efficiency matters, calls to here should be pretty optimized)).
	o := img.PixOffset(x, y)
	if o < 0 || o >= len(img.Pix) {
		return //fmt.Errorf("point %v is outside the size of the image %v", image.Point{X: x, Y: y}, img.Rect.Size())
	}

	s := img.Pix[o : o+len(pixelBytes) : o+len(pixelBytes)]

	for i := 0; i < len(pixelBytes); i++ {
		s[i] = pixelBytes[i]
	}

	//return nil
}
