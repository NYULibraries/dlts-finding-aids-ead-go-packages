package validate

import "fmt"

func ValidateEAD(ead []byte) []string {
	return []string{
		"RED LIGHT #1",
		"RED LIGHT #2",
	}
}

func makeMissingRequiredElementErrorMessage(elementName string) string {
	return fmt.Sprintf("Required element %s not found.", elementName)
}
