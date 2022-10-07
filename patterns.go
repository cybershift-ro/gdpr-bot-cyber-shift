package main

import (
	"regexp"
	"strings"
)

var keywords = []string{
	`(?i)amend(ă|a)`,
	`(?i)amend(ă|a)\s+pentru`,
	`(?i)sanc(ţ|t)iune`,
	`(?i)sanc(ţ|t)iune\s+pentru`,
	`(?i)(î|i)nc(ă|a)lcarea`,
	`(?i)(î|i)nc(ă|a)lcarea\s+RGPD`,
}

var companyRegex = []string{
	`societatea\s+[a-zA-Z]+\s+[a-zA-Z]+\s+[a-zA-Z]+\sa\sfost`,
	`operatorul\s+[a-zA-Z]+\s+[a-zA-Z]+\s+[a-zA-Z]+\sa\sfost`,
	`operatorul\s+[a-zA-Z]+\s+[a-zA-Z]+\s+(S\.A\.|SA|S\.A)\s+și\sa`,
}

var sanctionRegex = []string{
	`(?i)cuantum\s+(.)+\d[\d,.]*\s+(ron|lei|euro|de euro)(\s+|,)`,
	`(?i)cuantum\s+(.)+\d[\d,.]*\s+(ron|lei|euro|de euro)(\s+|,|\.)`,
}

var compiledCompanyRegex = []*regexp.Regexp{}
var compiledSanctionRegex = []*regexp.Regexp{}

var compiledKeyWords = []*regexp.Regexp{}

func compileAllKeywords() {
	for _, triggers := range keywords {
		compiled := regexp.MustCompile(triggers)
		compiledKeyWords = append(compiledKeyWords, compiled)
	}

	for _, triggers := range companyRegex {
		compiled := regexp.MustCompile(triggers)
		compiledCompanyRegex = append(compiledCompanyRegex, compiled)
	}

	for _, triggers := range sanctionRegex {
		compiled := regexp.MustCompile(triggers)
		compiledSanctionRegex = append(compiledSanctionRegex, compiled)
	}
}

func isSaction(title string) bool {

	for _, keyword := range compiledKeyWords {
		if keyword.MatchString(title) {
			return true
		}
	}

	return false
}

func extractCompanyName(paragraph string) string {
	for _, pattern := range compiledCompanyRegex {

		match := pattern.FindStringSubmatch(paragraph)

		if len(match) > 0 {
			words := strings.Fields(match[0])

			// Exclude first word and last 2 words
			words = words[1 : len(words)-2]

			// Return a single string that should contain the company name
			return strings.Join(words[:], " ")
		}
	}

	return ""
}

func extractSanctionSum(paragraph string) string {
	for _, pattern := range compiledSanctionRegex {

		match := pattern.FindStringSubmatch(paragraph)

		if len(match) > 0 {
			words := strings.Fields(match[0])

			// Exclude first 2 words and last word
			words = words[2 : len(words)-1]

			// Edge case '9.671,40 lei, echivalentul a 2.000'
			if len(words) > 2 {
				words = words[:2]
				words[1] = strings.TrimSuffix(words[1], ",")
			}

			// Return a single string that should contain the company name
			return strings.Join(words[:], " ")
		}
	}

	return ""
}

func leaveOnlyNumbers(s []byte) string {
	j := 0
	for _, b := range s {
		if ('0' <= b && b <= '9') || b == '.' || b == ',' {
			s[j] = b
			j++
		}
	}
	return string(s[:j])
}
