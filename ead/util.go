package ead

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/nyulibraries/dlts-finding-aids-ead-go-packages/ead/validate"
)

func getConvertedTextWithTags(text string) ([]byte, error) {
	return _getConvertedTextWithTags(text, true)
}

func getConvertedTextWithTagsNoLBConversion(text string) ([]byte, error) {
	return _getConvertedTextWithTags(text, false)
}

func _getConvertedTextWithTags(text string, convertLBTags bool) ([]byte, error) {
	decoder := xml.NewDecoder(strings.NewReader(text))

	var result string
	needClosingTag := true
	for {
		token, err := decoder.Token()

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		switch token := token.(type) {
		case xml.StartElement:
			switch token.Name.Local {
			default:
				result += _getConvertedTextWithTagsDefault(token.Name.Local)
			case "emph":
				{
					var render string
					for i := range token.Attr {
						if token.Attr[i].Name.Local == "render" {
							render = token.Attr[i].Value
							break
						}
					}
					result += fmt.Sprintf("<span class=\"%s\">", "ead-emph ead-emph-"+render)
				}
			case "lb":
				{
					if convertLBTags {
						result += "<br>"
						needClosingTag = false
					} else {
						result += _getConvertedTextWithTagsDefault(token.Name.Local)
					}
				}
			}

		case xml.EndElement:
			if needClosingTag {
				result += "</span>"
			} else {
				// Reset
				needClosingTag = true
			}
		case xml.CharData:
			result += strings.ReplaceAll(string(token), "\n", " ")
		}
	}

	return []byte(cleanupWhitespace(result)), nil
}

func _getConvertedTextWithTagsDefault(tagName string) string {
	return fmt.Sprintf("<span class=\"ead-%s\">", tagName)
}

func getRelatorAuthoritativeLabel(relatorID string) (string, error) {
	if authoritativeLabel, ok := validate.RelatorAuthoritativeLabelMap[relatorID]; ok {
		return authoritativeLabel, nil
	} else {
		return "", errors.New(fmt.Sprintf("Unknown relator code \"%s\"", relatorID))
	}
}

func regexpReplaceAllLiteralStringInAccessTermSourceSlice(accessTermWithRoleSlice []AccessTermWithRole, re *regexp.Regexp, replacementString string) {
	accessTermWithRoleSliceWithSubfieldDelimitersConverted := accessTermWithRoleSlice[:0]
	for _, accessTermWithRole := range accessTermWithRoleSlice {
		accessTermWithRole.Value = re.ReplaceAllLiteralString(accessTermWithRole.Value, replacementString)
		accessTermWithRoleSliceWithSubfieldDelimitersConverted = append(
			accessTermWithRoleSliceWithSubfieldDelimitersConverted,
			accessTermWithRole,
		)
	}
}

func regexpReplaceAllLiteralStringInTextSlice(textSlice []string, re *regexp.Regexp, replacementString string) {
	accessTermSWithRoleliceWithSubfieldDelimitersConverted := textSlice[:0]
	for _, text := range textSlice {
		accessTermSWithRoleliceWithSubfieldDelimitersConverted = append(
			accessTermSWithRoleliceWithSubfieldDelimitersConverted,
			re.ReplaceAllLiteralString(text, replacementString),
		)
	}
}

// RunInfo stores data related to the parsing/JSON generation process
type RunInfo struct {
	PkgVersion string    `json:"libversion"`
	TimeStamp  time.Time `json:"timestamp"`
	SourceFile string    `json:"sourcefile"`
}

func (r *RunInfo) SetRunInfo(version string, t time.Time, sourceFile string) {
	r.PkgVersion = version
	r.TimeStamp = t
	r.SourceFile = sourceFile
}

// Donors is slice containing Donor names
type Donors []FilteredString

// PubInfo stores data used by the publication system
type PubInfo struct {
	ThemeID string `json:"themeid"`
}

func (p *PubInfo) SetPubInfo(themeid string) {
	p.ThemeID = themeid
}

// FilteredString provides a centralized string cleanup mechanism
type FilteredString string

func (s FilteredString) String() string {
	return cleanupWhitespace(string(s))
}

func (s FilteredString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func cleanupWhitespace(s string) string {
	// find occurrences of one or more consecutive \r, \n, \t, " "
	re := regexp.MustCompile(`\r+|\n+|\t+|( )+`)
	// replace occurences with a single space
	result := re.ReplaceAllString(s, " ")
	// run again in case previous operation resulted in
	// consecutive spaces
	result = re.ReplaceAllString(result, " ")
	// clean off any leading/trailing whitespace
	result = strings.TrimSpace(result)
	return result
}

type FilteredLabelString FilteredString

func (s FilteredLabelString) MarshalJSON() ([]byte, error) {

	return json.Marshal(cleanupWhitespace(removeBracketedText(string(s))))
}

func removeBracketedText(s string) string {
	// find bracketed text
	re := regexp.MustCompile(`\[.+\]`)
	// remove occurences
	result := re.ReplaceAllString(s, "")
	return result
}

func flattenTitleProper(titleProper []*TitleProper) ([]byte, error) {

	var titleToFlatten *TitleProper

	// capture first titleProper not of type "filing"
	for _, t := range titleProper {
		titleToFlatten = t
		// ignore "filing" title
		if titleToFlatten.Type != "filing" {
			break
		}
	}

	// we only found the "filing" title. This is a problem!
	if titleToFlatten.Type == "filing" {
		return nil, fmt.Errorf("Unable to find correct title\n")
	}

	return getConvertedTextWithTagsNoLBConversion(titleToFlatten.Value)
}
