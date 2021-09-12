package graphics

// Todo: everything assumes r,g,b,a, and SetPixel assumes RGBA type byte order

// DrawCircleBorder draws a rasterized circle border (ring 1 pixel wide), centered on (cx, cy) and of the
// color provided by r,g,b,a, using the Midpoint Circle algorithm.
func (img *Image) DrawCircleBorder(cx, cy, rad int, r, g, b, a uint8) {
	// If circle falls entirely outside the image, return
	if (cx+rad < 0 || cx-rad > img.Bounds().Dx()) && (cy+rad < 0 || cy-rad > img.Bounds().Dy()) {
		return
	}

	dx, dy, ex, ey := rad-1, 0, 1, 1
	err := ex - (rad * 2)

	for dx > dy {
		img.SetPixel(cx+dx, cy+dy, r, g, b, a)
		img.SetPixel(cx+dy, cy+dx, r, g, b, a)
		img.SetPixel(cx-dy, cy+dx, r, g, b, a)
		img.SetPixel(cx-dx, cy+dy, r, g, b, a)
		img.SetPixel(cx-dx, cy-dy, r, g, b, a)
		img.SetPixel(cx-dy, cy-dx, r, g, b, a)
		img.SetPixel(cx+dy, cy-dx, r, g, b, a)
		img.SetPixel(cx+dx, cy-dy, r, g, b, a)

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

// DrawFilledCircle draws a filled-in (rasterized) circle, centered on (cx, cy) and of the color provided by r,g,b,a,
// using a (heavy) modification to the Midpoint Circle algorithm.
// This method is adapted from https://stackoverflow.com/q/10878209/5061881.
//
// License(s):
// https://creativecommons.org/licenses/by-sa/3.0/
func (img *Image) DrawFilledCircle(cx, cy, rad int, r, g, b, a uint8) {
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

		img.drawTwoCenteredLines(cx, cy, x, lastY, r, g, b, a)

		if err >= 0 {
			if x != lastY {
				img.drawTwoCenteredLines(cx, cy, lastY, x, r, g, b, a)
			}

			err -= x
			x--
			err -= x
		}
	}
}

// drawTwoCenteredLines draws two lines of length 2*dx+1, centered on (cx,cy) and of the color provided by r,g,b,a,
// and with a gap of 2*dx-1 rows/pixels between them (that is, the line at cy and dy-1 lines to either side of it are
// not drawn).
// This is used by DrawFilledCircle. See attribution there.
func (img *Image) drawTwoCenteredLines(cx, cy, dx, dy int, r, g, b, a uint8) {
	img.DrawHLine(cx-dx, cy+dy, cx+dx, r, g, b, a)
	if dy != 0 {
		img.DrawHLine(cx-dx, cy-dy, cx+dx, r, g, b, a)
	}
}

// DrawHLine draws a horizontal line from (x0,y0) to (x1,y0), of the color provided by r,g,b,a.
func (img *Image) DrawHLine(x0, y0, x1 int, r, g, b, a uint8) {
	for x := x0; x <= x1; x++ {
		img.SetPixel(x, y0, r, g, b, a)
	}
}

// DrawVLine draws a vertical line from (x0,y0) to (x0,y1), of the color provided by r,g,b,a.
func (img *Image) DrawVLine(x0, y0, y1 int, r, g, b, a uint8) {
	for y := y0; y <= y1; y++ {
		img.SetPixel(x0, y, r, g, b, a)
	}
}

// SetPixel sets the pixel at (x,y) on s.img to the color provided by r,g,b,a.
func (img *Image) SetPixel(x, y int, r, g, b, a uint8) {
	// Setting the pixel color bytes in the back-buffer is >5x the speed of img.Set()
	o := img.PixOffset(x, y)
	if o < 0 || o >= len(img.Pix) {
		return
	}

	// Locks are only necessary if multithreading (and not then if very rare write failures are acceptable - it's just a slice)
	//s.imgLock.Lock()
	img.Pix[o], img.Pix[o+1], img.Pix[o+2], img.Pix[o+3] = r, g, b, a
	//s.imgLock.Unlock()
}
