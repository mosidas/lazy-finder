package ui

import (
	"fmt"
	"sort"
	"strings"
)

// Theme is a palette of hex colours. The UI paints an explicit background so
// it keeps the theme's look regardless of the terminal's own background.
type Theme struct {
	BG     string // editor background
	BGSel  string // cursor row
	BGMark string // parent selection / status chip
	FG     string // primary text
	Accent string // focused border, marks
	Blue   string // directories, titles
	Cyan   string // symlinks
	Green  string // success
	Red    string // errors
	Muted  string // muted text
	Border string // inactive border
}

// DefaultThemeName is the theme used when none is chosen.
const DefaultThemeName = "min-light"

var themes = map[string]Theme{
	"ayu-light": {
		BG:     "#FAFAFA",
		BGSel:  "#D6E5F3",
		BGMark: "#E6E8EB",
		FG:     "#5C6166",
		Accent: "#FA8D3E",
		Blue:   "#399EE6",
		Cyan:   "#55B4D4",
		Green:  "#86B300",
		Red:    "#E65050",
		Muted:  "#8A9199",
		Border: "#C5C9CE",
	},
	// Min Light from https://github.com/miguelsolorio/min-theme
	"min-light": {
		BG:     "#ffffff",
		BGSel:  "#e0e0e0",
		BGMark: "#eeeeee",
		FG:     "#212121",
		Accent: "#424242",
		Blue:   "#6871ff",
		Cyan:   "#4dbf99",
		Green:  "#77cc00",
		Red:    "#d32f2f",
		Muted:  "#757575",
		Border: "#d0d0d0",
	},
}

// ThemeNames returns the builtin theme names in sorted order.
func ThemeNames() []string {
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// LookupTheme returns the builtin theme called name.
func LookupTheme(name string) (Theme, error) {
	t, ok := themes[name]
	if !ok {
		return Theme{}, fmt.Errorf("unknown theme %q (available: %s)", name, strings.Join(ThemeNames(), ", "))
	}
	return t, nil
}

// ResolveThemeName picks the theme name: the flag wins over the environment,
// and DefaultThemeName applies when both are empty.
func ResolveThemeName(flag, env string) string {
	if flag != "" {
		return flag
	}
	if env != "" {
		return env
	}
	return DefaultThemeName
}
