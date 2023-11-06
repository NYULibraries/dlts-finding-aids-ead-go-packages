## Technical Notes

#### DLTS-hosted `xlink.xsd` file:

In order to prevent rate-limiting behavior by the Library of Congress  
encountered when validating thousands of EADs, the `validation` command  
uses a [modified `EAD 2002` schema](./ead/validate/schema/ead-2002-20210412-dlts.xsd) that references an `xlink.xsd` file  
hosted by NYU DLTS.

Therefore, the EAD validation functionality requires that the following  
URL is online: http://dlts-support-files.s3.amazonaws.com/findingaids/xsd/xlink.xsd  

The EAD 2002 schema is included in the executable using the  
Golang Embedded FS functionality.  Unfortunately, I (jgpawletko) was unable  
to figure out a way to embed the `xlink.xsd` schema as well.  I think this is  
because the `libxml2` C code performs the validation, and does not have direct  
access to the file system embedded in the executable.  

Going forward, the tradeoffs between a single-file deployment that requires that  
the `xlink.xsd` URL is available vs. a multi-file deployment where the schema files  
are deployed to the filesystem should be re-evaluated.  


#### Selective stream parsing:

The Golang XML parser does not respect the order of elements encountered in the EAD files.  
This is a [known issue](https://pkg.go.dev/encoding/xml#pkg-note-BUG):   
> Mapping between XML elements and data structures is inherently flawed:   
> an XML element is an order-dependent collection of anonymous values,  
> while a data structure is an order-independent collection of named values.  
> ... 

To preserve order when it was important, selective stream parsing was implemented.  
Please see [here](./ead/ead_decoder.go) for the implementation.  



