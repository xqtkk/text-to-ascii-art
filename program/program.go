package program

import (
	"bufio"
	"embed"
	"fmt"
	"strings"
)

//go:embed banners/*.txt
var bannerFiles embed.FS

var allowedBanners = map[string]string{
	"standard":   "standard.txt",
	"shadow":     "shadow.txt",
	"thinkertoy": "thinkertoy.txt",
}

// Split divides text into lines while preserving intentional empty lines.
func Split(text string) []string {
	return strings.Split(text, "\n")
}

func GetFont(name string) (map[rune][]string, error) {
	fileName, ok := allowedBanners[name]
	if !ok {
		return nil, fmt.Errorf("unknown font %q", name)
	}

	file, err := bannerFiles.Open("banners/" + fileName)
	if err != nil {
		return nil, fmt.Errorf("open embedded font: %w", err)
	}
	defer file.Close()

	font := make(map[rune][]string, 95)
	scanner := bufio.NewScanner(file)
	lines := make([]string, 0, 9)
	character := rune(' ')
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) != 9 {
			continue
		}
		font[character] = append([]string(nil), lines[1:]...)
		character++
		lines = lines[:0]
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read embedded font: %w", err)
	}
	if len(font) != 95 {
		return nil, fmt.Errorf("font %q is incomplete", name)
	}
	return font, nil
}

func TextToASCIIArt(text, fontName string) (string, error) {
	font, err := GetFont(fontName)
	if err != nil {
		return "", err
	}
	for _, character := range text {
		if character == '\n' {
			continue
		}
		if character < ' ' || character > '~' {
			return "", fmt.Errorf("unsupported character %q: use printable ASCII characters", character)
		}
	}

	inputLines := Split(text)
	if strings.HasSuffix(text, "\n") {
		inputLines = inputLines[:len(inputLines)-1]
	}
	rows := make([]string, 0, len(inputLines)*8)
	for _, inputLine := range inputLines {
		for row := 0; row < 8; row++ {
			var output strings.Builder
			for _, character := range inputLine {
				output.WriteString(font[character][row])
			}
			rows = append(rows, output.String())
		}
	}
	return strings.Join(rows, "\n") + "\n", nil
}
