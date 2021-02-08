package blocklist

import (
	"errors"
	"strings"
)

type SpamhauseErrorCode error

var (
	SPLData                  SpamhauseErrorCode = errors.New(" SPL Zone: SPL Data")
	SPLCSSDAta               SpamhauseErrorCode = errors.New("SPL Zone: SPL CSS Data")
	CBLData                  SpamhauseErrorCode = errors.New("XBL Zone: CBL data")
	DROPEDROPData            SpamhauseErrorCode = errors.New("SBL Zone: Drop/EDROP Data")
	ISPMaintained            SpamhauseErrorCode = errors.New("PBL Zone: ISP Maintained")
	SpamhauseMaintained      SpamhauseErrorCode = errors.New("PBL Zone: Spamhause maintained")
	SpamhauseDomainBlockList SpamhauseErrorCode = errors.New("API does not support Domain Block List")
	TypingError              SpamhauseErrorCode = errors.New("DNSBL typing error")
	PublicResolver           SpamhauseErrorCode = errors.New("public resolver")
	ExcessNumberOfQueries    SpamhauseErrorCode = errors.New("excessive number of queries")
	UnsupportedError         SpamhauseErrorCode = errors.New("unknown error")
)

func (bl *BlockList) lookError(errorCode string) SpamhauseErrorCode {
	var returnError SpamhauseErrorCode

	if strings.Contains(errorCode, "127.0.1") {
		returnError = SpamhauseDomainBlockList
	}

	switch errorCode {
	case "127.0.0.2":
		returnError = SPLData
	case "127.0.0.3":
		returnError = SPLCSSDAta
	case "127.0.0.4":
		returnError = CBLData
	case "127.0.0.9":
		returnError = DROPEDROPData
	case "127.0.0.10":
		returnError = ISPMaintained
	case "127.0.0.11":
		returnError = SpamhauseMaintained
	case "127.255.255.252":
		returnError = TypingError
	case "127.255.255.254":
		returnError = PublicResolver
	case "127.255.255.255":
		returnError = ExcessNumberOfQueries

	default:
		returnError = UnsupportedError

	}

	return returnError
}
