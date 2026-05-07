package twin

import (
	"reflect"
	"testing"

	"gotest.tools/v3/assert"
)

func TestTrimSpaceRight(t *testing.T) {
	// Empty
	assert.Assert(t, reflect.DeepEqual(
		TrimSpaceRight(
			[]StyledRune{},
		),
		[]StyledRune{}))

	// Single non-space
	assert.Assert(t, reflect.DeepEqual(
		TrimSpaceRight(
			[]StyledRune{{Rune: 'x'}},
		),
		[]StyledRune{{Rune: 'x'}}))

	// Single space
	assert.Assert(t, reflect.DeepEqual(
		TrimSpaceRight(
			[]StyledRune{{Rune: ' '}},
		),
		[]StyledRune{}))

	// Non-space plus space
	assert.Assert(t, reflect.DeepEqual(
		TrimSpaceRight(
			[]StyledRune{{Rune: 'x'}, {Rune: ' '}},
		),
		[]StyledRune{{Rune: 'x'}}))
}

func TestRuneWidth(t *testing.T) {
	assert.Equal(t, NewStyledRune('x', Style{}).Width(), 1)
	assert.Equal(t, NewStyledRune('午', Style{}).Width(), 2)
}

// Go's unicode tables (15.0.0 as of Go 1.25) lag behind the latest Unicode
// release. These are blocks added in Unicode 15.1 and 16.0 that
// unicode.IsPrint() does not yet recognize.
func TestPrintableUnicodePost15(t *testing.T) {
	cases := []struct {
		name string
		r    rune
	}{
		// Unicode 15.1 (2023)
		{"CJK Ext I start", 0x2EBF0},
		{"CJK Ext I end", 0x2EE5F},

		// Unicode 16.0 (2024)
		{"Todhri start", 0x105C0},
		{"Todhri end", 0x105F3},
		{"Garay start", 0x10D40},
		{"Garay end", 0x10D8E},
		{"Tulu-Tigalari start", 0x11380},
		{"Tulu-Tigalari end", 0x113D5},
		{"Sunuwar start", 0x11BC0},
		{"Sunuwar end", 0x11BF2},
		{"Egyptian Hieroglyphs Ext-A start", 0x13460},
		{"Egyptian Hieroglyphs Ext-A end", 0x143FA},
		{"Gurung Khema start", 0x16100},
		{"Gurung Khema end", 0x16139},
		{"Kirat Rai start", 0x16D40},
		{"Kirat Rai end", 0x16D79},
		{"Legacy Computing Supplement start", 0x1CC00},
		{"Large Type Piece (used by jj)", 0x1CE1A},
		{"Large Type Piece end", 0x1CE50},
		{"Legacy Computing Supplement end", 0x1CEBF},
		{"Ol Onal start", 0x1E5D0},
		{"Ol Onal end", 0x1E5FA},

		// Unicode 17.0 (2025)
		{"Sidetic start", 0x10940},
		{"Sidetic end", 0x1095F},
		{"Sharada Supplement start", 0x11B60},
		{"Sharada Supplement end", 0x11B7F},
		{"Tolong Siki start", 0x11DB0},
		{"Tolong Siki end", 0x11DEF},
		{"Beria Erfe start", 0x16EA0},
		{"Beria Erfe end", 0x16EDF},
		{"Tangut Components Supplement start", 0x18D80},
		{"Tangut Components Supplement end", 0x18DFF},
		{"Misc Symbols Supplement start", 0x1CEC0},
		{"Misc Symbols Supplement end", 0x1CEFF},
		{"Tai Yo start", 0x1E6C0},
		{"Tai Yo end", 0x1E6FF},
		{"CJK Ext J start", 0x323B0},
		{"CJK Ext J end", 0x3347F},
	}

	for _, tc := range cases {
		assert.Assert(t, Printable(tc.r),
			"expected U+%04X (%s) to be printable", tc.r, tc.name)
	}
}

// Mix of ASCII (the dominant case in real input), CJK, an emoji, an
// unprintable control char, and a Unicode 16+ rune that exercises the new
// range table.
// Binary search in Printable() depends on the table being sorted by `lo`
// with no overlaps. Catch ordering mistakes that the existing membership
// tests can miss (sort.Search returns 0 for an out-of-place leading entry,
// which silently misses lookups).
func TestUnicodePost15PrintableRangesSorted(t *testing.T) {
	prevHi := rune(-1)
	for _, r := range unicodePost15PrintableRanges {
		assert.Assert(t, r.lo > prevHi,
			"range %X..%X overlaps or is out of order with previous (hi=%X)",
			r.lo, r.hi, prevHi)
		assert.Assert(t, r.lo <= r.hi,
			"range %X..%X has lo > hi", r.lo, r.hi)
		prevHi = r.hi
	}
}

var benchPrintableInput = []rune{
	'a', 'b', 'c', ' ', '1', '\t', '\n', // ASCII / common
	'午',     // CJK
	'🚀',     // emoji
	0x07,    // BEL — unprintable
	0xa0,    // NBSP
	0x1CE1A, // Large Type Piece (Unicode 16, only printable via the new table)
}

func BenchmarkPrintable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, r := range benchPrintableInput {
			_ = Printable(r)
		}
	}
}
