package helper

import (
	"net/mail"
	"strings"
)

//Return the string between the argument "a" and argument "b" in a string value.
func Between(value string, a string, b string) string {
    // Get substring between two strings.
    posFirst := strings.Index(value, a)
    if posFirst == -1 {
        return ""
    }
    posLast := strings.Index(value, b)
    if posLast == -1 {
        return ""
    }
    posFirstAdjusted := posFirst + len(a)
    if posFirstAdjusted >= posLast {
        return ""
    }
    return value[posFirstAdjusted:posLast]
}

//Checks if slice contains the given value
func Contains(slice []string, val string) bool {
	for _, elem := range slice {
		if elem == val {
			return true
		}
	}
	return false
}

//Checks the validity of the email address
func Valid_email(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}