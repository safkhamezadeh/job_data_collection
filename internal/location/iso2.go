package location

import "strings"

type CountryISO2 string

var iso2Countries = map[CountryISO2]struct{}{
	"us": {}, "gb": {}, "de": {}, "fr": {}, "nl": {}, "be": {}, "es": {}, "it": {}, "pl": {},
	"se": {}, "no": {}, "fi": {}, "dk": {}, "ie": {}, "pt": {}, "cz": {}, "sk": {}, "hu": {},
	"ro": {}, "bg": {}, "hr": {}, "si": {}, "ee": {}, "lv": {}, "lt": {}, "gr": {},

	"ca": {}, "mx": {},

	"br": {}, "ar": {}, "cl": {}, "co": {}, "pe": {}, "ve": {},

	"cn": {}, "jp": {}, "kr": {}, "in": {}, "sg": {}, "my": {}, "th": {}, "id": {}, "ph": {},
	"vn": {}, "hk": {}, "tw": {},

	"au": {}, "nz": {},

	"za": {}, "eg": {}, "ng": {}, "ke": {}, "ma": {},

	"ae": {}, "sa": {}, "il": {}, "tr": {}, "qa": {}, "kw": {},
}

func IsValidISO2(country CountryISO2) bool {
	if country == "" {
		return false
	}

	c := CountryISO2(strings.ToLower(string(country)))
	_, exists := iso2Countries[c]
	return exists
}

func NormalizeISO2(country string) CountryISO2 {
	return CountryISO2(strings.ToLower(strings.TrimSpace(country)))
}
