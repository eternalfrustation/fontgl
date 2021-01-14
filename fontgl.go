package fontgl

import (
	"bytes"
	"fmt"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"io"
	"os"
)

type Fnt struct {
	Glyphsegmap map[rune]sfnt.Segments
	buffer      *sfnt.Buffer
	Font        *sfnt.Font
}

func Setup(path string) *Fnt {
	Font := new(Fnt)
	Font.Glyphsegmap = make(map[rune]sfnt.Segments)
	fontfile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	fontbuffer := new(bytes.Buffer)
	io.Copy(fontbuffer, fontfile)
	fmt.Println(len(fontbuffer.Bytes()))
	font, err := sfnt.Parse(fontbuffer.Bytes())
	if err != nil {
		panic(err)
	}
	gbuf := new(sfnt.Buffer)
	for i := 0; i < font.NumGlyphs(); i++ {
		gi, err := font.GlyphIndex(gbuf, rune(i))
		//fmt.Println(gi)
		if err != nil {
			panic(err)
		}
		tempglyph, err := font.LoadGlyph(gbuf, gi, fixed.I(16), nil)
		if err != nil {
			panic(err)
		}
		Font.Glyphsegmap[rune(i)] = tempglyph
	}
	Font.buffer = gbuf
	Font.Font = font
	return Font
}

// Recieves a font from Setup(), string and x,y coords and gives an slice of float 32 in the format of
// {X, Y, 0, R, G, B, 0, 0, 1} for use with openGL
// R, G, B, are currently always 1
func GetTriText(todraw string, font *Fnt, x, y float32) (LArr, PtArr []float32) {
	//	var LArr, TriArr, PtArr []float32
	//	var prevMove bool
	var WeAtX, WeAtY float32
	for _, curRune := range todraw {
		glyph := font.Glyphsegmap[curRune]
		gi, _ := font.Font.GlyphIndex(font.buffer, curRune)
		Bounds, Advance26, err := font.Font.GlyphBounds(font.buffer, gi, fixed.I(16), 2)
		//fmt.Println(gi)
		if err != nil {
			panic(err)
		}
		maxX := float32(
			Bounds.
				Sub(Bounds.Min).
				Max.X)
		maxY := float32(
			Bounds.
				Sub(Bounds.Min).
				Max.Y)
		x += float32(
			Advance26) / (maxX)
		//fmt.Println(maxX, Advance26, x)
		for _, curSeg := range glyph {
			var MapX []float32
			var MapY []float32
			for _, arg := range curSeg.Args {
				mX := (float32(arg.X) / maxX) - 1
				mY := (float32(arg.Y) / maxY) - 1
				fmt.Println(mX, mY)
				MapX = append(MapX, mX)
				MapY = append(MapY, mY)
			}
			switch curSeg.Op {
			case sfnt.SegmentOpMoveTo:
				WeAtX = MapX[0]
				WeAtY = MapY[0]
			case sfnt.SegmentOpLineTo:
				LArr = append(LArr, WeAtX, WeAtY, 0, 1, 1, 1, 0, 0, 1)
				LArr = append(LArr, MapX[0], MapY[0], 0, 1, 1, 1, 0, 0, 1)
				fmt.Println("Line")
			case sfnt.SegmentOpQuadTo:
				LArr = append(LArr, WeAtX, WeAtY, 0, 1, 1, 1, 0, 0, 1)
				LArr = append(LArr, BezCurve(10, WeAtX, WeAtY, MapX[0], MapY[0], MapX[1], MapY[1])...)
			case sfnt.SegmentOpCubeTo:
				LArr = append(LArr, WeAtX, WeAtY, 0, 1, 1, 1, 0, 0, 1)
				LArr = append(LArr, BezCurve(10, WeAtX, WeAtY, MapX[0], MapY[0], MapX[1], MapY[1], MapX[2], MapY[2])...)
			}
		}
	}
	//	fmt.Println(len(LArr), len(PtArr))
	return LArr, PtArr
}

func BezCurve(t int, xs ...float32) []float32 {
	var Xa []float32
	for i := float32(0); int(i) < t; i += 1 / float32(t) {
		Xa = append(Xa, float32((1-i)*(1-i))*xs[0]+2*(1-i)*i*xs[2]+i*i*xs[4], float32((1-i)*(1-i))*xs[1]+2*(1-i)*i*xs[3]+i*i*xs[5], 0, 1, 1, 1, 0, 0, 1)
	}
	return Xa
}

func CubBezCurve(t int, xs ...float32) []float32 {
	var Xa []float32
	for j := 0; j < t; j += 1 {
		i := float32(j)/float32(t)
		Xa = append(Xa, float32((1-i)*(1-i)*(1-i))*xs[0]+3*(1-i)*(1-i)*i*xs[2]+3*(1-i)*i*i*xs[4]+i*i*i*xs[6], float32((1-i)*(1-i)*(1-i))*xs[1]+3*(1-i)*(1-i)*i*xs[3]+3*(1-i)*i*i*xs[5]+i*i*i*xs[7], 0, 1, 1, 1, 0, 0, 1)
	}
	return Xa
}
