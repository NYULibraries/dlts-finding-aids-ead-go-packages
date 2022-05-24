package validate

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/nyulibraries/dlts-finding-aids-ead-go-packages/ead"
)

func ValidateEAD(bytes []byte) ([]string, error) {
	var validationErrors = []string{}

	validationErrors = append(validationErrors, validateXML(bytes)...)

	var ead ead.EAD
	err := xml.Unmarshal(bytes, &ead)
	if err != nil {
		return validationErrors, err
	}

	validationErrors = append(validationErrors, validateRequiredEADElements(ead)...)

	return validationErrors, err
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

func validateRequiredEADElements(ead ead.EAD) []string {
	var validationErrors = []string{}

	// Even if the file contains only "<ead></ead>", ead.EADHeader.EADID.Value will
	// be non-nil.  Test for empty string.
	if ead.EADHeader.EADID.Value == "" {
		validationErrors = append(validationErrors,
			makeMissingRequiredElementErrorMessage("<eadid>"))
	}

	if ead.ArchDesc != nil {
		// If ead.ArchDesc exists, DID will be non-nil, so can move on to testing Repository.
		if ead.ArchDesc.DID.Repository != nil {
			if ead.ArchDesc.DID.Repository.CorpName == nil {
				validationErrors = append(validationErrors,
					makeMissingRequiredElementErrorMessage("<archdesc>/<did>/<repository>/<corpname>"))
			}
		} else {
			validationErrors = append(validationErrors,
				makeMissingRequiredElementErrorMessage("<archdesc>/<did>/<repository>"))
		}
	} else {
		validationErrors = append(validationErrors,
			makeMissingRequiredElementErrorMessage("<archdesc>"))
	}

	return validationErrors
}

// This is not so straightforward: https://stackoverflow.com/questions/53476012/how-to-validate-a-xml:
//
// Note that the answer from GodsBoss given in
//
// func IsValid(input string) bool {
//    decoder := xml.NewDecoder(strings.NewReader(input))
//    for {
//        err := decoder.Decode(new(interface{}))
//        if err != nil {
//            return err == io.EOF
//        }
//    }
// }
//
// ...doesn't work for our fixture file containing simply:
//
// This is not XML!
//
// ...because the first err returned is io.EOF, perhaps because no open tag was ever encountered?

// Note that the xml.Unmarshal solution we use below, from the same StackOverflow page,
// will not detect invalid XML that occurs after the start element has closed.
// e.g. <something>This is not XML!</something><<<
// ...is not well-formed, but Unmarshal never deals with the "<<<" after
// element <something>.

// There are 3rd party libraries for validating against a schema, but they
// require CGO, which we'd like to avoid for now.
// * https://github.com/krolaw/xsd
// * https://github.com/lestrrat-go/libxml2
// * https://github.com/terminalstatic/go-xsd-validate/blob/master/libxml2.go
func validateXML(bytes []byte) []string {
	var validationErrors = []string{}

	// Not perfect, but maybe good enough for now.
	if xml.Unmarshal(bytes, new(interface{})) != nil {
		validationErrors = append(validationErrors, makeInvalidXMLErrorMessage())
	}

	return validationErrors
}
