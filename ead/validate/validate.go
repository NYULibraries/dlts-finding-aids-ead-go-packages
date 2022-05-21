package validate

import (
	"bytes"
	"fmt"
	"strings"
)

func ValidateEAD(ead []byte) []string {
	return []string{
		"RED LIGHT #1",
		"RED LIGHT #2",
		"RED LIGHT #3",
	}
}

func makeAudienceInternalErrorMessage(elementsAudienceInternal []string) string {
	return fmt.Sprintf(`Private data detected

The EAD file contains unpublished material.  The following EAD elements have attribute audience="internal" and must be removed:

%s`, strings.Join(elementsAudienceInternal, "\n"))
}

func makeInvalidXMLErrorMessage() string {
	return "The XML in this file is not valid.  Please check it using an XML validator."
}

func makeMissingRequiredElementErrorMessage(elementName string) string {
	return fmt.Sprintf("Required element %s not found.", elementName)
}

func makeInvalidEADIDErrorMessage(eadid string, invalidCharacters []byte) string {
	return fmt.Sprintf(`Invalid <eadid>

<eadid> value "%s" does not conform to the Finding Aids specification.
There must be between 2 to 8 character groups joined by underscores.
The following characters are not allowed in character groups: %s
`, eadid, string(bytes.Join([][]byte{invalidCharacters}, []byte(" "))))
}

func makeInvalidRepositoryErrorMessage(repositoryName string) string {
	return fmt.Sprintf(`Invalid <repository>

	<repository> contains unknown repository name "%s".
		The repository name must match a value from this list:

Akkasah: Center for Photography (NYU Abu Dhabi)
	New York University Archives
	Center for Brooklyn History
	Fales Library and Special Collections
	Villa La Pietra
	New-York Historical Society
	NYU Abu Dhabi, Archives and Special Collections
	Poly Archives at Bern Dibner Library of Science and Technology
	Tamiment Library and Robert F. Wagner Labor Archives`, repositoryName)
}

func makeUnrecognizedRelatorCodesErrorMessage(unrecognizedRelatorCodes [][]string) string {
	var unrecognizedRelatorCodeSlice []string
	for _, elementAttributePair := range unrecognizedRelatorCodes {
		unrecognizedRelatorCodeSlice = append(unrecognizedRelatorCodeSlice,
			fmt.Sprintf(`%s has role="%s"`, elementAttributePair[0], elementAttributePair[1]))
	}

	return fmt.Sprintf(`Unrecognized relator codes

The EAD file contains elements with role attributes containing unrecognized relator codes:

%s`, strings.Join(unrecognizedRelatorCodeSlice, "\n"))
}
