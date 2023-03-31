# CHANGELOG

#### v0.23.1
  - Remove unneeded `UnitDate` code after successfully moving  
    the functionality to the `FASB`

#### v0.23.0
  - Update `DAO MarshalJSON()` as follows:  
     if `DAO.Role` is empty,  
       and `DAO.Href` IS     a valid URL then `role = "external-link"`  
       and `DAO.Href` IS NOT a valid URL then `role = "non-url"`  

  - remove `UnitDate` custom marshaling  
    Remove custom marshaling of `UnitDate` because the functionality  
      required to properly process the `UnitDate` values is better  
      performed in the `FASB`.

#### v0.22.0
  - add `Normal` to `UnitDate`
  - add custom marshaling for `UnitDate` using the following logic:
    - if there is a `Value`, then output the `Value`  
	  else if the `<unitdate @normal>` attribute is populated,   
	  then use the `@normal` attribute value converted as follows:  
      `@normal="dateA/dateB"`  
        if `dateA == dateB` then `Value = "dateA"`  
        if `dateA != dateB` then `Value = "dateA-dateB"`

#### v0.21.0
  - add `ArchRef` and `TitleValue` to `ExtRef` type
  - add `Date` to `Creation` type to align with [FADESIGN-29 data model](https://github.com/nyudlts/fadesign_29-data-model/blob/main/models.csv)
  - convert `ControlAccess.Title` to `[]*Title` instead of `[]*AccessTermWithRole`
  - add `DateChar` to `UnitDate`

#### v0.20.0
  - Update `FormattedNoteWithHead` data type to better support stream
    parsing when there are presentation elements, (e.g., `<emph>`, `<lb>`), 
	embedded in regular XML element values.
  - Tweak EADID-validation error message
  - Add custom `MarshalJSON()` function for `FormattedNoteWithHead` type.
	If a `FormattedNoteWithHead` variable has `Children`, then the variable
	will be marshaled as usual. If the `Children` slice is empty, then the 
	`innerxml` captured during parsing will be flattend/stringified, loaded
	into an `EADChild` variable, and added to the `FormattedNoteWithHead`
	`Children` slice.  This updated variable is then marshaled as usual.

#### v0.19.0
  - Add `ArchDesc.Index`, `IndexEntry.Title` and `IndexEntry.Ref`
  - Add support for multiple `<language>` child elements of
    `<langmaterial>` and `<langusage>`
  - Change parsing and JSON marshaling of `<titlestmt>` children
    `<author>, <sponsor>, and <subtitle>`.  The new implementation
    captures the `innerxml` for the child elements during parsing, and
    flattens and converts any presentation elements in the `innerxml`
    to `<span>` tags during JSON marshaling.
  - implement finalized EAD validation criteria:
	- EADs must be valid XML
	- EADs must pass validation against the EAD 2002 schema
	- The EAD ID field must:
	  - have at least two character groups, separated by an `_`
	  - the characters in the groups must be from the set `0-9a-z`
	  - the maximum length of an EAD ID is 251 characters
    - the `<did><repository><corpname>` value must be from a
      controlled vocabulary
    - the EAD top level element must be `<archdesc level="collection">`
    - the EAD must not contain any elements with the `@audience="internal"`
    - the EAD file must be smaller than 100 MB in size
  - add `ValidateEADFromFilePath()` to allow for EAD file-size checks

#### v0.18.1
  - security patch: upgrade `golang.org/x/net/http2` to `v0.7.0`

#### v0.18.0
  - add `DAO.Width` and `DAO.Height`

#### v0.17.0
  - change `<extref>` element stringification so that instead of a
	`<span class="ead-extref"...>` element (that is later escaped) the
	`<extref>` is converted to `<a class="ead-extref" href="..." target="...">`

  - Add `ExtRef` to the following types to conform with the [FADESIGN data model](https://github.com/nyudlts/fadesign_29-data-model/blob/main/models.csv):
	- `FormattedNoteWithHead`
	- `Item`
	- `PhysFacet`
	- `PhysLoc`

  - Correct typo:
	- change
	  ```
	  AltFormAvailable  []*FormattedNoteWithHead `xml:"altformavailable" json:"altformavailable,omitempty"`
	  ```
	  to
	  ```
		AltFormAvail  []*FormattedNoteWithHead `xml:"altformavail" json:"altformavail,omitempty"`
	  ```
  - Add `AppVersion` and `SourceFileHash` to `RunInfo` type
  - Remove `SetRunInfo()`

#### v0.16.0
  - correct XML and JSON tag errors: `userrestrict` --> `userestrict`
  - implement stream parsing for the following types:
	- `Bibliography`
	- `FormattedNoteWithHead`
  - add units, if present, to `Extent` during `JSON` marshaling
  - add `ReposID` to `PubInfo` type
  - add valid-HREF assertion to `ValidateEAD()`
  - accessor functions:
	- `GuideTitle()`
	- `TitleProper()`
	- `ThemeID()`
	- `RepoID()`
	- `EADID()`
  - add `DAOInfo` type, `InitDAOInfo()` and accessor functions:
	- `AllDAOCount()`
	- `AudioDAOCount()`
	- `VideioDAOCount()`
	- `ImageDAOCount()`
	- `ExternalLinkDAOCount()`
	- `ElectronicRecordsReadingRoomDAOCount()`
	- `AudioReadingRoomDAOCount()`
	- `VideoReadingRoomDAOCount()`
  - add `DAOGrpInfo` type, `InitDAOGrpInfo()` and accessor functions:
	- `AllDAOGrpCount()`
  - update `SetPubInfo()` to accept `themeID` and `repoID` arguments
  - update `ControlAccess` type, change the following members to `[]*AccessTermWithRole`:
	- `Function`
	- `GenreForm`
	- `GeogName`
	- `Occupation`
	- `Subject`
	- `Title`
  - update `Index` type, add member `P []*P`
  - make `CountDIDDAO*` and `CountDAOGrp*` functions public
  - add functions `DAOInfo.Clear()` and `DAOGrpInfo.Clear()`
  - add presentation container functionality via `InitPresentationContainers()`

#### v0.15.2
  - change the xlink.xsd `schemaLocation` in the local EAD 2002 schema
	to use a DLTS URL to avoid rate-limiting behavior observed when
	pulling the `xlink.xsd` schema from the Library of Congress
	server (https://www.loc.gov/standards/xlink/xlink.xsd).

#### v0.15.1
  - add in Free() calls to validation functions

#### v0.15.0
  - add in EAD validation functionality against both project-specific
	criteria and the EAD 2002 schema

#### v0.14.0
  - update to use Go modules

#### v0.13.0
  - add DOType and Count members to DAO type

#### v0.12.0
  - set DAO Role to "external-link" if DAO Role is empty

#### v0.11.0
  - add FilteredString String() function

#### v0.10.1
  - change Donors type to []FilteredString (from []string)

#### v0.10.0
  - add Donors member to EAD struct

#### v0.9.0
  - flatten TitleProper from an array into a single FilteredString

#### v0.8.0
  - add PubInfo type
  - add PubInfo member to EAD struct

#### v0.7.0
  - remove bracketed text from Label values
  - remove leading/trailing spaces from FilteredString values

#### v0.6.0
  - change Address in PublicationStmt to an array
  - change ChronList in FormattedNoteWithHead to an array
  - change ControlAccess in ArchDesc to an array
  - change CorpName in IndexEntry to an array
  - change CorpName in Origination to an array
  - change CorpName in Repository to an array
  - change Date in Change to an array
  - change Date in ChronItem to an array
  - change FamName in Origination to an array
  - change Item in Change to an array
  - change Item in DefItem to an array
  - change Name in IndexEntry to an array
  - change PersName in Origination to an array
  - change PhysLoc in ArchRef to an array
  - change Subject in IndexEntry to an array
  - change Title in Event to an array

#### v0.5.0
  - replace all instances of \r, \t, \n, and consecutive spaces in
	EAD element values with a single space

#### v0.4.0
  - add RunInfo.SourceFile to record the source EAD file path

#### v0.3.0
  - add FilteredString type to strip out newlines from string fields
  - add RunInfo type to capture JSON-creation timestamp and EAD package version
  - add P.PersName
  - add P.GeogName
  - add P.ChronList
  - remove P.ID field
  - remove Head.ExtPtr
  - add Extref.Actuate
  - correct parsing tag for Extent.Unit (it is an attr).
  - remove Head from DSC per data model v8.0.1
  - remove ID   from DID per data model v8.0.1
  - add AltFormAvailable to type C
  - rename AltFormatAvailable to AltFormAvailable
  - remove Abtract.Label to reflect updated data model
  - rename `AltFormatAvailable` to `AltFormAvailable`, correct XML tag, JSON tag
  - add `Date` field to `Creation` struct (matches data model v8.0.1)

#### v0.2.0
  - replace instances of `\n` with spaces in `value` fields processed by `_getConvertedTextWithTags`
