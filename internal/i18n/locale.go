package i18n

import "strings"

type Locale string

const (
	LocaleEnUS Locale = "en-US"
	LocalePtBR Locale = "pt-BR"
)

func ParseLocale(value string) (Locale, bool) {
	switch Locale(value) {
	case LocaleEnUS, LocalePtBR:
		return Locale(value), true
	default:
		return "", false
	}
}

func Locales() []Locale {
	return []Locale{LocaleEnUS, LocalePtBR}
}

func (locale Locale) NativeLabel() string {
	switch locale {
	case LocalePtBR:
		return "Português (pt-BR)"
	default:
		return "English (en-US)"
	}
}

func Detect(lcAll, lang string) Locale {
	for _, value := range []string{lcAll, lang} {
		if locale, ok := localeFromEnv(value); ok {
			return locale
		}
	}
	return LocaleEnUS
}

func Resolve(saved, lcAll, lang string) Locale {
	if locale, ok := ParseLocale(saved); ok {
		return locale
	}
	return Detect(lcAll, lang)
}

func localeFromEnv(value string) (Locale, bool) {
	value = strings.TrimSpace(value)
	if value == "" || value == "C" || value == "POSIX" {
		return "", false
	}
	base, _, _ := strings.Cut(value, ".")
	base = strings.ReplaceAll(base, "-", "_")
	lower := strings.ToLower(base)
	switch {
	case strings.HasPrefix(lower, "pt"):
		return LocalePtBR, true
	case strings.HasPrefix(lower, "en"):
		return LocaleEnUS, true
	default:
		return "", false
	}
}
