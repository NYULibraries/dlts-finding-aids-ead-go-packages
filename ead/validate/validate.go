package validate

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/nyulibraries/dlts-finding-aids-ead-go-packages/ead"
)

const ValidEADIDRegexpString = "^[a-z0-9]+(?:_[a-z0-9]+){1,7}$"

// Map elements are from
// https://github.com/nyudlts/fadesign_29-data-model/blob/1fa578e90a0239431154d150fe053f7982795e10/relator-authoritative-labels.csv
var RelatorAuthoritativeLabelMap = map[string]string{
	"abr": "Abridger",
	"acp": "Art copyist",
	"act": "Actor",
	"adi": "Art director",
	"adp": "Adapter",
	"aft": "Author of afterword, colophon, etc.",
	"anl": "Analyst",
	"anm": "Animator",
	"ann": "Annotator",
	"ant": "Bibliographic antecedent",
	"ape": "Appellee",
	"apl": "Appellant",
	"app": "Applicant",
	"aqt": "Author in quotations or text abstracts",
	"arc": "Architect",
	"ard": "Artistic director",
	"arr": "Arranger",
	"art": "Artist",
	"asg": "Assignee",
	"asn": "Associated name",
	"ato": "Autographer",
	"att": "Attributed name",
	"auc": "Auctioneer",
	"aud": "Author of dialog",
	"aui": "Author of introduction, etc.",
	"aus": "Screenwriter",
	"aut": "Author",
	"bdd": "Binding designer",
	"bjd": "Bookjacket designer",
	"bkd": "Book designer",
	"bkp": "Book producer",
	"blw": "Blurb writer",
	"bnd": "Binder",
	"bpd": "Bookplate designer",
	"brd": "Broadcaster",
	"brl": "Braille embosser",
	"bsl": "Bookseller",
	"cas": "Caster",
	"ccp": "Conceptor",
	"chr": "Choreographer",
	"cli": "Client",
	"cll": "Calligrapher",
	"clr": "Colorist",
	"clt": "Collotyper",
	"cmm": "Commentator",
	"cmp": "Composer",
	"cmt": "Compositor",
	"cnd": "Conductor",
	"cng": "Cinematographer",
	"cns": "Censor",
	"coe": "Contestant-appellee",
	"col": "Collector",
	"com": "Compiler",
	"con": "Conservator",
	"cor": "Collection registrar",
	"cos": "Contestant",
	"cot": "Contestant-appellant",
	"cou": "Court governed",
	"cov": "Cover designer",
	"cpc": "Copyright claimant",
	"cpe": "Complainant-appellee",
	"cph": "Copyright holder",
	"cpl": "Complainant",
	"cpt": "Complainant-appellant",
	"cre": "Creator",
	"crp": "Correspondent",
	"crr": "Corrector",
	"crt": "Court reporter",
	"csl": "Consultant",
	"csp": "Consultant to a project",
	"cst": "Costume designer",
	"ctb": "Contributor",
	"cte": "Contestee-appellee",
	"ctg": "Cartographer",
	"ctr": "Contractor",
	"cts": "Contestee",
	"ctt": "Contestee-appellant",
	"cur": "Curator",
	"cwt": "Commentator for written text",
	"dbp": "Distribution place",
	"dfd": "Defendant",
	"dfe": "Defendant-appellee",
	"dft": "Defendant-appellant",
	"dgg": "Degree granting institution",
	"dgs": "Degree supervisor",
	"dis": "Dissertant",
	"dln": "Delineator",
	"dnc": "Dancer",
	"dnr": "Donor",
	"dpc": "Depicted",
	"dpt": "Depositor",
	"drm": "Draftsman",
	"drt": "Director",
	"dsr": "Designer",
	"dst": "Distributor",
	"dtc": "Data contributor",
	"dte": "Dedicatee",
	"dtm": "Data manager",
	"dto": "Dedicator",
	"dub": "Dubious author",
	"edc": "Editor of compilation",
	"edm": "Editor of moving image work",
	"edt": "Editor",
	"egr": "Engraver",
	"elg": "Electrician",
	"elt": "Electrotyper",
	"eng": "Engineer",
	"enj": "Enacting jurisdiction",
	"etr": "Etcher",
	"evp": "Event place",
	"exp": "Expert",
	"fac": "Facsimilist",
	"fds": "Film distributor",
	"fld": "Field director",
	"flm": "Film editor",
	"fmd": "Film director",
	"fmk": "Filmmaker",
	"fmo": "Former owner",
	"fmp": "Film producer",
	"fnd": "Funder",
	"fpy": "First party",
	"frg": "Forger",
	"gis": "Geographic information specialist",
	"his": "Host institution",
	"hnr": "Honoree",
	"hst": "Host",
	"ill": "Illustrator",
	"ilu": "Illuminator",
	"ins": "Inscriber",
	"inv": "Inventor",
	"isb": "Issuing body",
	"itr": "Instrumentalist",
	"ive": "Interviewee",
	"ivr": "Interviewer",
	"jud": "Judge",
	"jug": "Jurisdiction governed",
	"lbr": "Laboratory",
	"lbt": "Librettist",
	"ldr": "Laboratory director",
	"led": "Lead",
	"lee": "Libelee-appellee",
	"lel": "Libelee",
	"len": "Lender",
	"let": "Libelee-appellant",
	"lgd": "Lighting designer",
	"lie": "Libelant-appellee",
	"lil": "Libelant",
	"lit": "Libelant-appellant",
	"lsa": "Landscape architect",
	"lse": "Licensee",
	"lso": "Licensor",
	"ltg": "Lithographer",
	"lyr": "Lyricist",
	"mcp": "Music copyist",
	"mdc": "Metadata contact",
	"med": "Medium",
	"mfp": "Manufacture place",
	"mfr": "Manufacturer",
	"mod": "Moderator",
	"mon": "Monitor",
	"mrb": "Marbler",
	"mrk": "Markup editor",
	"msd": "Musical director",
	"mte": "Metal-engraver",
	"mtk": "Minute taker",
	"mus": "Musician",
	"nrt": "Narrator",
	"opn": "Opponent",
	"org": "Originator",
	"orm": "Organizer",
	"osp": "Onscreen presenter",
	"oth": "Other",
	"own": "Owner",
	"pan": "Panelist",
	"pat": "Patron",
	"pbd": "Publishing director",
	"pbl": "Publisher",
	"pdr": "Project director",
	"pfr": "Proofreader",
	"pht": "Photographer",
	"plt": "Platemaker",
	"pma": "Permitting agency",
	"pmn": "Production manager",
	"pop": "Printer of plates",
	"ppm": "Papermaker",
	"ppt": "Puppeteer",
	"pra": "Praeses",
	"prc": "Process contact",
	"prd": "Production personnel",
	"pre": "Presenter",
	"prf": "Performer",
	"prg": "Programmer",
	"prm": "Printmaker",
	"prn": "Production company",
	"pro": "Producer",
	"prp": "Production place",
	"prs": "Production designer",
	"prt": "Printer",
	"prv": "Provider",
	"pta": "Patent applicant",
	"pte": "Plaintiff-appellee",
	"ptf": "Plaintiff",
	"pth": "Patent holder",
	"ptt": "Plaintiff-appellant",
	"pup": "Publication place",
	"rbr": "Rubricator",
	"rcd": "Recordist",
	"rce": "Recording engineer",
	"rcp": "Addressee",
	"rdd": "Radio director",
	"red": "Redaktor",
	"ren": "Renderer",
	"res": "Researcher",
	"rev": "Reviewer",
	"rpc": "Radio producer",
	"rps": "Repository",
	"rpt": "Reporter",
	"rpy": "Responsible party",
	"rse": "Respondent-appellee",
	"rsg": "Restager",
	"rsp": "Respondent",
	"rsr": "Restorationist",
	"rst": "Respondent-appellant",
	"rth": "Research team head",
	"rtm": "Research team member",
	"sad": "Scientific advisor",
	"sce": "Scenarist",
	"scl": "Sculptor",
	"scr": "Scribe",
	"sds": "Sound designer",
	"sec": "Secretary",
	"sgd": "Stage director",
	"sgn": "Signer",
	"sht": "Supporting host",
	"sll": "Seller",
	"sng": "Singer",
	"spk": "Speaker",
	"spn": "Sponsor",
	"spy": "Second party",
	"srv": "Surveyor",
	"std": "Set designer",
	"stg": "Setting",
	"stl": "Storyteller",
	"stm": "Stage manager",
	"stn": "Standards body",
	"str": "Stereotyper",
	"tcd": "Technical director",
	"tch": "Teacher",
	"ths": "Thesis advisor",
	"tld": "Television director",
	"tlp": "Television producer",
	"trc": "Transcriber",
	"trl": "Translator",
	"tyd": "Type designer",
	"tyg": "Typographer",
	"uvp": "University place",
	"vac": "Voice actor",
	"vdg": "Videographer",
	"wac": "Writer of added commentary",
	"wal": "Writer of added lyrics",
	"wam": "Writer of accompanying material",
	"wat": "Writer of added text",
	"wdc": "Woodcutter",
	"wde": "Wood engraver",
	"win": "Writer of introduction",
	"wit": "Witness",
	"wpr": "Writer of preface",
	"wst": "Writer of supplementary textual content",
}

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

func ValidateEAD(data []byte) ([]string, error) {
	var validationErrors = []string{}

	validationErrors = append(validationErrors, validateXML(data)...)

	var ead ead.EAD
	err := xml.Unmarshal(data, &ead)
	if err != nil {
		return validationErrors, err
	}

	validationErrors = append(validationErrors, validateRequiredEADElements(ead)...)
	validationErrors = append(validationErrors, validateRepository(ead)...)
	validationErrors = append(validationErrors, validateEADID(ead)...)

	validateNoUnpublishedMaterialValidationErrors, err := validateNoUnpublishedMaterial(data)
	if err != nil {
		return validationErrors, err
	}

	validationErrors = append(validationErrors, validateNoUnpublishedMaterialValidationErrors...)
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

	validationErrors = append(validationErrors, makeAudienceInternalErrorMessage(audienceInternalElements))

	return validationErrors, nil
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
func validateXML(data []byte) []string {
	var validationErrors = []string{}

	// Not perfect, but maybe good enough for now.
	if xml.Unmarshal(data, new(interface{})) != nil {
		validationErrors = append(validationErrors, makeInvalidXMLErrorMessage())
	}

	return validationErrors
}
