package i18n

import (
	"encoding/json"
	"fmt"
	"path"
)

const localesDir = "locales"

func init() {
	catalogs, localeFormats = loadBundles()
}

func loadBundles() (map[Locale]map[string]string, map[Locale]formats) {
	messages := make(map[Locale]map[string]string)
	loadedFormats := make(map[Locale]formats)
	for _, locale := range Locales() {
		messages[locale] = mustLoadMessages(locale)
		loadedFormats[locale] = mustLoadFormats(locale)
	}
	return messages, loadedFormats
}

func mustLoadMessages(locale Locale) map[string]string {
	var nested map[string]any
	mustDecodeLocaleFile(locale, "messages.json", &nested)
	return flattenMessages(nested, "")
}

func flattenMessages(raw map[string]any, prefix string) map[string]string {
	out := make(map[string]string)
	flattenInto(raw, prefix, out)
	return out
}

func flattenInto(raw map[string]any, prefix string, out map[string]string) {
	for key, value := range raw {
		full := key
		if prefix != "" {
			full = prefix + "." + key
		}
		switch typed := value.(type) {
		case string:
			out[full] = typed
		case map[string]any:
			flattenInto(typed, full, out)
		default:
			panic(fmt.Sprintf("unsupported message value at %s", full))
		}
	}
}

func mustLoadFormats(locale Locale) formats {
	var loaded formats
	mustDecodeLocaleFile(locale, "formats.json", &loaded)
	return loaded
}

func mustDecodeLocaleFile(locale Locale, name string, dest any) {
	filePath := path.Join(localesDir, string(locale), name)
	data, err := localeFS.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("load %s: %v", filePath, err))
	}
	if err := json.Unmarshal(data, dest); err != nil {
		panic(fmt.Sprintf("parse %s: %v", filePath, err))
	}
}
