package ui

import "testing"

func TestLookupThemeBuiltins(t *testing.T) {
	cases := map[string]string{
		"ayu-light": "#FAFAFA",
		"min-light": "#ffffff",
	}
	for name, bg := range cases {
		th, err := LookupTheme(name)
		if err != nil {
			t.Fatalf("LookupTheme(%q): %v", name, err)
		}
		if th.BG != bg {
			t.Errorf("LookupTheme(%q).BG = %q, want %q", name, th.BG, bg)
		}
	}
}

func TestLookupThemeUnknown(t *testing.T) {
	_, err := LookupTheme("nope")
	if err == nil {
		t.Fatal("LookupTheme(\"nope\") returned no error")
	}
	want := `unknown theme "nope" (available: ayu-light, min-light)`
	if err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

func TestResolveThemeName(t *testing.T) {
	cases := []struct {
		flag, env, want string
	}{
		{"ayu-light", "min-light", "ayu-light"},
		{"", "ayu-light", "ayu-light"},
		{"", "", "min-light"},
	}
	for _, c := range cases {
		if got := ResolveThemeName(c.flag, c.env); got != c.want {
			t.Errorf("ResolveThemeName(%q, %q) = %q, want %q", c.flag, c.env, got, c.want)
		}
	}
}

func TestDefaultThemeIsBuiltin(t *testing.T) {
	if _, err := LookupTheme(DefaultThemeName); err != nil {
		t.Fatalf("DefaultThemeName %q is not a builtin: %v", DefaultThemeName, err)
	}
}
