## Technical Notes

#### The EAD validation functionality in this package requires that  
     the following URL is online and serves the xlink.xsd file :  
	 http://ead-publisher.dlib.nyu.edu/support/xsd/xlink.xsd  


#### Code layout:
* `ead.go` : contains EAD types
* `generate.go` : generation code
* how to handle providing arrays of child elements?
  could use MarshalJSON to output array
  but would be better to load array on the way in...
  (again, this is where stream parsing comes to mind...)

MAPPING:
runinfo: 
	libversion
	timestamp
	sourcefile
		 
archdesc:
	level
	accessrestrict [
		id
		head 
			value
		legalstatus
			value
			id
		list [
			numeration
			type
			head
				value
			item [
				value
					bibref [
						value
						title [
							value
							]
							
		
ARRAYIFICATION:
	option 1.) when hit element that needs to be arrayified, 
	           assemble all elements into an array
			   THIS HAPPENS DURING THE PARSING OF THE PARENT ELEMENT

	option 2.) dunno
	

|---------------|-----------------------|--------|
| element       | parent                | repeat |
|---------------|-----------------------|--------|
| address       | p                     | TRUE   |
| address       | publicationstmt       | FALSE  |
|---------------|-----------------------|--------|
| chronlist     | p                     | TRUE   |
| chronlist     | formattednotewithhead | FALSE  |
|---------------|-----------------------|--------|
| controlaccess | c                     | TRUE   |
| controlaccess | archdesc              | FALSE  |
|---------------|-----------------------|--------|
| corpname      | controlaccess         | TRUE   |
| corpname      | item                  | TRUE   |
| corpname      | p                     | TRUE   |
| corpname      | unittitle             | TRUE   |
| corpname      | indexentry            | FALSE  |
| corpname      | origination           | FALSE  |
| corpname      | repository            | FALSE  |
|---------------|-----------------------|--------|
| date          | p                     | TRUE   |
| date          | change                | FALSE  |
| date          | chronitem             | FALSE  |
| date          | creation              | FALSE  |
|---------------|-----------------------|--------|
| famname       | controlaccess         | TRUE   |
| famname       | origination           | FALSE  |
|---------------|-----------------------|--------|
| item          | list                  | TRUE   |
| item          | change                | FALSE  |
| item          | defitem               | FALSE  |
|---------------|-----------------------|--------|
| name          | item                  | TRUE   |
| name          | p                     | TRUE   |
| name          | unittitle             | TRUE   |
| name          | indexentry            | FALSE  |
|---------------|-----------------------|--------|
| persname      | controlaccess         | TRUE   |
| persname      | item                  | TRUE   |
| persname      | p                     | TRUE   |
| persname      | unittitle             | TRUE   |
| persname      | origination           | FALSE  |
|---------------|-----------------------|--------|
| physloc       | did                   | TRUE   |
| physloc       | archref               | FALSE  |
|---------------|-----------------------|--------|
| subject       | controlaccess         | TRUE   |
| subject       | p                     | TRUE   |
| subject       | indexentry            | FALSE  |
|---------------|-----------------------|--------|
| title         | abstract              | TRUE   |
| title         | bibref                | TRUE   |
| title         | controlaccess         | TRUE   |
| title         | item                  | TRUE   |
| title         | p                     | TRUE   |
| title         | unittitle             | TRUE   |
| title         | event                 | FALSE  |
|---------------|-----------------------|--------|


Try go playground for parsing 
	<publicationstmt>
        <publisher>Tamiment Library and Robert F. Wagner Labor Archives</publisher>
        <p><date>March 2021</date></p>
        <address>
          <addressline>Elmer Holmes Bobst Library</addressline>
          <addressline>70 Washington Square South</addressline>
          <addressline>2nd Floor</addressline>
          <addressline>New York, NY 10012</addressline>
          <addressline>special.collections@nyu.edu</addressline>
          <addressline>URL: <extptr
              xlink:href="http://library.nyu.edu/about/collections/special-collections-and-archives/special-collections/"
              xlink:show="new"
              xlink:title="http://library.nyu.edu/about/collections/special-collections-and-archives/special-collections/"
              xlink:type="simple"/></addressline>
        </address>
      </publicationstmt>
	  
	  
	  
REFERENCES:
https://stackoverflow.com/questions/30256729/how-to-traverse-through-xml-data-in-golang
	https://play.golang.org/p/d9BkGclp-1

NYUDLTS:
	publicationstmt example: non-array code:
		https://play.golang.org/p/palkkIZXiCl
	publicationstmt example: ARRAY code:
		https://play.golang.org/p/Hx-mWF0ROEh
		
	publicationstmt example: ARRAY code, MULTIPLE addresses:
		https://play.golang.org/p/ThimHDWuMvV	



COULD HAVE A CONVERSION FUNTION IN TITLE PROPER, AND MAPPING IN
TITLESTMT THAT DOESN'T HAVE A MAPPING FOR THE TITLEPROPER ARRAY, BUT
DOES HAVE A MAPPING FOR THE CONVERTED TITLEPROPER STRING!

E.G., 
type TitleStmt struct {
        Author      FilteredString `xml:"author" json:"author,omitempty"`
        Sponsor     FilteredString `xml:"sponsor" json:"sponsor,omitempty"`
        SubTitle    FilteredString `xml:"subtitle" json:"subtitle,omitempty"`
        TitleProper []*TitleProper `xml:"titleproper""`
        TitleProperFlattened FilteredString `json:"titleproper,omitempty"`
}


and then update the MarshalJSON to update the TitleProperFlattened
variable before outputing the data.

This would allow you to get rid of the extra Type.
