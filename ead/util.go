package ead

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
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

func GetConvertedTextWithTags(text string) ([]byte, error) {
	return getConvertedTextWithTags(text)
}
func getConvertedTextWithTags(text string) ([]byte, error) {
	return _getConvertedTextWithTags(text, true)
}

func GetConvertedTextWithTagsNoLBConversion(text string) ([]byte, error) {
	return getConvertedTextWithTagsNoLBConversion(text)
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
			case "extref":
				{
					var href string
					var target string

					for i := range token.Attr {
						if token.Attr[i].Name.Local == "href" {
							href = token.Attr[i].Value
						}
						if token.Attr[i].Name.Local == "show" {
							target = token.Attr[i].Value
						}
					}
					result += fmt.Sprintf("<a class=\"%s\" href=\"%s\" target=\"%s\">", "ead-extref", href, target)
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

			if token.Name.Local == "extref" {
				result += "</a>"
				needClosingTag = false
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
		// if relator code is not recognized, trust the archivists and drop through
		return relatorID, nil
	}
}

// RunInfo stores data related to the parsing/JSON generation process
type RunInfo struct {
	PkgVersion     string    `json:"libversion"`
	AppVersion     string    `json:"appversion,omitempty"`
	TimeStamp      time.Time `json:"timestamp"`
	SourceFile     string    `json:"sourcefile"`
	SourceFileHash string    `json:"sourcefilehash,omitempty"`
}

// DAOInfo stores data related to the digital objects in the parsed EAD
// https://jira.nyu.edu/browse/FADESIGN-138
type DAOInfo struct {
	AllDAOCount                       uint32
	AudioCount                        uint32
	VideoCount                        uint32
	ImageCount                        uint32
	ExternalLinkCount                 uint32
	ElectronicRecordsReadingRoomCount uint32
	AudioReadingRoomCount             uint32
	VideoReadingRoomCount             uint32
	AllDAOs                           []*DAO
	AudioDAOs                         []*DAO
	VideoDAOs                         []*DAO
	ImageDAOs                         []*DAO
	ExternalLinkDAOs                  []*DAO
	ElectronicRecordsReadingRoomDAOs  []*DAO
	AudioReadingRoomDAOs              []*DAO
	VideoReadingRoomDAOs              []*DAO
}

type DAOGrpInfo struct {
	AllDAOGrpCount uint32
	AllDAOGrps     []*DAOGrp
}

// Donors is slice containing Donor names
type Donors []FilteredString

// PubInfo stores data used by the publication system
type PubInfo struct {
	ThemeID string `json:"themeid"`
	RepoID  string `json:"reposidentifier"`
}

func (p *PubInfo) SetPubInfo(themeid string, repoid string) {
	p.ThemeID = themeid
	p.RepoID = repoid
}

// FilteredString provides a centralized string cleanup mechanism
type FilteredString string

func (s FilteredString) String() string {
	return cleanupWhitespace(string(s))
}

func (s FilteredString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func cleanupWhitespace(inputString string) string {
	s := html.UnescapeString(inputString)

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

func flattenCDATA(cdata CDATA) ([]byte, error) {
	return getConvertedTextWithTags(cdata.Value)
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
		return nil, fmt.Errorf("unable to find correct title")
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

func (e *EAD) RepoID() string {
	return e.PubInfo.RepoID
}

func (e *EAD) InitDAOCounts() {
	CountDIDDAOs(&e.ArchDesc.DID, &e.DAOInfo)
	CountCsDAOs(e.ArchDesc.DSC.C, &e.DAOInfo)
}

func (e *EAD) EADID() string {
	return e.EADHeader.EADID.Value
}

// process an array of containers
func CountCsDAOs(cs []*C, daoInfo *DAOInfo) {
	for _, c := range cs {
		CountCDAOs(c, daoInfo)
	}
}

// process a container
func CountCDAOs(c *C, daoInfo *DAOInfo) {
	// configured to perform a breadth-first aggregation.
	// switch the CountDIDDAOs/CountCsDAOs order if you want to do
	//   a depth-first aggregation.
	CountDIDDAOs(&c.DID, daoInfo)
	CountCsDAOs(c.C, daoInfo)
}

func CountDIDDAOs(did *DID, daoInfo *DAOInfo) {
	daos := did.DAO

	// https://jira.nyu.edu/browse/FADESIGN-138
	for _, dao := range daos {
		// init parent pointer
		dao.ParentDID = did

		// collect all DAOs
		daoInfo.AllDAOCount += 1
		appendDAO(dao, &daoInfo.AllDAOs)

		switch dao.Role {
		case "audio-service":
			daoInfo.AudioCount += 1
			appendDAO(dao, &daoInfo.AudioDAOs)
		case "video-service":
			daoInfo.VideoCount += 1
			appendDAO(dao, &daoInfo.VideoDAOs)
		case "image-service":
			daoInfo.ImageCount += 1
			appendDAO(dao, &daoInfo.ImageDAOs)
		case "external-link":
			daoInfo.ExternalLinkCount += 1
			appendDAO(dao, &daoInfo.ExternalLinkDAOs)
		case "electronic-records-reading-room":
			daoInfo.ElectronicRecordsReadingRoomCount += 1
			appendDAO(dao, &daoInfo.ElectronicRecordsReadingRoomDAOs)
		case "audio-reading-room":
			daoInfo.AudioReadingRoomCount += 1
			appendDAO(dao, &daoInfo.AudioReadingRoomDAOs)
		case "video-reading-room":
			daoInfo.VideoReadingRoomCount += 1
			appendDAO(dao, &daoInfo.VideoReadingRoomDAOs)
		default:
			// the strategy for DAOs without roles is to treat them as external links
			if len(strings.TrimSpace(string(dao.Role))) == 0 {
				daoInfo.ExternalLinkCount += 1
				appendDAO(dao, &daoInfo.ExternalLinkDAOs)
			}
		}
	}
}

func appendDAO(dao *DAO, daoSlice *[]*DAO) {
	*daoSlice = append(*daoSlice, dao)
}

func (e *EAD) AllDAOCount() uint32 {
	return e.DAOInfo.AllDAOCount
}

func (e *EAD) AudioDAOCount() uint32 {
	return e.DAOInfo.AudioCount
}

func (e *EAD) VideoDAOCount() uint32 {
	return e.DAOInfo.VideoCount
}

func (e *EAD) ImageDAOCount() uint32 {
	return e.DAOInfo.ImageCount
}

func (e *EAD) ExternalLinkDAOCount() uint32 {
	return e.DAOInfo.ExternalLinkCount
}

func (e *EAD) ElectronicRecordsReadingRoomDAOCount() uint32 {
	return e.DAOInfo.ElectronicRecordsReadingRoomCount
}

func (e *EAD) AudioReadingRoomDAOCount() uint32 {
	return e.DAOInfo.AudioReadingRoomCount
}

func (e *EAD) VideoReadingRoomDAOCount() uint32 {
	return e.DAOInfo.VideoReadingRoomCount
}

func (e *EAD) InitDAOGrpCount() {
	CountDAOGrps(e.ArchDesc.DID.DAOGrp, &e.DAOGrpInfo)
	CountCsDAOGrps(e.ArchDesc.DSC.C, &e.DAOGrpInfo)
}

func CountDAOGrps(daoGrps []*DAOGrp, daoGrpInfo *DAOGrpInfo) {
	for _, daoGrp := range daoGrps {
		daoGrpInfo.AllDAOGrpCount += 1
		appendDAOGrp(daoGrp, &daoGrpInfo.AllDAOGrps)
	}
}

// process an array of containers
func CountCsDAOGrps(cs []*C, daoGrpInfo *DAOGrpInfo) {
	for _, c := range cs {
		countCDAOGrps(c, daoGrpInfo)
	}
}

// process a container
func countCDAOGrps(c *C, daoGrpInfo *DAOGrpInfo) {
	CountCsDAOGrps(c.C, daoGrpInfo)
	CountDAOGrps(c.DID.DAOGrp, daoGrpInfo)
}

func appendDAOGrp(daoGrp *DAOGrp, daoGrpSlice *[]*DAOGrp) {
	*daoGrpSlice = append(*daoGrpSlice, daoGrp)
}

func (e *EAD) AllDAOGrpCount() uint32 {
	return e.DAOGrpInfo.AllDAOGrpCount
}

func (di *DAOInfo) Clear() {
	di.AllDAOCount = 0
	di.AudioCount = 0
	di.VideoCount = 0
	di.ImageCount = 0
	di.ExternalLinkCount = 0
	di.ElectronicRecordsReadingRoomCount = 0
	di.AudioReadingRoomCount = 0
	di.VideoReadingRoomCount = 0
	di.AllDAOs = nil
	di.AudioDAOs = nil
	di.VideoDAOs = nil
	di.ImageDAOs = nil
	di.ExternalLinkDAOs = nil
	di.ElectronicRecordsReadingRoomDAOs = nil
	di.AudioReadingRoomDAOs = nil
	di.VideoReadingRoomDAOs = nil
}

func (dgi *DAOGrpInfo) Clear() {
	dgi.AllDAOGrpCount = 0
	dgi.AllDAOGrps = nil
}

func (e *EAD) InitPresentationContainers() {
	e.ArchDesc.DSC.C = addPresentationContainers(&e.ArchDesc.DSC.C)
}

func addPresentationContainers(csp *[]*C) []*C {

	cs := *csp

	// return immediately if there aren't any containers
	if len(cs) == 0 {
		return cs
	}

	var collapsedCs []*C

	// inRun            : currently in a run of Cs to collapse
	// pcCount          : presentation container count, used to init the presentation container IDs
	// keepStartIdx     : the starting index of elements to keep
	// collapseStartIdx : the starting index of elements to collapse
	// collapseEndIdx   : the ending   index of elements to collapse

	inRun := false
	keepStartIdx := -1
	collapseStartIdx := 0
	collapseEndIdx := 0
	pcCount := 0

	for idx, c := range cs {

		//DEBUG fmt.Printf("idx: %03d, len: %03d\n", idx, len(collapsedCs))
		// is this a container to collapse?
		if shouldCollapseContainer(c) {
			//DEBUG fmt.Println("we shouldCollapseContainer")
			// are we already in a run of containers to collapse?
			if !inRun {
				//DEBUG fmt.Println("---> not in a run")
				// now we're in a run
				inRun = true
				collapseStartIdx = idx
				// if this run comes after some containers to keep...
				if keepStartIdx != -1 {
					//DEBUG fmt.Printf("------> keepStartIdx != -1")
					// need to append all of the known Cs since the last collapse
					collapsedCs = append(collapsedCs, cs[keepStartIdx:collapseStartIdx]...)
				}
			}

		} else {
			// this is a container to keep
			// have we already seen a container to keep?
			if keepStartIdx == -1 {
				// if not, store the index
				keepStartIdx = idx
			}

			// if we're currently in a run, then we have found the end of it
			if inRun {
				inRun = false
				collapseEndIdx = idx
				pc := new(C)
				pc.Level = "dl-presentation"
				pcCount += 1
				pc.ID = FilteredString(fmt.Sprintf("items%03d", pcCount))
				pc.DID.UnitTitle = &UnitTitle{Value: "View Inventory"}
				pc.C = cs[collapseStartIdx:collapseEndIdx]
				keepStartIdx = collapseEndIdx
				collapsedCs = append(collapsedCs, pc)
			}
		}
	} // end of container loop

	if inRun {
		// if we ended on a run, then we need to finish the collapse operation
		inRun = false
		pc := new(C)
		pc.Level = "dl-presentation"
		pcCount += 1
		pc.ID = FilteredString(fmt.Sprintf("items%03d", pcCount))
		pc.DID.UnitTitle = &UnitTitle{Value: "View Inventory"}
		pc.C = cs[collapseStartIdx:] // to the end of the slice because we ended on a run
		collapsedCs = append(collapsedCs, pc)
	} else {
		// did NOT end in a run, so we need to append the remaining "keep" Cs to the collapsedCs
		// need to append all of the "keep" Cs since the last collapse
		collapsedCs = append(collapsedCs, cs[keepStartIdx:]...)
	}

	// return collapsed Cs
	return collapsedCs
}

func shouldCollapseContainer(c *C) bool {
	switch c.Level {
	case "series", "otherlevel", "recordgrp", "dl-presentation":
		return false
	}
	return true
}
