package i18n_app

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
)

func GetLocales(lang string) *i18n.Localizer {
	var esPath string
	var enPath string
	var bundle *i18n.Bundle

	//For testing purposes
	if flag.Lookup("test.v") == nil {
		esPath = "./i18n-app/catalog/active.es.toml"
		enPath = "./i18n-app/catalog/active.en.toml"
	} else {
		esPath = "../../i18n-app/catalog/active.es.toml"
		enPath = "../../i18n-app/catalog/active.en.toml"
	}

	switch lang {
	case "es-ES":
		bundle = i18n.NewBundle(language.Spanish)
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		_, err := bundle.LoadMessageFile(esPath)
		if err != nil {
			log.Println(err)
		}
	default:
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		_, err := bundle.LoadMessageFile(enPath)
		if err != nil {
			log.Println(err)
		}
	}
	return i18n.NewLocalizer(bundle, lang)
}
