package twin

import (
	"fmt"
	"sort"
	"unicode"

	"github.com/rivo/uniseg"
)

// StyledRune is a rune with a style to be written to a one or more cells on the
// screen. Note that a StyledRune may use more than one cell on the screen ('午'
// for example).
type StyledRune struct {
	Rune  rune
	Style Style
}

func NewStyledRune(char rune, style Style) StyledRune {
	return StyledRune{
		Rune:  char,
		Style: style,
	}
}

func (styledRune StyledRune) String() string {
	return fmt.Sprint("rune='", string(styledRune.Rune), "' ", styledRune.Style)
}

// How many screen cells will this rune cover? Most runes cover one, but some
// like '午' will cover two.
func (styledRune StyledRune) Width() int {
	return uniseg.StringWidth(string(styledRune.Rune))
}

func (styledRune StyledRune) Equal(other StyledRune) bool {
	return styledRune.Rune == other.Rune && styledRune.Style.Equal(other.Style)
}

// Returns a slice of cells with trailing whitespace cells removed
func TrimSpaceRight(runes []StyledRune) []StyledRune {
	for i := len(runes) - 1; i >= 0; i-- {
		cell := runes[i]
		if !unicode.IsSpace(cell.Rune) {
			return runes[0 : i+1]
		}

		// That was a space, keep looking
	}

	// All whitespace, return empty
	return []StyledRune{}
}

// Returns a slice of cells with leading whitespace cells removed
func TrimSpaceLeft(runes []StyledRune) []StyledRune {
	for i := range runes {
		cell := runes[i]
		if !unicode.IsSpace(cell.Rune) {
			return runes[i:]
		}

		// That was a space, keep looking
	}

	// All whitespace, return empty
	return []StyledRune{}
}

// Blocks added in Unicode 15.1 (2023), 16.0 (2024), and 17.0 (2025). Go's
// unicode package lags behind the latest Unicode release (15.0.0 as of Go
// 1.25), so unicode.IsPrint() does not yet recognize these. We let the
// terminal render any unassigned code points within these blocks as tofu
// rather than mask real characters with '?'.
//
// Must be sorted by `lo` ascending; entries must not overlap. Binary search
// (sort.Search below) relies on this invariant.
var unicodePost15PrintableRanges = []struct {
	lo, hi rune
}{
	{0x105C0, 0x105FF}, // Todhri (16.0)
	{0x10940, 0x1095F}, // Sidetic (17.0)
	{0x10D40, 0x10D8F}, // Garay (16.0)
	{0x11380, 0x113FF}, // Tulu-Tigalari (16.0)
	{0x11B60, 0x11B7F}, // Sharada Supplement (17.0)
	{0x11BC0, 0x11BFF}, // Sunuwar (16.0)
	{0x11DB0, 0x11DEF}, // Tolong Siki (17.0)
	{0x13460, 0x143FF}, // Egyptian Hieroglyphs Extended-A (16.0)
	{0x16100, 0x1613F}, // Gurung Khema (16.0)
	{0x16D40, 0x16D7F}, // Kirat Rai (16.0)
	{0x16EA0, 0x16EDF}, // Beria Erfe (17.0)
	{0x18D80, 0x18DFF}, // Tangut Components Supplement (17.0)
	{0x1CC00, 0x1CEFF}, // Symbols for Legacy Computing Supplement (16.0) + Misc Symbols Supplement (17.0)
	{0x1E5D0, 0x1E5FF}, // Ol Onal (16.0)
	{0x1E6C0, 0x1E6FF}, // Tai Yo (17.0)
	{0x2EBF0, 0x2EE5F}, // CJK Unified Ideographs Extension I (15.1)
	{0x323B0, 0x3347F}, // CJK Unified Ideographs Extension J (17.0)
}

func Printable(char rune) bool {
	if unicode.IsPrint(char) {
		return true
	}

	if unicode.Is(unicode.Co, char) {
		// Co == "Private Use": https://www.compart.com/en/unicode/category
		//
		// This space is used by Font Awesome, for "fa-battery-empty" for
		// example: https://fontawesome.com/v4/icon/battery-empty
		//
		// So we want to print these and let the rendering engine deal with
		// outputting them in a way that's helpful to the user.
		return true
	}

	if char == 0xa0 {
		// 0xa0 is a non-breaking space, which is printable, despite what
		// unicode.IsPrint() says.
		return true
	}

	i := sort.Search(len(unicodePost15PrintableRanges), func(i int) bool {
		return unicodePost15PrintableRanges[i].lo > char
	})
	if i > 0 && char <= unicodePost15PrintableRanges[i-1].hi {
		return true
	}

	return false
}
