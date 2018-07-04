package utils

import (
	"regexp"
	"strings"

	"github.com/maddevsio/go-idmatch/log"
	"github.com/maddevsio/go-idmatch/templates"
)

const ErrorMessage = "(recognition failed)"

func gender(gender string) string {
	if strings.ContainsAny(gender, "m M э Э м М 3 9 5") {
		return "Э"
	} else if strings.ContainsAny(gender, "f F а А ж Ж") {
		return "А"
	}
	return ""
}

func Sanitize(documentMap map[string]interface{}, card templates.Card) {
	regex := "[^а-яa-zА-ЯA-Z0-9№ ]+"

	for _, v := range card.Structure {
		if documentMap[v.Field] == nil {
			continue
		}
		text := documentMap[v.Field].(string)

		switch v.Type {
		case "cyrillic":
			regex = "[^а-яА-Я№ ]+"
		case "latin":
			regex = "[^a-zA-Z ]+"
		case "number":
			regex = "[^0-9]+"
		case "gender":
			text = gender(text)
			regex = "[^а-яА-Я]$"
		}

		if n := strings.Index(text, "\n"); n > 0 {
			text = text[:n]
		}
		reg, err := regexp.Compile(regex)
		if err != nil {
			log.Print(log.ErrorLevel, err.Error())
		}
		clearText := reg.ReplaceAllString(text, "")
		if len(clearText) == 0 {
			clearText = ErrorMessage
			continue
		}
		if v.Length != 0 && len(clearText) > v.Length {
			clearText = clearText[len(clearText)-v.Length:]
		}
		if v.Prefix != "" {
			clearText = v.Prefix + clearText
		}
		// else if text != clearText {
		// 	clearText += " (?)"
		// }
		documentMap[v.Field] = strings.ToUpper(clearText)
	}
}
