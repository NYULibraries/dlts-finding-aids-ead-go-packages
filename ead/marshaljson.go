package ead

import (
	"encoding/json"
//	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// var unitDateNormalRegexp = regexp.MustCompile(`(.+)/(.+)`)
// var unitDateNormalExpectedStringSubmatchCount = 3

// Note that this custom marshalling for DID will prevent PhysDesc from having a Value field
// that is all whitespace if Extent is nil, but won't prevent PhysDesc from having
// a Value field that is all whitespace if Extent is not nil.
// We need to convert Value field values like "\n    \n    \n" to empty string
// so they can be removed by omitempty struct tag.  This is done in the PhysDesc.MarshalJSON
// in marshal-generated.go.
func (did *DID) MarshalJSON() ([]byte, error) {
	type DIDWithNoEmptyPhysDesc DID

	containsNonWhitespaceRegexp := regexp.MustCompile(`\S`)
	var physDescNoEmpties []*PhysDesc
	for _, el := range did.PhysDesc {
		if el.Extent != nil || containsNonWhitespaceRegexp.MatchString(el.Value) {
			physDescNoEmpties = append(physDescNoEmpties, el)
		}
	}

	var jsonData []byte
	var err error
	if physDescNoEmpties != nil {
		jsonData, err = json.Marshal(&struct {
			PhysDesc []*PhysDesc `xml:"physdesc" json:"physdesc,omitempty"`
			*DIDWithNoEmptyPhysDesc
		}{
			PhysDesc:               physDescNoEmpties,
			DIDWithNoEmptyPhysDesc: (*DIDWithNoEmptyPhysDesc)(did),
		})
	} else {
		jsonData, err = json.Marshal(&struct {
			PhysDesc []*PhysDesc `xml:"physdesc" json:"physdesc,omitempty"`
			*DIDWithNoEmptyPhysDesc
		}{
			PhysDesc:               nil,
			DIDWithNoEmptyPhysDesc: (*DIDWithNoEmptyPhysDesc)(did),
		})
	}

	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (accessTermWithRole *AccessTermWithRole) MarshalJSON() ([]byte, error) {
	type accessTermWithRoleWithTranslatedRelatorCode AccessTermWithRole

	var (
		role string
		err  error
	)
	if accessTermWithRole.Role != "" {
		role, err = getRelatorAuthoritativeLabel(accessTermWithRole.Role)
		if err != nil {
			return nil, err
		}
	}

	result, err := getConvertedTextWithTags(accessTermWithRole.Value)
	if err != nil {
		return nil, err
	}
	accessTermWithRole.Value = string(result)

	jsonData, err := json.Marshal(&struct {
		Role string `xml:"role,attr" json:"role,omitempty"`
		*accessTermWithRoleWithTranslatedRelatorCode
	}{
		Role: role,
		accessTermWithRoleWithTranslatedRelatorCode: (*accessTermWithRoleWithTranslatedRelatorCode)(accessTermWithRole),
	})
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (titleproper *TitleProper) MarshalJSON() ([]byte, error) {
	type TitleProperWithTags TitleProper

	result, err := getConvertedTextWithTagsNoLBConversion(titleproper.Value)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(&struct {
		Value string `json:"value,omitempty"`
		*TitleProperWithTags
	}{
		Value:               string(result),
		TitleProperWithTags: (*TitleProperWithTags)(titleproper),
	})
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (titleStmt *TitleStmt) MarshalJSON() ([]byte, error) {
	type TitleStmtAlias TitleStmt

	flattenedAuthor, err := flattenCDATA(titleStmt.Author)
	if err != nil {
		return nil, err
	}

	flattenedSponsor, err := flattenCDATA(titleStmt.Sponsor)
	if err != nil {
		return nil, err
	}

	flattenedSubTitle, err := flattenCDATA(titleStmt.SubTitle)
	if err != nil {
		return nil, err
	}

	flattenedTitleProper, err := flattenTitleProper(titleStmt.TitleProper)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		*TitleStmtAlias
		FlattenedAuthor      FilteredString `json:"author,omitempty"`
		FlattenedSponsor     FilteredString `json:"sponsor,omitempty"`
		FlattenedSubTitle    FilteredString `json:"subtitle,omitempty"`
		FlattenedTitleProper FilteredString `json:"titleproper,omitempty"`
	}{
		TitleStmtAlias:       (*TitleStmtAlias)(titleStmt),
		FlattenedAuthor:      FilteredString(flattenedAuthor),
		FlattenedSponsor:     FilteredString(flattenedSponsor),
		FlattenedSubTitle:    FilteredString(flattenedSubTitle),
		FlattenedTitleProper: FilteredString(flattenedTitleProper),
	})
}

func (indexEntry *IndexEntry) MarshalJSON() ([]byte, error) {
	type IndexEntryAlias IndexEntry

	flattenedRef, err := flattenCDATA(indexEntry.Ref)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		*IndexEntryAlias
		FlattenedRef FilteredString `json:"ref,omitempty"`
	}{
		IndexEntryAlias: (*IndexEntryAlias)(indexEntry),
		FlattenedRef:    FilteredString(flattenedRef),
	})
}

// set blank DAO Role attributes to "external-link"
// if a DAO has a non-URL HREF, then set the role to "non-url"
func (dao *DAO) MarshalJSON() ([]byte, error) {
	// if DAO Role is empty, set it to external link
	type DAOAlias DAO
	if len(strings.TrimSpace(string(dao.Role))) == 0 {
		dao.Role = "external-link"
	}

	_, err := url.ParseRequestURI(string(dao.Href))
	if err != nil {
		dao.Role = "non-url"
	}

	return json.Marshal(&struct {
		*DAOAlias
	}{
		DAOAlias: (*DAOAlias)(dao),
	})
}

func (extent *Extent) MarshalJSON() ([]byte, error) {
	type ExtentWithTags Extent

	// this code tempararily adds the unit string, if present
	// to the extent.Value for Marshaling
	valueSave := extent.Value
	if extent.Unit != "" {
		extent.Value = extent.Value + " " + extent.Unit.String()
	}

	result, err := getConvertedTextWithTags(extent.Value)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(&struct {
		Value string `json:"value,omitempty"`
		*ExtentWithTags
	}{
		Value:          string(result),
		ExtentWithTags: (*ExtentWithTags)(extent),
	})
	if err != nil {
		return nil, err
	}

	// restore the saved extent.Value
	extent.Value = valueSave

	return jsonData, nil
}

func (fnwh *FormattedNoteWithHead) MarshalJSON() ([]byte, error) {
	type FormattedNoteWithHeadAlias FormattedNoteWithHead

	// if there are no children then create a child from innerxml...
	if len(fnwh.Children) == 0 {
		// Children array is empty, therefore flatten innerXML
		flattenedValue, err := getConvertedTextWithTags(fnwh.Value)
		if err != nil {
			return nil, err
		}

		// create and add child element
		// the nesting of "value": is for consistency with the marshaled
		// JSON of regular stream-parsed children
		child := EADChild{}
		child.Name = "div"
		child.Value = &struct {
			Value FilteredString `json:"value,omitempty"`
		}{
			Value: FilteredString(flattenedValue),
		}
		fnwh.Children = append(fnwh.Children, &child)
	}

	return json.Marshal(&struct {
		*FormattedNoteWithHeadAlias
	}{
		FormattedNoteWithHeadAlias: (*FormattedNoteWithHeadAlias)(fnwh),
	})
}

// func (unitdate *UnitDate) MarshalJSON() ([]byte, error) {
// 	type UnitDateWithTags UnitDate

// 	result, err := getConvertedTextWithTags(unitdate.Value)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// clean up any blank space
// 	result = []byte(strings.TrimSpace(string(result)))

// 	// if result is empty...
// 	if len(result) == 0 {
// 		// check if we have a date value in Normal
// 		if len(strings.TrimSpace(string(unitdate.Normal))) == 0 {
// 			// nothing here, so omit this unitdate by setting the
// 			// unitdate variable to an empty struct, which will
// 			// be omitted during marshaling
// 			unitdate = &UnitDate{}
// 		} else {
// 			// ok, we found something. Let's convert the value...
// 			matches := unitDateNormalRegexp.FindStringSubmatch(string(unitdate.Normal))

// 			if len(matches) < unitDateNormalExpectedStringSubmatchCount {
// 				return nil, fmt.Errorf("problem parsing UnitDate.Normal")
// 			}
// 			// extract the values and configure result accordingly
// 			dateA := matches[1]
// 			dateB := matches[2]

// 			if dateA == dateB {
// 				result = []byte(dateA)
// 			} else {
// 				result = []byte(dateA + "-" + dateB)
// 			}
// 		}
// 	}
// 	jsonData, err := json.Marshal(&struct {
// 		Value string `json:"value,omitempty"`
// 		*UnitDateWithTags
// 	}{
// 		Value:            string(result),
// 		UnitDateWithTags: (*UnitDateWithTags)(unitdate),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return jsonData, nil
// }
