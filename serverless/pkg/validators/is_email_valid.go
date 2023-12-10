package validators 

import "regexp"

const email_min_length = 3
const email_max_length = 254

func IsEmailValid(email string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]{1,64}@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	// MustCompile method parses regular expression and if successful, returns a Regexp object that in this case, can be used to match against email

	if len(email) < email_min_length || len(email) > email_min_length || !rxEmail.MatchString(email){
		// if email length is less than min length or greater than max length or does not match regex, return false
		return false
	}
	
	return true
}