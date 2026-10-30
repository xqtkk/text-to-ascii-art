package program

import (
	"strings"
	"testing"
)

func TestGetFontIncludesAllPrintableASCII(t *testing.T) {
	for _, name := range []string{"standard", "shadow", "thinkertoy"} {
		font, err := GetFont(name)
		if err != nil {
			t.Fatalf("GetFont(%q): %v", name, err)
		}
		if len(font) != 95 {
			t.Errorf("GetFont(%q) returned %d glyphs, want 95", name, len(font))
		}
		for character := rune(' '); character <= '~'; character++ {
			if rows := font[character]; len(rows) != 8 {
				t.Errorf("font %q glyph %q has %d rows, want 8", name, character, len(rows))
			}
		}
	}
}

func TestTextToASCIIArtKeepsLineBreaksAndTrailingNewline(t *testing.T) {
	font, err := GetFont("standard")
	if err != nil {
		t.Fatal(err)
	}
	first := strings.Join(font['A'], "\n") + "\n"
	second := strings.Join(font['B'], "\n") + "\n"
	want := first + second

	got, err := TextToASCIIArt("A\nB", "standard")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("TextToASCIIArt output differs from banner glyphs\nwant: %q\n got: %q", want, got)
	}
}

func TestTextToASCIIArtRejectsUnsupportedCharactersAndFonts(t *testing.T) {
	if _, err := TextToASCIIArt("hello ☃", "standard"); err == nil {
		t.Fatal("expected unsupported character to return an error")
	}
	if _, err := TextToASCIIArt("hello", "../../README.md"); err == nil {
		t.Fatal("expected unknown font to return an error")
	}
}

func TestSplitPreservesEmptyLines(t *testing.T) {
	got := Split("first\n\nthird")
	want := []string{"first", "", "third"}
	if len(got) != len(want) {
		t.Fatalf("Split returned %d lines, want %d", len(got), len(want))
	}
	for index := range want {
		if got[index] != want[index] {
			t.Errorf("Split line %d = %q, want %q", index, got[index], want[index])
		}
	}
}
