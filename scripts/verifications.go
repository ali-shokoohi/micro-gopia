package scripts

import "unicode"

// VerifyPassword check passwords must be bigger than 8 character with uppercase, lowercase, digits and symbols
func VerifyPassword(s string) bool {
	letters := 0
	eightOrMore, number, upper, special := false, false, false, false
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
		}
	}
	eightOrMore = letters >= 8
	if eightOrMore && number && upper && special {
		return true
	}
	return false
}
