## dlts-finding-aids-ead-go-packages

This repo contains code for working with Encoded Archival Description
XML (EAD) [version 2002](https://loc.gov/ead/).

There are a few core areas of functionality:
1. EAD Parsing/JSON generation:  
This package parses EAD 2002 XML files and generates two forms of JSON that are input into the [Hugo static site generator](https://gohugo.io) application:  
    a. "Intermediate JSON" (`iJSON`)  
    b. Hugo Content files (`hJSON`)  
2. EAD Validation:  
This package validates EAD XML files per the [EAD 2002 schema](https://loc.gov/ead/eadschema.html) and the [EAD Validation Criteria for Publishing](https://github.com/nyudlts/findingaids_docs/blob/main/user/EAD_Validation_Criteria_for_Publishing.pdf)  
3. "FABifying" EADs:  
This package has code that will modify an incoming EAD so that it is compatible with the ["Finding Aids Bridge" (FAB) discovery application](https://github.com/NYULibraries/specialcollections/tree/master) indexer

##### WARNING:
The major version of this package is `0`.

Per [Semantic Versioning 2.0.0](https://semver.org/#spec-item-4):
```
Major version zero (0.y.z) is for initial development. 
Anything MAY change at any time. 
The public API SHOULD NOT be considered stable.
```

Thanks!

For additional information, see the [technical notes](./TECH-NOTES.md).