package fontgl

import (
	"bytes"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"io"
	"os"
)

var (
	glyphsegmap   = make(map[rune]sfnt.Segments)
	glyphsegments []sfnt.Segments
)

func Setup(path string) {
	fontfile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	fontbuffer := new(bytes.Buffer)
	io.Copy(fontbuffer, fontfile)
	font, err := sfnt.Parse(fontbuffer.Bytes())
	if err != nil {
		panic(err)
	}
	for i := 0; i < font.NumGlyphs(); i++ {
		gbuf := new(sfnt.Buffer)
		tgbuf := new(sfnt.Buffer)
		gi, err := font.GlyphIndex(gbuf, rune(i))
		if err != nil {
			panic(err)
		}
		tempglyph, err := font.LoadGlyph(tgbuf, gi, fixed.Int26_6(0), nil)
		if err != nil {
			panic(err)
		}
		glyphsegmap[rune(i)] = tempglyph
	}
}

func GetTriText(todraw string, x, y float32) (LArr, TriArr, PtArr []float32) {
//	var LArr, TriArr, PtArr []float32
	for _, curRune := range todraw {
		glyph := glyphsegmap[curRune]
		x += float32(glyph.Bounds().Max.X)
		y += float32(glyph.Bounds().Max.Y)
		for _, curSeg := range glyph {
			var MapX []float32
			var MapY []float32
			for _, arg := range curSeg.Args {
				mX := float32(arg.X)/float32(glyph.Bounds().Max.X) + x
				mY := float32(arg.Y)/float32(glyph.Bounds().Max.Y) + y
				MapX = append(MapX, mX)
				MapY = append(MapY, mY)
			}
			switch curSeg.Op {
			case sfnt.SegmentOpMoveTo:
				LArr = append(LArr, MapX[0], MapY[0], 0, 1, 1, 1, 0, 0, 1)
				TriArr = append(TriArr, MapX[0], MapY[0], 0, 1, 1, 1, 0, 0, 1)
			case sfnt.SegmentOpLineTo:
				LArr = append(LArr, MapX[0], MapY[0], 0, 1, 1, 1, 0, 0, 1)
				TriArr = TriArr[:len(TriArr)-9]
			case sfnt.SegmentOpQuadTo:
				TriArr = TriArr[:len(TriArr)-9]
				p1x := LArr[len(LArr)-9]
				p1y := LArr[len(LArr)-8]
				LArr = LArr[:len(LArr)-9]
				PtArr = append(PtArr, BezCurve(20, p1x, p1y, MapX[0], MapY[0], MapX[1], MapY[1])...)
			case sfnt.SegmentOpCubeTo:
				TriArr = TriArr[:len(TriArr)-9]
				p1x := LArr[len(LArr)-9]
				p1y := LArr[len(LArr)-8]
				LArr = LArr[:len(LArr)-9]
				PtArr = append(PtArr, BezCurve(20, p1x, p1y, MapX[0], MapY[0], MapX[1], MapY[1], MapX[2], MapY[2])...)
			}
		}
	}
	return LArr, TriArr, PtArr
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
	for i := float32(0); int(i) < t; i += 1 / float32(t) {
		Xa = append(Xa, float32((1-i)*(1-i)*(1-i))*xs[0] + 3*(1-i)*(1-i)*i*xs[2] + 3*(1-i)*i*i*xs[4] + i*i*i*xs[6], float32((1-i)*(1-i)*(1-i))*xs[1] + 3*(1-i)*(1-i)*i*xs[3] + 3*(1-i)*i*i*xs[5] + i*i*i*xs[7], 0, 1, 1, 1, 0, 0, 1)
	}
	return Xa
}
