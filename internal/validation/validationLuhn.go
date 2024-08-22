package validation

import "strconv"

// ValidateLuhn - checks order for right type
func ValidateLuhn(orderNumber string) bool {
	sum := 0
	isSecondDigit := false

	for i := len(orderNumber) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(orderNumber[i]))
		if err != nil {
			return false
		}

		if isSecondDigit {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isSecondDigit = !isSecondDigit
	}

	return sum%10 == 0
}
