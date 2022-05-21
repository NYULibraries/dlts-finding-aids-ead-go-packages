package validate

import "fmt"

func ValidateEAD(ead []byte) []string {
	return []string{
		"RED LIGHT #1",
		"RED LIGHT #2",
		"RED LIGHT #3",
	}
}

func makeInvalidXMLErrorMessage() string {
	return "The XML in this file is not valid.  Please check it using an XML validator."
}

func makeMissingRequiredElementErrorMessage(elementName string) string {
	return fmt.Sprintf("Required element %s not found.", elementName)
}
