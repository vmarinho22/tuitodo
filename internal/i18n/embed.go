package i18n

import "embed"

//go:embed locales/*/messages.json locales/*/formats.json
var localeFS embed.FS
