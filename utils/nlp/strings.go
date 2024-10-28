package nlp

import "strings"

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

// replaceCharFunc omits a character from the entire string.
// For example, www.a-b-c.com becomes www.abc.com
func ReplaceCharFunc(str, old, new string) (results []string) {
	results = append(results, strings.Replace(str, old, new, -1))
	return
}


func Paddding(str, old, new string) (results []string) {
	results = append(results, strings.Replace(str, old, new, -1))
	return
}
