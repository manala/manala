package asciify

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

const (
	ansiReset             = "\x1b[m"
	ansiCSI               = "\x1b["
	ansiRGBIntroducer     = 2
	ansiForegroundColor   = 38
	ansiBackgroundColor   = 48
	ansiDefaultBackground = 49

	upperHalfBlock = "▀"
	lowerHalfBlock = "▄"
)

// Asciify converts an NRGBA image into an ASCII art string using half-block characters with ANSI RGB colors.
// Each pair of vertical pixels is rendered as a single character with foreground and background colors.
// Transparent pixels are rendered as spaces. The output includes ANSI escape codes for terminal display.
func Asciify(src *image.NRGBA) string {
	pix := src.Pix
	stride := src.Stride
	bounds := src.Bounds()
	rows := bounds.Dy()
	cols := bounds.Dx()

	// Pre-size buffer: ~width chars per row/2, plus ANSI escapes overhead
	var buf strings.Builder
	buf.Grow(cols * rows / 2 * 24)

	var curFg, curBg color.NRGBA
	hasBg := false
	styled := false
	ansi := make([]byte, 0, 64)

	for y := 0; y < rows; y += 2 {
		topOff := y * stride
		botOff := (y + 1) * stride

		for x := range cols {
			px := x * 4
			topColor := color.NRGBA{R: pix[topOff+px], G: pix[topOff+px+1], B: pix[topOff+px+2], A: pix[topOff+px+3]}
			botColor := color.NRGBA{R: pix[botOff+px], G: pix[botOff+px+1], B: pix[botOff+px+2], A: pix[botOff+px+3]}
			topOpaque := topColor.A > 0
			botOpaque := botColor.A > 0

			// Determine desired fg, bg and character
			var wantFg, wantBg color.NRGBA
			var ch string
			switch {
			case topOpaque && botOpaque:
				wantFg, wantBg, ch = topColor, botColor, upperHalfBlock
			case topOpaque:
				wantFg, wantBg, ch = topColor, color.NRGBA{}, upperHalfBlock
			case botOpaque:
				wantFg, wantBg, ch = botColor, color.NRGBA{}, lowerHalfBlock
			default:
				if styled {
					buf.WriteString(ansiReset)
					curFg, curBg = color.NRGBA{}, color.NRGBA{}
					hasBg = false
					styled = false
				}
				buf.WriteRune(' ')

				continue
			}

			// Emit ANSI only when fg or bg changes
			if !styled || wantFg != curFg || wantBg != curBg {
				wantHasBg := wantBg.A > 0
				fgChanged := !styled || wantFg != curFg
				bgChanged := !styled || wantHasBg != hasBg || wantBg != curBg

				ansi = ansi[:0]
				if fgChanged {
					ansi = fmt.Appendf(ansi, "\x1b[%d;%d;%d;%d;%d",
						ansiForegroundColor, ansiRGBIntroducer, wantFg.R, wantFg.G, wantFg.B)
				}
				if bgChanged {
					if len(ansi) == 0 {
						ansi = append(ansi, ansiCSI...)
					} else {
						ansi = append(ansi, ';')
					}
					if wantHasBg {
						ansi = fmt.Appendf(ansi, "%d;%d;%d;%d;%d",
							ansiBackgroundColor, ansiRGBIntroducer, wantBg.R, wantBg.G, wantBg.B)
					} else {
						ansi = fmt.Appendf(ansi, "%d", ansiDefaultBackground)
					}
				}
				ansi = append(ansi, 'm')
				buf.Write(ansi)

				curFg, curBg = wantFg, wantBg
				hasBg = wantHasBg
				styled = true
			}
			buf.WriteString(ch)
		}

		buf.WriteRune('\n')
	}

	// Final reset if still styled
	if styled {
		buf.WriteString(ansiReset)
	}

	return buf.String()
}
