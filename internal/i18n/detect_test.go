package i18n

import "testing"

func TestDetect(t *testing.T) {
	cases := []struct {
		name  string
		lcAll string
		lang  string
		want  Locale
	}{
		{name: "pt_BR from LANG", lang: "pt_BR.UTF-8", want: LocalePtBR},
		{name: "en_US from LANG", lang: "en_US.UTF-8", want: LocaleEnUS},
		{name: "french falls back", lang: "fr_FR.UTF-8", want: LocaleEnUS},
		{name: "empty falls back", want: LocaleEnUS},
		{name: "LC_ALL wins over LANG", lcAll: "pt_BR.UTF-8", lang: "en_US.UTF-8", want: LocalePtBR},
		{name: "C locale falls through to LANG", lcAll: "C", lang: "pt_PT", want: LocalePtBR},
		{name: "pt hyphenated", lang: "pt-BR", want: LocalePtBR},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got := Detect(testCase.lcAll, testCase.lang)
			if got != testCase.want {
				t.Fatalf("Detect(%q, %q) = %q, want %q", testCase.lcAll, testCase.lang, got, testCase.want)
			}
		})
	}
}

func TestResolvePrefersSavedLocale(t *testing.T) {
	got := Resolve("pt-BR", "en_US.UTF-8", "en_US.UTF-8")
	if got != LocalePtBR {
		t.Fatalf("Resolve saved = %q, want %q", got, LocalePtBR)
	}
}

func TestResolveInvalidSavedUsesDetection(t *testing.T) {
	got := Resolve("de-DE", "", "pt_BR.UTF-8")
	if got != LocalePtBR {
		t.Fatalf("Resolve invalid = %q, want %q", got, LocalePtBR)
	}
}

func TestParseLocale(t *testing.T) {
	if loc, ok := ParseLocale("en-US"); !ok || loc != LocaleEnUS {
		t.Fatalf("ParseLocale(en-US) = %q %v", loc, ok)
	}
	if _, ok := ParseLocale("en"); ok {
		t.Fatal("expected en to be invalid stored locale")
	}
}
