package i18n

import (
	"errors"
	"fmt"
	"time"

	"tuitodo/internal/domain"
	"tuitodo/internal/store"
)

type formats struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

var (
	catalogs      map[Locale]map[string]string
	localeFormats map[Locale]formats
)

type Catalog struct {
	locale   Locale
	messages map[Locale]map[string]string
}

func NewCatalog(locale Locale) *Catalog {
	if _, ok := ParseLocale(string(locale)); !ok {
		locale = LocaleEnUS
	}
	return &Catalog{
		locale:   locale,
		messages: catalogs,
	}
}

func (catalog *Catalog) Locale() Locale {
	return catalog.locale
}

func (catalog *Catalog) SetLocale(locale Locale) {
	if _, ok := ParseLocale(string(locale)); !ok {
		locale = LocaleEnUS
	}
	catalog.locale = locale
}

func (catalog *Catalog) T(key string, args ...any) string {
	value, ok := catalog.lookup(catalog.locale, key)
	if !ok {
		value, ok = catalog.lookup(LocaleEnUS, key)
	}
	if !ok {
		return key
	}
	if len(args) == 0 {
		return value
	}
	return fmt.Sprintf(value, args...)
}

func (catalog *Catalog) lookup(locale Locale, key string) (string, bool) {
	if catalog.messages == nil {
		return "", false
	}
	bundle, ok := catalog.messages[locale]
	if !ok {
		return "", false
	}
	value, ok := bundle[key]
	return value, ok
}

func (catalog *Catalog) Error(err error) string {
	switch {
	case err == nil:
		return ""
	case errors.Is(err, domain.ErrEmptyTitle):
		return catalog.T(KeyEmptyTitle)
	case errors.Is(err, store.ErrCategoryInUse):
		return catalog.T(KeyCategoryInUse)
	case errors.Is(err, domain.ErrNestedSubtask):
		return catalog.T(KeyNestedSubtask)
	case errors.Is(err, domain.ErrNotAParent):
		return catalog.T(KeyNotAParent)
	case errors.Is(err, domain.ErrSubtaskNotFound):
		return catalog.T(KeySubtaskNotFound)
	default:
		return catalog.T(KeyUnexpected)
	}
}

func (catalog *Catalog) FormatDate(value time.Time) string {
	return value.Format(catalog.currentFormats().Date)
}

func (catalog *Catalog) FormatTime(value time.Time) string {
	return value.Format(catalog.currentFormats().Time)
}

func (catalog *Catalog) currentFormats() formats {
	if current, ok := localeFormats[catalog.locale]; ok {
		return current
	}
	return localeFormats[LocaleEnUS]
}
