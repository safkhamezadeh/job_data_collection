package adzuna

type Iso2CountryCode string

const (
	GB Iso2CountryCode = "GB"
	US Iso2CountryCode = "US"
	AT Iso2CountryCode = "AT"
	AU Iso2CountryCode = "AU"
	BE Iso2CountryCode = "BE"
	BR Iso2CountryCode = "BR"
	CA Iso2CountryCode = "CA"
	CH Iso2CountryCode = "CH"
	DE Iso2CountryCode = "DE"
	ES Iso2CountryCode = "ES"
	FR Iso2CountryCode = "FR"
	IN Iso2CountryCode = "IN"
	IT Iso2CountryCode = "IT"
	MX Iso2CountryCode = "MX"
	NL Iso2CountryCode = "NL"
	NZ Iso2CountryCode = "NZ"
	PL Iso2CountryCode = "PL"
	SG Iso2CountryCode = "SG"
	ZA Iso2CountryCode = "ZA"
)

// Map for fast O(1) lookups
var WHITELISTEDCOUNTRIES = map[Iso2CountryCode]struct{}{
	GB: {},
	US: {},
	AT: {},
	AU: {},
	BE: {},
	BR: {},
	CA: {},
	CH: {},
	DE: {},
	ES: {},
	FR: {},
	IN: {},
	IT: {},
	MX: {},
	NL: {},
	NZ: {},
	PL: {},
	SG: {},
	ZA: {},
}

// Helper function to check whitelist
func IsWhitelisted(country Iso2CountryCode) bool {
	_, ok := WHITELISTEDCOUNTRIES[country]
	return ok
}
