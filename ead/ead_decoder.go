package ead

import (
	"encoding/xml"
	"fmt"
)

type EADChild struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value,omitempty"`
}

func (eadChild *EADChild) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	name := start.Name.Local
	switch name {
	case "accessrestrict", "accruals", "acqinfo", "altformavail", "appraisal", "arrangement", "bioghist",
		"custodhist", "odd", "otherfindaid", "originalsloc", "phystech", "prefercite",
		"processinfo", "relatedmaterial", "scopecontent", "separatedmaterial", "userestrict":
		e := FormattedNoteWithHead{}
		return decodeElement(eadChild, &e, d, start)
	case "bibliography":
		e := Bibliography{}
		return decodeElement(eadChild, &e, d, start)
	case "bibref":
		e := BibRef{}
		return decodeElement(eadChild, &e, d, start)
	case "controlaccess":
		e := ControlAccess{}
		return decodeElement(eadChild, &e, d, start)
	case "chronlist":
		e := ChronList{}
		return decodeElement(eadChild, &e, d, start)
	case "defitem":
		e := DefItem{}
		return decodeElement(eadChild, &e, d, start)
	case "did":
		e := DID{}
		return decodeElement(eadChild, &e, d, start)
	case "dsc":
		e := DSC{}
		return decodeElement(eadChild, &e, d, start)
	case "extref":
		e := ExtRef{}
		return decodeElement(eadChild, &e, d, start)
	case "legalstatus":
		e := LegalStatus{}
		return decodeElement(eadChild, &e, d, start)
	case "list":
		e := List{}
		return decodeElement(eadChild, &e, d, start)
	case "p":
		e := P{}
		return decodeElement(eadChild, &e, d, start)
	default:
		return fmt.Errorf("unsupported element error: %s", name)
	}
}

func decodeElement(eadChild *EADChild, strct any, d *xml.Decoder, start xml.StartElement) error {
	if err := d.DecodeElement(&strct, &start); err != nil {
		return err
	}
	eadChild.Name = start.Name.Local
	eadChild.Value = strct
	return nil
}
