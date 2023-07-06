package requests

import (
    "regexp"
    "strings"
)

func HTMLhighlight(htmlCode string) string {
	// Define the regular expressions for highlighting
	tagPattern := regexp.MustCompile(`(<\w+|\w+>|<\/\w+>)`)
	attributePattern := regexp.MustCompile(`([\w-]+)(=)`)
	commentPattern := regexp.MustCompile(`<!--[\s\S]*?-->`)
    quotePattern := regexp.MustCompile(`(["][a-zA-Z0-9_//\ \)\(-;:.\{\}&?]*)(.*?)`)
    
	// Apply syntax highlighting
	htmlCode = tagPattern.ReplaceAllString(htmlCode, "\x1b[34m$1\x1b[0m")         // Blue 
	htmlCode = attributePattern.ReplaceAllString(htmlCode, "\x1b[31m$0\x1b[0m")  // Red 
    htmlCode = commentPattern.ReplaceAllString(htmlCode, "\x1b[32m$0\x1b[0m")   //  green
    htmlCode = quotePattern.ReplaceAllString(htmlCode, "\x1b[36m$0\x1b[0m")     // cyan
    
	return htmlCode
}



func JSONhighlight(jsonCode string) string {
   	var highlighted strings.Builder
	inString := false
	isKey := true

	for _, char := range jsonCode {
		switch char {
		case '{', '[':
			if !inString {
				highlighted.WriteString(Cyan + string(char) + Reset )
				
			} else {
				highlighted.WriteString(string(char))
			}
			isKey = true
		case '}', ']':
			if !inString {
				highlighted.WriteString(Cyan + string(char) + Reset)
			} else {
				highlighted.WriteString(string(char))
				
			}
			isKey = false
		case ',':
			if !inString {
				highlighted.WriteString(Green + string(char) + Reset )
			} else {
				highlighted.WriteString(string(char))
			}
			isKey = true
		case ':':
			if !inString {
				highlighted.WriteString(Red + string(char) + Reset)
			} else {
				highlighted.WriteString(string(char))
			}
			isKey = false
		case '"':
			highlighted.WriteString(Green + string(char))
			inString = !inString
			
		default:
			if isKey {
				highlighted.WriteString(Yellow + string(char) + Reset)
			} else {
				highlighted.WriteString(string(char))
			}
		}
	}

	return highlighted.String()
}
