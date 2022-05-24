package validate

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/nyulibraries/dlts-finding-aids-ead-go-packages/ead"
)

const ValidEADIDRegexpString = "^[a-z0-9]+(?:_[a-z0-9]+){1,7}$"

var ValidEADIDRegexp *regexp.Regexp

var ValidRepositoryNames = []string{
	"Akkasah: Center for Photography (NYU Abu Dhabi)",
	"Center for Brooklyn History",
	"Fales Library and Special Collections",
	"NYU Abu Dhabi, Archives and Special Collections",
	"New York University Archives",
	"New-York Historical Society",
	"Poly Archives at Bern Dibner Library of Science and Technology",
	"Tamiment Library and Robert F. Wagner Labor Archives",
	"Villa La Pietra",
}

func init() {
	var err error
	ValidEADIDRegexp, err = regexp.Compile(ValidEADIDRegexpString)
	if err != nil {
		// TODO: Figure out what to do here...in theory this can't ever fail because
		// we're compiling a constant.  If it does fail, might want to avoid panic()
		// calls because might be in use in the FAM API server, which in theory
		// should be able to trap panic calls, but what if it (or whatever client)
		// doesn't?
	}
}

func ValidateEAD(bytes []byte) ([]string, error) {
	var validationErrors = []string{}

	validationErrors = append(validationErrors, validateXML(bytes)...)

	var ead ead.EAD
	err := xml.Unmarshal(bytes, &ead)
	if err != nil {
		return validationErrors, err
	}

	validationErrors = append(validationErrors, validateRequiredEADElements(ead)...)
	validationErrors = append(validationErrors, validateRepository(ead)...)
	validationErrors = append(validationErrors, validateEADID(ead)...)
	validationErrors = append(validationErrors, validateNoUnpublishedMaterial(ead)...)
	validationErrors = append(validationErrors, validateRoleAttributes(ead)...)

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

func makeInvalidEADIDErrorMessage(eadid string, invalidCharacters []rune) string {
	return fmt.Sprintf(`Invalid <eadid>

<eadid> value "%s" does not conform to the Finding Aids specification.
There must be between 2 to 8 character groups joined by underscores.
The following characters are not allowed in character groups: %s
`, eadid, string(invalidCharacters))
}

func makeInvalidRepositoryErrorMessage(repositoryName string) string {
	return fmt.Sprintf(`Invalid <repository>

<repository> contains unknown repository name "%s".
The repository name must match a value from this list:

%s
`, repositoryName, strings.Join(ValidRepositoryNames, "\n"))
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

func validateEADID(ead ead.EAD) []string {
	var validationErrors = []string{}

	var EADID = ead.EADHeader.EADID.Value

	match := ValidEADIDRegexp.Match([]byte(EADID))
	if !match {
		var invalidCharacters = []rune{}
		charMap := make(map[rune]uint, len(EADID))
		for _, r := range EADID {
			charMap[r]++
		}
		for char, _ := range charMap {
			if !(unicode.IsLower(char) || unicode.IsDigit(char)) {
				invalidCharacters = append(invalidCharacters, char)
			}
		}
		validationErrors = append(validationErrors, makeInvalidEADIDErrorMessage(EADID, invalidCharacters))
	}

	return validationErrors
}

func validateNoUnpublishedMaterial(ead ead.EAD) []string {
	var validationErrors = []string{}

	return validationErrors
}

func validateRepository(ead ead.EAD) []string {
	var validationErrors = []string{}

	var repositoryName = ead.ArchDesc.DID.Repository.CorpName[0].Value

	for _, validRepository := range ValidRepositoryNames {
		if repositoryName == validRepository {
			return []string{}
		}
	}

	validationErrors = append(validationErrors, makeInvalidRepositoryErrorMessage(repositoryName))

	return validationErrors
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

func validateRoleAttributes(ead ead.EAD) []string {
	var validationErrors = []string{}

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
