# Changelog

#### v0.3.0
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
  - replace instances of `\n` with spaces when processing `xml.CharData`
  
