package validate

import (
	"bytes"
	"embed"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/xsd"
	"golang.org/x/text/message"

	"github.com/nyulibraries/dlts-finding-aids-ead-go-packages/ead"
)

//go:embed schema
var schemas embed.FS

const ValidEADIDRegexpString = "^[a-z0-9]+(?:_[a-z0-9]+){1,}$"
const MAXIMUM_EADID_LENGTH = 251
const MAXIMUM_FILE_SIZE = 100_000_000 // 100 MB
const ARCHDESC_REQUIRED_LEVEL = "collection"

var ValidRepositoryNames = []string{
	"Akkasah: Photography Archive (NYU Abu Dhabi)",
	"Center for Brooklyn History",
	"Fales Library and Special Collections",
	"NYU Abu Dhabi, Archives and Special Collections",
	"New York University Archives",
	"New-York Historical Society",
	"Poly Archives at the Bern Dibner Library of Science and Technology, NYU Libraries",
	"Tamiment Library and Robert F. Wagner Labor Archives",
	"Villa La Pietra",
}

// this function is required to perform file-level checks,
// like maximum file size
func ValidateEADFromFilePath(filepath string) ([]string, error) {
	var validationErrors = []string{}

	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return validationErrors, err
	}

	if fileInfo.Size() > MAXIMUM_FILE_SIZE {
		return append(validationErrors, makeFileTooBigErrorMessage(filepath, fileInfo.Size())), nil
	}

	EADXML, err := os.ReadFile(filepath)
	if err != nil {
		return validationErrors, err
	}

	return ValidateEAD(EADXML)
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
	validationErrors = append(validationErrors, validateArchDescLevel(ead)...)

	validateNoUnpublishedMaterialValidationErrors, err := validateNoUnpublishedMaterial(data)
	if err != nil {
		return validationErrors, err
	}

	validationErrors = append(validationErrors, validateNoUnpublishedMaterialValidationErrors...)

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
There must be a minimum of 2 character groups joined by an underscore.
There is no maximum number of character groups, however, the <eadid>
value must have at most %d characters.
The following characters found in the eadid value are not allowed in
character groups: "%s"
`, eadid, MAXIMUM_EADID_LENGTH, string(invalidCharacters))
}

func makeEADIDTooLongErrorMessage(eadid string) string {
	return fmt.Sprintf(`<eadid> length too long

	The <eadid> value in this EAD "%s" is %d characters.
	This exceeds the maximum allowed length of %d characters.`,
		eadid, len(eadid), MAXIMUM_EADID_LENGTH)
}

func makeFileTooBigErrorMessage(filepath string, size int64) string {
	// https://pkg.go.dev/golang.org/x/text/message
	p := message.NewPrinter(message.MatchLanguage("en"))

	return p.Sprintf(`ead file too big

	The size of the EAD file "%s"
	is %d bytes. The maximum allowed file size 
	is %d bytes.`,
		filepath, size, MAXIMUM_FILE_SIZE)
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

func makeInvalidArchDescLevelErrorMessage(level string) string {
	return fmt.Sprintf(`Invalid <archdesc> level

	The archdesc level attribute must be set to "%s".
	This EAD's archdesc level attribute is set to "%s"`,
		ARCHDESC_REQUIRED_LEVEL, level)
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
			for char := range charMap {
				if !(unicode.IsLower(char) || unicode.IsDigit(char) || char == '_') {
					invalidCharacters = append(invalidCharacters, char)
				}
			}
			validationErrors = append(validationErrors, makeInvalidEADIDErrorMessage(EADID, invalidCharacters))
		}

		if len(trimmedEADID) > MAXIMUM_EADID_LENGTH {
			validationErrors = append(validationErrors, makeEADIDTooLongErrorMessage(trimmedEADID))
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

func validateArchDescLevel(ead ead.EAD) []string {
	var validationErrors = []string{}

	if ead.ArchDesc != nil {
		level := string(ead.ArchDesc.Level)
		if level != ARCHDESC_REQUIRED_LEVEL {
			return append(validationErrors, makeInvalidArchDescLevelErrorMessage(level))
		}
	}
	return validationErrors
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
