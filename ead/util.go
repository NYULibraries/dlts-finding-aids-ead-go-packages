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
)

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
	unit := ""
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
			case "extent":
				{
					// if there is a "unit" attribute, save the
					// unit attribute value to append to the result
					for i := range token.Attr {
						if token.Attr[i].Name.Local == "unit" {
							unit = token.Attr[i].Value
						}

					}
					// default processing of opening tag
					result += _getConvertedTextWithTagsDefault(token.Name.Local)
				}
			}

		case xml.EndElement:
			// Add "unit" attribute value to extent if present
			if token.Name.Local == "extent" {
				if unit != "" {
					result += fmt.Sprintf(" %s", unit)
				}
				// reset unit
				unit = ""
			}

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
	if authoritativeLabel, ok := RelatorAuthoritativeLabelMap[relatorID]; ok {
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

// DAOInfo stores data related to the digital objects in the parsed EAD
// https://jira.nyu.edu/browse/FADESIGN-138
type DAOInfo struct {
	AudioCount                        uint32
	VideoCount                        uint32
	ImageCount                        uint32
	ExternalLinkCount                 uint32
	ElectronicRecordsReadingRoomCount uint32
	AudioReadingRoomCount             uint32
	VideoReadingRoomCount             uint32
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

func (e *EAD) GuideTitle() string {
	return e.ArchDesc.DID.UnitTitle.Value
}

func (e *EAD) TitleProper() string {
	flattenedTitleProper, _ := flattenTitleProper(e.EADHeader.FileDesc.TitleStmt.TitleProper)
	return string(flattenedTitleProper)
}

func (e *EAD) ThemeID() string {
	return e.PubInfo.ThemeID
}

func (e *EAD) InitDAOCounts() {
	countArchDescDAOs(e)
}

func countArchDescDAOs(e *EAD) {
	countDAOs(e.ArchDesc.DID.DAO, &e.DAOInfo)
}

func countDAOs(daos []*DAO, daoInfo *DAOInfo) {
	// https://jira.nyu.edu/browse/FADESIGN-138
	for _, dao := range daos {
		switch dao.Role {
		case "audio-service":
			daoInfo.AudioCount += 1
		case "video-service":
			daoInfo.VideoCount += 1
		case "image-service":
			daoInfo.ImageCount += 1
		case "external-link":
			daoInfo.ExternalLinkCount += 1
		case "electronic-records-reading-room":
			daoInfo.ElectronicRecordsReadingRoomCount += 1
		case "audio-reading-room":
			daoInfo.AudioReadingRoomCount += 1
		case "video-reading-room":
			daoInfo.AudioReadingRoomCount += 1
		}
	}
}
