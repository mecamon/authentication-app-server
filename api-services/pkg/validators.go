package pkg

import "regexp"

func ValidEmail(email string) bool {
	pattern := "^.{2,}@.{2,}\\..{2,}$"
	output, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}

	return output
}

func ValidPassword(password string) bool {
	regx1 := regexp.MustCompile("(.*[a-z])")
	regx2 := regexp.MustCompile("(.*[A-Z])")
	regx3 := regexp.MustCompile("(.*[\\d])")
	regx4 := regexp.MustCompile("^.{8,}$")

	output1 := regx1.MatchString(password)
	output2 := regx2.MatchString(password)
	output3 := regx3.MatchString(password)
	output4 := regx4.MatchString(password)

	return output1 && output2 && output3 && output4
}
