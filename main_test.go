package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path"
	"testing"
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

func TestSrcImageOverDstImage(t *testing.T) {
	dst := readImage(t, path.Join("data", "tenchi-mato-no-kamae-street.png"))
	src := readImage(t, path.Join("data", "tsuugaku_jitensya_girl_sailor.png"))
	draw.Draw(dst, src.Bounds(), src, image.ZP, draw.Src)
	outfile, err := os.Create("./data/src-over-dst-tenchi-mato-street.png")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, dst); err != nil {
		t.Fatal(err)
	}
}
