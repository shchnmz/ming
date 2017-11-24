package ming

import (
	"regexp"
)

// Valid8DigitTelephoneNum checks if phone number matches the format:
// 1. Starts with 8 digital number.
// 2. Can have one or more '.' as sufix.
func Valid8DigitTelephoneNum(phoneNum string) bool {
	p := `^\d{8}\.*$`
	re := regexp.MustCompile(p)
	return re.MatchString(phoneNum)
}

// ValidMobilePhoneNum checks if phone number matches the format:
// 1. Starts with 11 digital number.
// 2. Can have one or more '.' as sufix.
func ValidMobilePhoneNum(phoneNum string) bool {
	p := `^\d{11}\.*$`
	re := regexp.MustCompile(p)
	return re.MatchString(phoneNum)
}

// ValidPhoneNum checks if phone number is 11-digit mobile phone number or 8-digit telephone number.
func ValidPhoneNum(phoneNum string) bool {
	if !Valid8DigitTelephoneNum(phoneNum) && !ValidMobilePhoneNum(phoneNum) {
		return false
	}
	return true
}
