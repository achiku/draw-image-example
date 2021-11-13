package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func TestSrcOverDst(t *testing.T) {
	dst := image.NewRGBA(image.Rect(0, 0, 640, 480))
	blue := color.RGBA{0, 0, 255, 255}

	/*
		dst=destination image
		r=destination rectangles
		src=source image
		sp=source coordinates (image.ZP is the zero point -- the origin.)

		drawing also needs to know three rectangles, one for each image.
		Since each rectangle has the same width and height, it suffices to pass a destination rectangle r and two points sp and mp:
		the source rectangle is equal to r translated so that r.Min in the destination image aligns with sp in the source image, and similarly for mp.
	*/

	// func Draw(dst Image, r image.Rectangle, src image.Image, sp image.Point, op Op) {
	draw.Draw(dst, dst.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	outfile, err := os.Create("./data/src-over-dst.png")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, dst); err != nil {
		t.Fatal(err)
	}
}

func TestSrcOverDstTransparent(t *testing.T) {
	dst := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(dst, dst.Bounds(), image.Transparent, image.ZP, draw.Src)

	outfile, err := os.Create("./data/src-over-dst-transparent.png")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, dst); err != nil {
		t.Fatal(err)
	}
}

func readImage(t *testing.T, p string) draw.Image {
	sf, err := os.Open(p)
	if err != nil {
		t.Fatal(err)
	}
	dimg, _, err := image.Decode(sf)
	if err != nil {
		t.Fatal(err)
	}
	img, ok := dimg.(draw.Image)
	if !ok {
		t.Fatal(err)
	}
	return img
}

func readImageSimple(t *testing.T, p string) image.Image {
	sf, err := os.Open(p)
	if err != nil {
		t.Fatal(err)
	}
	dimg, _, err := image.Decode(sf)
	if err != nil {
		t.Fatal(err)
	}
	return dimg
}

func TestSrcImageOverDstImageSrc(t *testing.T) {
	dst := readImage(t, path.Join("data", "tenchi-mato-no-kamae-street.png"))
	src := readImage(t, path.Join("data", "tsuugaku_jitensya_girl_sailor.png"))
	draw.Draw(dst, src.Bounds(), src, image.ZP, draw.Src)
	outfile, err := os.Create("./data/src-src-dst-tenchi-mato-street.png")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, dst); err != nil {
		t.Fatal(err)
	}
}

func TestSrcImageOverDstImageOver(t *testing.T) {
	dst := readImage(t, path.Join("data", "tenchi-mato-no-kamae-street.png"))
	src := readImage(t, path.Join("data", "tsuugaku_jitensya_girl_sailor.png"))
	draw.Draw(dst, src.Bounds(), src, image.ZP, draw.Over)
	outfile, err := os.Create("./data/src-over-dst-tenchi-mato-street.png")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, dst); err != nil {
		t.Fatal(err)
	}
}

func TestSrcImageOverDstImageWithCenterCoordinates(t *testing.T) {
	dst := readImage(t, path.Join("data", "tenchi-mato-no-kamae-street.png"))
	src := readImage(t, path.Join("data", "tsuugaku_jitensya_girl_sailor.png"))
	dstRect := dst.Bounds()
	srcRect := src.Bounds()
	targetMin := image.Pt(dstRect.Max.X/2-srcRect.Max.X/2, dstRect.Max.Y/2-srcRect.Max.Y/2)
	targetMax := targetMin.Add(srcRect.Size())
	targetRect := image.Rectangle{Max: targetMax, Min: targetMin}
	draw.Draw(dst, targetRect, src, image.ZP, draw.Over)
	outfile, err := os.Create("./data/src-over-dst-tenchi-mato-street-center.png")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, dst); err != nil {
		t.Fatal(err)
	}
}

func TestSrcImageOverDstImageWithRightCoordinates1(t *testing.T) {
	dst := readImage(t, path.Join("data", "tenchi-mato-no-kamae-street.png"))
	src := readImage(t, path.Join("data", "tsuugaku_jitensya_girl_sailor.png"))
	dstRect := dst.Bounds()
	srcRect := src.Bounds()
	offset := 220
	targetMin := image.Pt((dstRect.Max.X/2-srcRect.Max.X/2)+offset, (dstRect.Max.Y/2-srcRect.Max.Y/2)+offset)
	// 2点を出すには以下でもいけるが、srcの画像全てを書き込むのであればSize()を利用した方が短く意図を伝えやすいと思う
	targetMax := image.Pt((dstRect.Max.X/2+srcRect.Max.X/2)+offset, (dstRect.Max.Y/2+srcRect.Max.Y/2)+offset)
	// targetMax := targetMin.Add(srcRect.Size())
	targetRect := image.Rectangle{Max: targetMax, Min: targetMin}
	t.Logf("destRect=%+v, srcRect=%+v, targetRect=%+v", dstRect, srcRect, targetRect)
	draw.Draw(dst, targetRect, src, image.ZP, draw.Over)
	outfile, err := os.Create("./data/src-over-dst-tenchi-mato-street-center-offset.png")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, dst); err != nil {
		t.Fatal(err)
	}
}

func TestSrcImageOverDstImageWithRightCoordinates2(t *testing.T) {
	dst := readImage(t, path.Join("data", "tenchi-mato-no-kamae-street.png"))
	src := readImage(t, path.Join("data", "tsuugaku_jitensya_girl_sailor.png"))
	dstRect := dst.Bounds()
	srcRect := src.Bounds()
	offset := 220
	targetMin := image.Pt((dstRect.Max.X/2-srcRect.Max.X/2)+offset, (dstRect.Max.Y/2-srcRect.Max.Y/2)+offset)
	// 2点を出すには以下でもいけるが、srcの画像全てを書き込むのであればSize()を利用した方が短く意図を伝えやすいと思う
	// targetMax := image.Pt((dstRect.Max.X/2+srcRect.Max.X/2)+offset, (dstRect.Max.Y/2+srcRect.Max.Y/2)+offset)
	targetMax := targetMin.Add(srcRect.Size())
	targetRect := image.Rectangle{Max: targetMax, Min: targetMin}
	t.Logf("destRect=%+v, srcRect=%+v, targetRect=%+v", dstRect, srcRect, targetRect)
	draw.Draw(dst, targetRect, src, image.ZP, draw.Over)
	outfile, err := os.Create("./data/src-over-dst-tenchi-mato-street-center-offset.png")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, dst); err != nil {
		t.Fatal(err)
	}
}

func TestSrcImageOverDstImageWithCoordinatesMultipleSource(t *testing.T) {
	dst := readImage(t, path.Join("data", "tenchi-mato-no-kamae-street.png"))
	src := readImage(t, path.Join("data", "tsuugaku_jitensya_girl_sailor.png"))
	dog := readImage(t, path.Join("data", "dog_corgi.png"))
	dstRect := dst.Bounds()
	srcRect := src.Bounds()
	offset := 220
	targetMin := image.Pt((dstRect.Max.X/2-srcRect.Max.X/2)+offset, (dstRect.Max.Y/2-srcRect.Max.Y/2)+offset)
	targetMax := targetMin.Add(srcRect.Size())
	targetRect := image.Rectangle{Max: targetMax, Min: targetMin}
	t.Logf("destRect=%+v, srcRect=%+v, targetRect=%+v", dstRect, srcRect, targetRect)
	draw.Draw(dst, targetRect, src, image.ZP, draw.Over)

	dogMin := image.Pt(300, 600)
	dogMax := dogMin.Add(dog.Bounds().Size())
	dogTarget := image.Rectangle{Max: dogMax, Min: dogMin}
	draw.Draw(dst, dogTarget, dog, image.ZP, draw.Over)
	outfile, err := os.Create("./data/src-over-dst-tenchi-mato-street-multiple-source.png")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, dst); err != nil {
		t.Fatal(err)
	}
}

func TestSrcImageOverDstImagePartialSrc(t *testing.T) {
	dst := readImage(t, path.Join("data", "tenchi-mato-no-kamae-street.png"))
	src := readImage(t, path.Join("data", "tsuugaku_jitensya_girl_sailor.png"))
	sb, db := src.Bounds(), dst.Bounds()
	offset := 220
	// dstに書き込む場所をsrcイメージの1/4に設定し、srcの1/4をそこに描画する
	targetMin := image.Pt((db.Max.X/2-sb.Max.X/4)+offset, (db.Max.Y/2-sb.Max.Y/4)+offset)
	targetMax := targetMin.Add(image.Pt(sb.Max.X/2, sb.Max.Y/2))
	targetRect := image.Rectangle{Max: targetMax, Min: targetMin}
	draw.Draw(dst, targetRect, src, image.ZP, draw.Over)
	outfile, err := os.Create("./data/src-over-dst-tenchi-mato-street-partial-src.png")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, dst); err != nil {
		t.Fatal(err)
	}
}

func TestSrcImageOverDstImageWithMaskOpaque(t *testing.T) {
	dst := readImage(t, path.Join("data", "tenchi-mato-no-kamae-street.png"))
	src := readImage(t, path.Join("data", "tsuugaku_jitensya_girl_sailor.png"))
	mask := readImage(t, path.Join("data", "mask-opaque.png"))
	offset := 220
	targetMin := image.Pt(
		(dst.Bounds().Max.X/2-src.Bounds().Max.X/2)+offset,
		(dst.Bounds().Max.Y/2-src.Bounds().Max.Y/2)+offset,
	)
	targetMax := targetMin.Add(src.Bounds().Size())
	targetRect := image.Rectangle{Max: targetMax, Min: targetMin}
	draw.DrawMask(dst, targetRect, src, image.ZP, mask, image.ZP, draw.Over)
	outfile, err := os.Create("./data/src-over-dst-tenchi-mato-street-mask-opaque.png")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, dst); err != nil {
		t.Fatal(err)
	}
}

// from https://go.dev/blog/image-draw
type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

func TestSrcImageOverDstImageWithMaskCircle(t *testing.T) {
	dst := readImage(t, path.Join("data", "tenchi-mato-no-kamae-street.png"))
	src := readImage(t, path.Join("data", "tsuugaku_jitensya_girl_sailor.png"))
	sb := src.Bounds()
	offset := 220
	targetMin := image.Pt(
		(dst.Bounds().Max.X/2-src.Bounds().Max.X/2)+offset,
		(dst.Bounds().Max.Y/2-src.Bounds().Max.Y/2)+offset,
	)
	targetMax := targetMin.Add(src.Bounds().Size())
	targetRect := image.Rectangle{Max: targetMax, Min: targetMin}
	c := circle{p: image.Pt(sb.Bounds().Max.X/2, sb.Bounds().Max.Y/2), r: 150}
	draw.DrawMask(dst, targetRect, src, image.ZP, &c, image.ZP, draw.Over)
	outfile, err := os.Create("./data/src-over-dst-tenchi-mato-street-mask-circle.png")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, dst); err != nil {
		t.Fatal(err)
	}
}

type fakeFontMask struct {
	i image.Image
}

func (f *fakeFontMask) ColorModel() color.Model {
	return color.AlphaModel
}

func (f *fakeFontMask) Bounds() image.Rectangle {
	return f.i.Bounds()
}

func (f *fakeFontMask) At(x, y int) color.Color {
	p := f.i.At(x, y)
	r, g, b, _ := p.RGBA()
	if r == g && g == b && r == uint32(0) {
		return color.Alpha{0}
	}
	return color.Alpha{255}
}

func TestSrcImageOverDstImageWithFakeFontMask(t *testing.T) {
	dst := readImage(t, path.Join("data", "tenchi-mato-no-kamae-street.png"))
	src := &image.Uniform{color.RGBA{0, 0, 255, 255}}
	m := readImageSimple(t, path.Join("data", "fake-font-mask-brack-white.png"))

	mask := &fakeFontMask{m}
	targetMin := image.Pt(600, 700)
	targetMax := targetMin.Add(image.Pt(50, 50))
	targetRect := image.Rectangle{Max: targetMax, Min: targetMin}
	maskp := image.Pt(10, 10)
	draw.DrawMask(dst, targetRect, src, image.ZP, mask, maskp, draw.Over)

	targetMin2 := targetMin.Add(image.Pt(100, 100))
	targetMax2 := targetMin2.Add(image.Pt(50, 50))
	targetRect2 := image.Rectangle{Max: targetMax2, Min: targetMin2}
	maskp2 := image.Pt(55, 10)
	draw.DrawMask(dst, targetRect2, src, image.ZP, mask, maskp2, draw.Over)
	outfile, err := os.Create("./data/src-over-dst-tenchi-mato-street-fake-font-mask.png")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, dst); err != nil {
		t.Fatal(err)
	}
}

func readFont(t *testing.T, p string) *truetype.Font {
	fontBytes, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal(err)
	}
	ft, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
	return ft
}

func writeImage(t *testing.T, p string, img image.Image) {
	outfile, err := os.Create(path.Join("data", p))
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, img); err != nil {
		log.Fatal(err)
	}
}

func TestDrawString(t *testing.T) {
	ff := "./Koruri-Bold.ttf"
	ft := readFont(t, ff)

	width, height, fontsize := 1200, 630, 100.0
	fontColor := image.Black
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	face := truetype.NewFace(ft, &truetype.Options{
		Size: fontsize,
	})
	dr := &font.Drawer{
		Dst:  img,
		Src:  fontColor,
		Face: face,
		Dot:  fixed.Point26_6{},
	}
	text := "あ"
	dr.Dot.X = (fixed.I(width) - dr.MeasureString(text)) / 2
	dr.Dot.Y = fixed.I((height + int(fontsize)/2) / 2)
	dr.DrawString(text)
	writeImage(t, "a.png", img)
}

func TestDrawLongString(t *testing.T) {
	ff := "./Koruri-Bold.ttf"
	ft := readFont(t, ff)

	width, height, fontsize := 1200, 630, 100.0
	fontColor := image.Black
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)
	face := truetype.NewFace(ft, &truetype.Options{
		Size: fontsize,
	})
	dr := &font.Drawer{
		Dst:  img,
		Src:  fontColor,
		Face: face,
		Dot:  fixed.Point26_6{},
	}
	text := "株式会社カンムはGoエンジニアを募集しております！一緒に良いプロダクトを作りませんか！"
	dr.Dot.X = fixed.I(100)
	dr.Dot.Y = fixed.I(200)
	dr.DrawString(text)
	writeImage(t, "we-are-hiring.png", img)
}

// DrawStringOpts options
type DrawStringOpts struct {
	ImageWidth       fixed.Int26_6
	ImageHeight      fixed.Int26_6
	VerticalMargin   fixed.Int26_6
	HorizontalMargin fixed.Int26_6
	FontSize         fixed.Int26_6
	LineSpace        fixed.Int26_6
	Verbose          bool
}

// DrawStringWrapped draw string wrapped
func DrawStringWrapped(d *font.Drawer, s string, opt *DrawStringOpts) {
	d.Dot.X = opt.HorizontalMargin
	d.Dot.Y = opt.FontSize + opt.VerticalMargin
	originalX, originalY := d.Dot.X, d.Dot.Y

	prevC := rune(-1)
	for _, c := range s {
		if prevC >= 0 {
			d.Dot.X += d.Face.Kern(prevC, c)
		}
		advance, ok := d.Face.GlyphAdvance(c)
		if !ok {
			// TODO: is falling back on the U+FFFD glyph the responsibility of
			// the Drawer or the Face?
			// TODO: set prevC = '\ufffd'?
			continue
		}
		if d.Dot.X+advance >= (opt.ImageWidth - opt.HorizontalMargin*2) {
			if opt.Verbose {
				fmt.Printf("### new line: %#U, x=%d, y=%d, ", c, d.Dot.X, d.Dot.Y)
			}
			d.Dot.Y = originalY + d.Dot.Y + opt.LineSpace
			d.Dot.X = originalX
		}
		dr, mask, maskp, advance, _ := d.Face.Glyph(d.Dot, c)
		if opt.Verbose {
			fmt.Printf(
				"%#U: maskp=%+v, advance=%d, X=%d, w=%d, Y=%d, h=%d, realW=%d\n",
				c, maskp, advance, d.Dot.X, opt.ImageWidth, d.Dot.Y, opt.ImageHeight, (opt.ImageWidth - opt.HorizontalMargin*2))
		}
		draw.DrawMask(d.Dst, dr, d.Src, image.Point{}, mask, maskp, draw.Over)
		d.Dot.X += advance
		prevC = c
	}
}

func TestDrawLongWrap(t *testing.T) {
	ff := "./Koruri-Bold.ttf"
	ft := readFont(t, ff)

	width, height, fontsize := 1200, 630, 100.0
	fontColor := image.Black
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)
	face := truetype.NewFace(ft, &truetype.Options{
		Size: fontsize,
	})
	dr := &font.Drawer{
		Dst:  img,
		Src:  fontColor,
		Face: face,
		Dot:  fixed.Point26_6{},
	}
	dOpt := &DrawStringOpts{
		ImageWidth:       fixed.I(width),
		ImageHeight:      fixed.I(height),
		Verbose:          false,
		FontSize:         fixed.I(int(fontsize)),
		LineSpace:        fixed.I(5),
		VerticalMargin:   fixed.I(10),
		HorizontalMargin: fixed.I(40),
	}

	text := "株式会社カンムはGoエンジニアを募集しております！一緒に良いプロダクトを作りませんか！"
	DrawStringWrapped(dr, text, dOpt)
	writeImage(t, "we-are-hiring-ok.png", img)
}
