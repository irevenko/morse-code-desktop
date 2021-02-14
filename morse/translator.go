package morse

import "strings"

func TranslateToMorse(text string) string {
	var morseSequence string

	text = strings.ToUpper(text)
	text = strings.Trim(text, " ")
	textSlice := strings.Split(text, "")

	for _, char := range textSlice {
		morseSequence += Alphabet[char] + " "
	}

	return morseSequence
}

func TranslateFromMorse(morseSequence string) string  {
	var translation string

	morseSequence = strings.Trim(morseSequence, " ")
	morseSlice := strings.Split(morseSequence, " ")

	for _, char := range morseSlice {
		translation += ReverseAlphabet[char]
	}

	return translation
}
