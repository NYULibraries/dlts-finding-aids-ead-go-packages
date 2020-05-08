package ead

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// Based on: "Data model for parsing EAD <archdesc> elements": https://jira.nyu.edu/jira/browse/FADESIGN-29
// ...using code generated by zek against https://github.com/nyudlts/samplesFromWeatherly/tree/216c4aab0eae53c1c06a1433b357b1d7df933bd6/FADESIGN-3
// as a starting point for non-<archdesc> elements.

const (
	Version = "0.1.0"
)

type EAD struct {
	EADHeader []*EADHeader `xml:"eadheader" json:"eadheader,omitempty"`
	ArchDesc  []*ArchDesc  `xml:"archdesc" json:"archdesc,omitempty"`
}

// https://jira.nyu.edu/jira/browse/FADESIGN-29 additions
type Abstract struct {
	ID    string `xml:"id,attr" json:"id,attr,omitempty"`
	Label string `xml:"label,attr" json:"label,attr,omitempty"`
	Value string `xml:",innerxml" json:"value,chardata,omitempty"`
}

type AddressLine struct {
	Value  string    `xml:",chardata" json:"value,chardata,omitempty"`
	ExtPtr []*ExtPtr `xml:"extptr" json:"extptr,omitempty"`
}

type ArchDesc struct {
	Level              string                   `xml:"level,attr" json:"level,attr,omitempty"`
	AcqInfo            []*FormattedNoteWithHead `xml:"acqinfo" json:"acqinfo,omitempty"`
	DID                []*DID                   `xml:"did" json:"did,omitempty"`
	DSC                []*DSC                   `xml:"dsc" json:"dsc,omitempty"`
	ScopeContent       []*FormattedNoteWithHead `xml:"scopecontent" json:"scopecontent,omitempty"`
	BiogHist           []*FormattedNoteWithHead `xml:"bioghist" json:"bioghist,omitempty"`
	AccessRestrict     []*FormattedNoteWithHead `xml:"accessrestrict" json:"accessrestrict,omitempty"`
	UserRestrict       []*FormattedNoteWithHead `xml:"userestrict" json:"userestrict,omitempty"`
	PreferCite         []*FormattedNoteWithHead `xml:"prefercite" json:"prefercite,omitempty"`
	ProcessInfo        []*FormattedNoteWithHead `xml:"processinfo" json:"processinfo,omitempty"`
	Arrangement        []*FormattedNoteWithHead `xml:"arrangement" json:"arrangement,omitempty"`
	ControlAccess      []*ControlAccess         `xml:"controlaccess" json:"controlaccess,omitempty"`
	CustodHist         []*FormattedNoteWithHead `xml:"custodhist" json:"custodhist,omitempty"`
	PhysTech           []*FormattedNoteWithHead `xml:"phystech" json:"phystech,omitempty"`
	Appraisal          []*FormattedNoteWithHead `xml:"appraisal" json:"appraisal,omitempty"`
	SeparatedMaterial  []*FormattedNoteWithHead `xml:"separatedmaterial" json:"separatedmaterial,omitempty"`
	RelatedMaterial    []*FormattedNoteWithHead `xml:"relatedmaterial" json:"relatedmaterial,omitempty"`
	Accruals           []*FormattedNoteWithHead `xml:"accruals" json:"accruals,omitempty"`
	AltFormatAvailable []*FormattedNoteWithHead `xml:"altformatavailable" json:"altformatavailable,omitempty"`
	Odd                []*FormattedNoteWithHead `xml:"odd" json:"odd,omitempty"`
	Bibliography       []*FormattedNoteWithHead `xml:"bibliography" json:"bibliography,omitempty"`
}

type C struct {
	ID         string `xml:"id,attr" json:"id,attr,omitempty"`
	Level      string `xml:"level,attr" json:"level,attr,omitempty"`
	OtherLevel string `xml:"otherlevel,attr" json:"otherlevel,attr,omitempty"`

	AccessRestrict  []*FormattedNoteWithHead `xml:"accessrestrict,omitempty" json:"accessrestrict,omitempty"`
	Accruals        []*FormattedNoteWithHead `xml:"accruals,omitempty" json:"accruals,omitempty"`
	Appraisal       []*FormattedNoteWithHead `xml:"appraisal,omitempty" json:"appraisal,omitempty"`
	Arrangement     []*FormattedNoteWithHead `xml:"arrangement,omitempty" json:"arrangement,omitempty"`
	BiogHist        []*FormattedNoteWithHead `xml:"bioghist,omitempty" json:"bioghist,omitempty"`
	C               []*C                     `xml:"c,omitempty" json:"c,omitempty"`
	ControlAccess   []*ControlAccess         `xml:"controlaccess" json:"controlaccess,omitempty"`
	DID             []*DID                   `xml:"did,omitempty" json:"did,omitempty"`
	DSC             []*DSC                   `xml:"dsc,omitempty" json:"dsc,omitempty"`
	CustodHist      []*FormattedNoteWithHead `xml:"custodhist" json:"custodhist,omitempty"`
	PhysTech        []*FormattedNoteWithHead `xml:"phystech,omitempty" json:"phystech,omitempty"`
	PreferCite      []*FormattedNoteWithHead `xml:"prefercite,omitempty" json:"prefercite,omitempty"`
	ProcessInfo     []*FormattedNoteWithHead `xml:"processinfo,omitempty" json:"processinfo,omitempty"`
	RelatedMaterial []*FormattedNoteWithHead `xml:"relatedmaterial,omitempty" json:"relatedmaterial,omitempty"`
	ScopeContent    []*FormattedNoteWithHead `xml:"scopecontent,omitempty" json:"scopecontent,omitempty"`
	UserRestrict    []*FormattedNoteWithHead `xml:"userrestrict,omitempty" json:"userrestrict,omitempty"`
}

type Change struct {
	Date  []*Date  `xml:"date" json:"date,omitempty"`
	Item  []*Item  `xml:"item" json:"item,omitempty"`
	Value string   `xml:",chardata" json:"value,chardata,omitempty"`
}

type ControlAccess struct {
	FamName   []*FamName   `xml:"famname" json:"famname,omitempty"`
	GenreForm []*GenreForm `xml:"genreform" json:"genreform,omitempty"`
	PersName  []*PersName  `xml:"persname" json:"persname,omitempty"`
	Subject   []*Subject   `xml:"subject" json:"subject,omitempty"`
	CorpName  []*CorpName  `xml:"corpname" json:"corpname,omitempty"`
}

type CorpName struct {
	Value string `xml:",chardata" json:"value,chardata,omitempty"`
}

type Creation struct {
	Value string  `xml:",chardata" json:"value,chardata,omitempty"`
	Date  []*Date `xml:"date" json:"date,omitempty"`
}

type Date struct {
	Value string `xml:",chardata" json:"value,chardata,omitempty"`
}

type DescRules struct {
	Value string `xml:",chardata" json:"value,chardata,omitempty"`
}

type DID struct {
	Abstract     []*Abstract     `xml:"abstract" json:"abstract,omitempty"`
	LangMaterial []*LangMaterial `xml:"langmaterial" json:"langmaterial,omitempty"`
	Origination  []*Origination  `xml:"origination" json:"origination,omitempty"`
	PhysDesc     []*PhysDesc     `xml:"physdesc" json:"physdesc,omitempty"`
	PhysLoc      []*PhysLoc      `xml:"physloc" json:"physloc,omitempty"`
	Repository   []*Repository   `xml:"repository" json:"repository,omitempty"`
	UnitTitle    []*UnitTitle    `xml:"unittitle" json:"unittitle,omitempty"`
	UnitID       []*UnitID       `xml:"unitid" json:"unitid,omitempty"`
	UnitDate     []*UnitDate     `xml:"unitdate" json:"unitdate,omitempty"`
}

type DSC struct {
	C []*C `xml:"c,omitempty" json:"c,omitempty"`
}

type EADHeader struct {
	EADID        []*EADID        `xml:"eadid" json:"eadid,omitempty"`
	FileDesc     []*FileDesc     `xml:"filedesc" json:"filedesc,omitempty"`
	ProfileDesc  []*ProfileDesc  `xml:"profiledesc" json:"profiledesc,omitempty"`
	RevisionDesc []*RevisionDesc `xml:"revisiondesc" json:"revisiondesc,omitempty"`
}

type EADID struct {
	Value          string `xml:",chardata" json:"value,chardata,omitempty"`
	CountryCode    string `xml:"countrycode,attr" json:"countrycode,attr,omitempty"`
	MainAgencyCode string `xml:"mainagencycode,attr" json:"mainagencycode,attr,omitempty"`
	URL            string `xml:"url,attr" json:"url,attr,omitempty"`
}

type Emph struct {
	Value  string `xml:",chardata" json:"value,chardata,omitempty"`
	Render string `xml:"render,attr" json:"render,attr,omitempty"`
}

type Extent struct {
	Value string `xml:",chardata" json:"value,chardata,omitempty"`
}

type ExtPtr struct {
	Href  string `xml:"href,attr" json:"href,attr,omitempty"`
	Show  string `xml:"show,attr" json:"show,attr,omitempty"`
	Title string `xml:"title,attr" json:"title,attr,omitempty"`
	Type  string `xml:"type,attr" json:"type,attr,omitempty"`
}

type FamName struct {
	Value string `xml:",chardata" json:"value,chardata,omitempty"`
}

type FileDesc struct {
	TitleStmt   TitleStmt                `xml:"titlestmt" json:"titlestmt,omitempty"`
	NoteStmt    []*FormattedNoteWithHead `xml:"notestmt" json:"notestmt,omitempty"`
	EditionStmt struct {
		P []*P `xml:"p" json:"p,omitempty"`
	} `xml:"editionstmt" json:"editionstmt,omitempty"`
}

// "eadnote" in current draft of the data model
type FormattedNoteWithHead struct {
	Head []Head `xml:"head,omitemtpy" json:"head,omitempty"`
	ID   string `xml:"id,attr" json:"id,attr,omitempty"`
	P    []*P   `xml:"p,omitempty" json:"p,omitempty"`
}

type GenreForm struct {
	Value string `xml:",chardata" json:"value,chardata,omitempty"`
}

type GeogName struct {
	Value  string `xml:",chardata" json:"value,chardata,omitempty"`
	Source string `xml:"source,attr" json:"source,attr,omitempty"`
}

type Head struct {
	Value  string    `xml:",innerxml" json:"value,chardata,omitempty"`
	ExtPtr []*ExtPtr `xml:"extptr" json:"extptr,omitempty"`
}

type Item struct {
	Value string `xml:",chardata" json:"value,chardata,omitempty"`
}

type LangMaterial struct {
	ID    	 string 	 `xml:"id,attr" json:"id,attr,omitempty"`
	Language []*Language `xml:"language" json:"language,omitempty"`
	Value 	 string 	 `xml:",chardata" json:"value,chardata,omitempty"`
}

type Language struct {
	LangCode   string `xml:"langcode,attr" json:"langcode,attr,omitempty"`
	ScriptCode string `xml:"scriptcode,attr" json:"scriptcode,attr,omitempty"`

	Value  string `xml:",chardata" json:"value,chardata,omitempty"`
}

type LangUsage struct {
	LangCode   string `xml:"langcode,attr" json:"langcode,attr,omitempty"`
}

type Occupation struct {
	Value  string `xml:",chardata" json:"value,chardata,omitempty"`
	Source string `xml:"source,attr" json:"source,attr,omitempty"`
}

type Origination struct {
	PersName []*PersName `xml:"persname" json:"persname,omitempty"`
}

type P struct {
	Value string `xml:",innerxml" json:"value,chardata,omitempty"`
}

type PersName struct {
	Value    string `xml:",chardata" json:"value,chardata,omitempty"`
	Audience string `xml:"audience,attr" json:"audience,attr,omitempty"`
	Source   string `xml:"source,attr" json:"source,attr,omitempty"`
	Rules    string `xml:"rules,attr" json:"rules,attr,omitempty"`
	Role     string `xml:"role,attr" json:"role,attr,omitempty"`
}

type PhysDesc struct {
	Extent []*Extent `xml:"extent" json:"extent,omitempty"`
}

type PhysLoc struct {
	Value string `xml:",chardata" json:"value,chardata,omitempty"`
}

type ProfileDesc struct {
	Creation  []*Creation  `xml:"creation" json:"creation,omitempty"`
	LangUsage []*LangUsage `xml:"langusage" json:"langusage,omitempty"`
	DescRules []*DescRules `xml:"descrules" json:"descrules,omitempty"`
}

type Repository struct {
	CorpName []*CorpName `xml:"corpname" json:"corpname,omitempty"`
	Emph     []*Emph     `xml:"emph" json:"emph,omitempty"`
}

type RevisionDesc struct {
	Change Change `xml:"change" json:"change,omitempty"`
}

type Subject struct {
	Value  string `xml:",chardata" json:"value,chardata,omitempty"`
	Source string `xml:"source,attr" json:"source,attr,omitempty"`
}

type TitleProper struct {
	Value string   `xml:",chardata" json:"value,chardata,omitempty"`
	Num   []string `xml:"num" json:"num,omitempty"`
	Emph  []*Emph  `xml:"emph" json:"emph,omitempty"`
	Lb    []string `xml:"lb" json:"lb,omitempty"`
}

type TitleStmt struct {
	TitleProper []*TitleProper `xml:"titleproper" json:"titleproper,omitempty"`
	Author      []string       `xml:"author" json:"author,omitempty"`
	Sponsor     []string       `xml:"sponsor" json:"sponsor,omitempty"`
}

type UnitDate struct {
	Value string `xml:",chardata" json:"value,chardata,omitempty"`
}

type UnitID struct {
	Value string `xml:",chardata" json:"value,chardata,omitempty"`
}

type UnitTitle struct {
	Value string `xml:",chardata" json:"value,chardata,omitempty"`
}

func (head *Head) MarshalJSON() ([]byte, error) {
	result, err := getMarshalledJSONForTextWithTags(head.Value)
	if err != nil {
		return nil, err
	}

	return (result), nil
}

func (p *P) MarshalJSON() ([]byte, error) {
	result, err := getMarshalledJSONForTextWithTags(p.Value)
	if err != nil {
		return nil, err
	}

	return (result), nil
}

func getMarshalledJSONForTextWithTags(text string) ([]byte, error) {
	decoder := xml.NewDecoder(strings.NewReader(text))

	var result string
	for {
		token, err := decoder.Token()

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		switch token := token.(type) {
		case xml.StartElement:
			var spanClasses string
			if token.Name.Local == "emph" {
				var render string
				for i := range token.Attr {
					if token.Attr[i].Name.Local == "render" {
						render = token.Attr[i].Value
						break
					}
				}
				spanClasses = "emph emph-" + render
			} else {
				spanClasses = token.Name.Local
			}

			result += fmt.Sprintf("<span class=\"%s\">", spanClasses)
		case xml.EndElement:
			result += "</span>"
		case xml.CharData:
			result += string(token)
		}
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
