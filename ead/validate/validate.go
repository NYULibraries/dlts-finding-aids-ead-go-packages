package validate

import (
	"bytes"
	"embed"
	_ "embed"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/antchfx/xmlquery"
	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/xsd"

	"github.com/nyulibraries/dlts-finding-aids-ead-go-packages/ead"
)

//go:embed schema
var schemas embed.FS

const ValidEADIDRegexpString = "^[a-z0-9]+(?:_[a-z0-9]+){1,7}$"

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

func ValidateEAD(data []byte) ([]string, error) {
	var validationErrors = []string{}

	// Performing a two stage validation here:
	// 1.) validateXML is just making sure that we have a valid XML document
	// 2.) validateEADAgainstSchema is confirming that the document conforms
	//     to the EAD schema
	//
	// This approach avoids unwanted libxml2 output when the data is not XML
	validationErrors = append(validationErrors, validateXML(data)...)
	// If the data is not valid XML there is no point doing any more checks.
	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	validationErrors = append(validationErrors, validateEADAgainstSchema(data)...)
	// If the data is not valid XML there is no point doing any more checks.
	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	var ead ead.EAD
	err := xml.Unmarshal(data, &ead)
	if err != nil {
		return validationErrors, err
	}

	validateEADIDValidationErrors, err := validateEADID(ead)
	if err != nil {
		return validationErrors, err
	}
	validationErrors = append(validationErrors, validateEADIDValidationErrors...)

	validationErrors = append(validationErrors, validateRepository(ead)...)

	validateNoUnpublishedMaterialValidationErrors, err := validateNoUnpublishedMaterial(data)
	if err != nil {
		return validationErrors, err
	}

	validationErrors = append(validationErrors, validateNoUnpublishedMaterialValidationErrors...)

	validateRoleAttributesValidationErrors, err := validateRoleAttributes(data)
	if err != nil {
		return validationErrors, err
	}
	validationErrors = append(validationErrors, validateRoleAttributesValidationErrors...)

	validationErrors = append(validationErrors, validateHREFs(ead)...)

	return validationErrors, err
}

func makeAudienceInternalErrorMessage(elementsAudienceInternal []string) string {
	return fmt.Sprintf(`Private data detected

The EAD file contains unpublished material.  The following EAD elements have attribute audience="internal" and must be removed:

%s`, strings.Join(elementsAudienceInternal, "\n"))
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

func makeInvalidXMLErrorMessage() string {
	return "The XML in this file is not valid.  Please check it using an XML validator."
}

func makeMissingRequiredElementErrorMessage(elementName string) string {
	return fmt.Sprintf("Required element %s not found.", elementName)
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

func validateEADID(ead ead.EAD) ([]string, error) {
	var validationErrors = []string{}

	var EADID = ead.EADHeader.EADID.Value

	// Even if the file contains only "<ead></ead>", ead.EADHeader.EADID.Value will
	// be not empty.  Test for empty string.
	if EADID != "" {
		trimmedEADID := strings.TrimSpace(EADID)
		match, err := regexp.Match(ValidEADIDRegexpString, []byte(trimmedEADID))
		if err != nil {
			return validationErrors, err
		}

		if !match {
			var invalidCharacters = []rune{}
			charMap := make(map[rune]uint, len(EADID))
			for _, r := range EADID {
				charMap[r]++
			}
			for char, _ := range charMap {
				if !(unicode.IsLower(char) || unicode.IsDigit(char) || char == '_') {
					invalidCharacters = append(invalidCharacters, char)
				}
			}
			validationErrors = append(validationErrors, makeInvalidEADIDErrorMessage(EADID, invalidCharacters))
		}
	} else {
		validationErrors = append(validationErrors,
			makeMissingRequiredElementErrorMessage("<eadid>"))
	}

	return validationErrors, nil
}

func validateNoUnpublishedMaterial(data []byte) ([]string, error) {
	var validationErrors = []string{}

	decoder := xml.NewDecoder(bytes.NewReader(data))

	audienceInternalElements := []string{}
	for {
		token, err := decoder.Token()
		if token == nil || err == io.EOF {
			break
		} else if err != nil {
			return []string{}, err
		}

		switch tokenType := token.(type) {
		case xml.StartElement:
			var elementName = tokenType.Name.Local
			for _, attribute := range tokenType.Attr {
				attributeName := attribute.Name.Local
				attributeValue := attribute.Value

				if attributeName == "audience" && attributeValue == "internal" {
					audienceInternalElements = append(audienceInternalElements,
						fmt.Sprintf("<%s>", elementName))
				}
			}
		default:
		}
	}

	if len(audienceInternalElements) > 0 {
		validationErrors = append(validationErrors, makeAudienceInternalErrorMessage(audienceInternalElements))
	}

	return validationErrors, nil
}

func validateRepository(ead ead.EAD) []string {
	var validationErrors = []string{}

	if ead.ArchDesc != nil {
		// If ead.ArchDesc exists, DID will be non-nil, so can move on to testing Repository.
		if ead.ArchDesc.DID.Repository != nil {
			if ead.ArchDesc.DID.Repository.CorpName == nil {
				validationErrors = append(validationErrors,
					makeMissingRequiredElementErrorMessage("<archdesc>/<did>/<repository>/<corpname>"))
			} else {
				var repositoryName = ead.ArchDesc.DID.Repository.CorpName[0].Value

				for _, validRepository := range ValidRepositoryNames {
					if repositoryName == validRepository {
						return []string{}
					}
				}

				validationErrors = append(validationErrors, makeInvalidRepositoryErrorMessage(repositoryName))
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

// See https://jira.nyu.edu/browse/FADESIGN-171.
func validateRoleAttributes(data []byte) ([]string, error) {
	var validationErrors = []string{}

	// Parent elements
	const controlAccessElementName = "controlaccess"
	const originationElementName = "origination"
	const repositoryElementName = "repository"

	// Child elements with `role` attributes that we need to test (within the
	// context of the above parent elements).
	const corpnameElementName = "corpname"
	const famnameElementName = "famname"
	const persnameElementName = "persname"

	// Note that we are testing role attributes only for very specific occurrences of
	// the child elements, hence the need for this 2-dimensional slice of slices.
	var elementsToTest = [][]string{
		{controlAccessElementName, corpnameElementName},
		{controlAccessElementName, famnameElementName},
		{controlAccessElementName, persnameElementName},

		{originationElementName, corpnameElementName},
		{originationElementName, famnameElementName},
		{originationElementName, persnameElementName},

		{repositoryElementName, corpnameElementName},
	}

	doc, err := xmlquery.Parse(strings.NewReader(string(data)))
	if err != nil {
		return validationErrors, err
	}

	// Slice of string slices, where an inner slice element is of the form:
	// {"<repository><corpname>NYU Archives</corpname></repository>", "grt"}.
	var unrecognizedRelatorCodes [][]string
	for _, elementToTest := range elementsToTest {
		var parentElementName = elementToTest[0]
		var childElementName = elementToTest[1]

		roleAttributes, err := xmlquery.QueryAll(doc, fmt.Sprintf("//%s/%s/@role", parentElementName, childElementName))
		if err != nil {
			return validationErrors, err
		}

		for _, roleAttribute := range roleAttributes {
			relatorCode := roleAttribute.FirstChild.Data
			_, ok := ead.RelatorAuthoritativeLabelMap[relatorCode]
			if !ok {
				// Example: "<repository><corpname>NYU Archives</corpname></repository>"
				elementDescription := fmt.Sprintf("<%s><%s>%s</%s></%s>",
					parentElementName,
					childElementName,
					roleAttribute.Parent.FirstChild.Data,
					childElementName,
					parentElementName,
				)

				unrecognizedRelatorCodes = append(unrecognizedRelatorCodes, []string{
					elementDescription,
					relatorCode,
				})
			}
		}
	}

	if len(unrecognizedRelatorCodes) > 0 {
		validationErrors = append(validationErrors, makeUnrecognizedRelatorCodesErrorMessage(unrecognizedRelatorCodes))
	}

	return validationErrors, nil
}

// The following comment and function validateXML() are from David Arjanik:
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

// I did a quick search for some 3rd party libraries for validating against a schema,
// which would allow for validation against https://www.loc.gov/ead/eadschema.html, but
// I have some reservations about using them -- see https://jira.nyu.edu/browse/FADESIGN-491.
func validateXML(data []byte) []string {
	var validationErrors = []string{}

	// Not perfect, but maybe good enough for now.
	if xml.Unmarshal(data, new(interface{})) != nil {
		validationErrors = append(validationErrors, makeInvalidXMLErrorMessage())
	}

	return validationErrors
}

// This function is largely borrowed from Don Mennerich's go-aspace package
// https://github.com/nyudlts/go-aspace
func validateEADAgainstSchema(data []byte) []string {
	var validationErrors = []string{}

	// initialize with a default error message
	validationErrors = append(validationErrors, makeInvalidXMLErrorMessage())

	schema, err := schemas.ReadFile("schema/ead-2002-20210412-dlts.xsd")
	if err != nil {
		return append(validationErrors, err.Error())
	}

	eadxsd, err := xsd.Parse(schema)
	if err != nil {
		return append(validationErrors, err.Error())
	}
	defer eadxsd.Free()

	p := parser.New()
	doc, err := p.Parse(data)
	if err != nil {
		validationErrors = append(validationErrors, "Unable to parse XML file")
		return append(validationErrors, err.Error())
	}
	defer doc.Free()

	err = eadxsd.Validate(doc)
	if err != nil {
		// capture the high-level error message
		validationErrors = append(validationErrors, err.Error())

		// capture the detailed error info:
		// the Validate function returns an xsd.SchemaValidationError that
		// has an underlying Errors() method with more detailed errors
		for _, e := range err.(xsd.SchemaValidationError).Errors() {
			validationErrors = append(validationErrors, e.Error())
		}
		return validationErrors
	}

	// all ok, return empty slice
	return []string{}
}

func validateHREFs(ead ead.EAD) []string {
	var validationErrors = []string{}

	// REQUIRED!
	ead.InitDAOCounts()

	for _, dao := range ead.DAOInfo.AllDAOs {
		// https://golang.cafe/blog/how-to-validate-url-in-go.html
		_, err := url.ParseRequestURI(string(dao.Href))
		if err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("Invalid HREF detected: '%s', Title: '%s'", []byte(dao.Href), dao.Title))
		}
	}

	return validationErrors
}
