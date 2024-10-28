package nlp

// missingCharFunc removes a character one at a time from the string.
// For example, wwwgoogle.com and www.googlecom
func MissingCharFunc(str, character string) (results []string) {
	for i, char := range str {
		if character == string(char) {
			results = append(results, str[:i]+str[i+1:])
		}
	}
	return
}