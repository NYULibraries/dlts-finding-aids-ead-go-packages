package ead

import (
	"encoding/json"
	"regexp"
	"strings"
)

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

	flattenedTitleProper, err := flattenTitleProper(titleStmt.TitleProper)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		*TitleStmtAlias
		FlattenedTitleProper FilteredString `json:"titleproper,omitempty"`
	}{
		TitleStmtAlias:       (*TitleStmtAlias)(titleStmt),
		FlattenedTitleProper: FilteredString(flattenedTitleProper),
	})
}

// set blank DAO Role attributes to "external-link"
func (dao *DAO) MarshalJSON() ([]byte, error) {
	// if DAO Role is empty, set it to external link
	type DAOAlias DAO
	if len(strings.TrimSpace(string(dao.Role))) == 0 {
		dao.Role = "external-link"
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
