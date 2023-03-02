package ead

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type iJSONTestParams struct {
	TestName              string
	EADFilePath           string
	JSONReferenceFilePath string
	JSONErrorFilePath     string
	PrePopulatedEAD       *EAD
}

var testFixturePath string = filepath.Join(".", "testdata")
var omegaTestFixturePath string = filepath.Join(testFixturePath, "omega", "v0.1.5")
var falesTestFixturePath string = filepath.Join(testFixturePath, "fales")
var nyhsTestFixturePath string = filepath.Join(testFixturePath, "nyhs")
var akkasahTestFixturePath string = filepath.Join(testFixturePath, "akkasah")
var presentationContainerPath string = filepath.Join(testFixturePath, "presentation-containers")

func testProcessDID(did *DID) {
	daos := did.DAO
	// populate digital object DOType and Count members
	for _, dao := range daos {
		switch dao.Href {
		case "https://hdl.handle.net/2333.1/xgxd28gq":
			dao.DOType = "image_set"
			dao.Count = 32
		case "https://hdl.handle.net/2333.1/m63xss7g":
			dao.DOType = "image_set"
			dao.Count = 6
		case "https://hdl.handle.net/2333.1/ttdz0j92":
			dao.DOType = "image_set"
			dao.Count = 3
		}
	}
}

func testProcessCs(cs []*C) {
	for _, c := range cs {
		testProcessCs(c.C)
		testProcessDID(&c.DID)
	}
}

func failOnError(t *testing.T, err error, label string) {
	if err != nil {
		t.Errorf("%s: %s", label, err)
		t.FailNow()
	}
}

func assertEqual(t *testing.T, want string, got string, label string) {
	if want != got {
		t.Errorf("%s Mismatch: want: %s, got: %s", label, want, got)
	}
}

func assertEqualUint32(t *testing.T, want uint32, got uint32, label string) {
	if want != got {
		t.Errorf("%s Mismatch: want: %d, got: %d", label, want, got)
	}
}

func assertFilteredStringSlicesEqual(t *testing.T, want []FilteredString, got []FilteredString, label string) {
	if len(want) != len(got) {
		t.Errorf("%s Mismatch: want: %v, got: %v", label, want, got)
	}
	for i := range want {
		if want[i] != got[i] {
			t.Errorf("%s Mismatch: want: %v, got: %v", label, want[i], got[i])
		}
	}
}

func getOmegaEAD(t *testing.T) EAD {
	EADXML, err := os.ReadFile(omegaTestFixturePath + "/" + "Omega-EAD.xml")
	failOnError(t, err, "Unexpected error")

	var ead EAD
	err = xml.Unmarshal([]byte(EADXML), &ead)
	failOnError(t, err, "Unexpected error")

	return ead
}

func getNYHSFoundling(t *testing.T) EAD {
	EADXML, err := os.ReadFile(nyhsTestFixturePath + "/" + "nyhs_foundling.xml")
	failOnError(t, err, "Unexpected error")

	var ead EAD
	err = xml.Unmarshal([]byte(EADXML), &ead)
	failOnError(t, err, "Unexpected error")

	return ead
}

func getPresentationContainerEAD(t *testing.T, filename string) EAD {
	EADXML, err := os.ReadFile(presentationContainerPath + "/" + filename)
	failOnError(t, err, "Unexpected error")

	var ead EAD
	err = xml.Unmarshal([]byte(EADXML), &ead)
	failOnError(t, err, "Unexpected error")

	return ead
}

func getTestEAD(t *testing.T, eadPath string) *EAD {
	EADXML, err := os.ReadFile(eadPath)
	failOnError(t, err, "Unexpected error")

	var ead EAD
	err = xml.Unmarshal([]byte(EADXML), &ead)
	failOnError(t, err, "Unexpected error")

	return &ead
}

func TestXMLParsing(t *testing.T) {
	t.Run("XML Parsing", func(t *testing.T) {
		ead := getOmegaEAD(t)

		want := "collection"
		got := string(ead.ArchDesc.Level)
		assertEqual(t, want, got, "ArchDesc.Level")
	})
}

func TestJSONMarshaling(t *testing.T) {
	var params iJSONTestParams

	params.TestName = "JSON Marshaling"
	params.EADFilePath = filepath.Join(omegaTestFixturePath, "Omega-EAD.xml")
	params.JSONReferenceFilePath = filepath.Join(omegaTestFixturePath, "mos_2021.json")
	params.JSONErrorFilePath = "./testdata/tmp/failing-marshal.json"

	runiJSONComparisonTest(t, &params)
}

func TestUpdateRunInfo(t *testing.T) {
	t.Run("Update RunInfo", func(t *testing.T) {
		var sut EAD

		want := "0001-01-01T00:00:00Z"
		got := sut.RunInfo.TimeStamp.Format(time.RFC3339)
		assertEqual(t, want, got, "Initial ead.RunInfo.TimeStamp")

		now := time.Now()
		pkgVersion := Version // from ead package constant
		appVersion := "v0.17.0"
		sourceFile := "/a/very/fine/path/to/an/ead.xml"
		sourceFileHash := "md5:9cacfec48461900f3170f3b5d69af527"

		sut.RunInfo.PkgVersion = pkgVersion
		sut.RunInfo.AppVersion = appVersion
		sut.RunInfo.TimeStamp = now
		sut.RunInfo.SourceFile = sourceFile
		sut.RunInfo.SourceFileHash = sourceFileHash

		want = pkgVersion
		got = sut.RunInfo.PkgVersion
		assertEqual(t, want, got, "Post-assignment ead.RunInfo.PkgVersion")

		want = appVersion
		got = sut.RunInfo.AppVersion
		assertEqual(t, want, got, "Post-assignment ead.RunInfo.AppVersion")

		want = now.Format(time.RFC3339)
		got = sut.RunInfo.TimeStamp.Format(time.RFC3339)
		assertEqual(t, want, got, "Post-assignment ead.RunInfo.TimeStamp")

		want = sourceFile
		got = sut.RunInfo.SourceFile
		assertEqual(t, want, got, "set ead.RunInfo.SourceFile")

		want = sourceFileHash
		got = sut.RunInfo.SourceFileHash
		assertEqual(t, want, got, "set ead.RunInfo.SourceFileHash")
	})
}

func TestUpdatePubInfo(t *testing.T) {
	t.Run("Update PubInfo", func(t *testing.T) {
		var sut EAD
		var want, got string

		want = ""
		got = sut.PubInfo.ThemeID
		assertEqual(t, want, got, "Initial ead.PubInfo.ThemeID")

		want = ""
		got = sut.PubInfo.RepoID
		assertEqual(t, want, got, "Initial ead.PubInfo.RepoID")

		themeid := "cdf80c84-2655-4a01-895d-fbf9a374c1df"
		repoid := "9d396ffa-1b3e-41f0-8bc9-e101a5a828bc"
		sut.PubInfo.SetPubInfo(themeid, repoid)

		want = themeid
		got = sut.PubInfo.ThemeID
		assertEqual(t, want, got, "Post-assignment ead.PubInfo.ThemeID")

		want = repoid
		got = sut.PubInfo.RepoID
		assertEqual(t, want, got, "Post-assignment ead.PubInfo.RepoID")

	})
}

func TestBarcodeRemovalFromLabels(t *testing.T) {
	var params iJSONTestParams

	params.TestName = "Barcode Removal from Labels"
	params.EADFilePath = filepath.Join(falesTestFixturePath, "mss_460.xml")
	params.JSONReferenceFilePath = filepath.Join(falesTestFixturePath, "mss_460.json")
	params.JSONErrorFilePath = "./testdata/tmp/failing-test-barcode-removal.json"

	runiJSONComparisonTest(t, &params)
}

func TestUpdateDonors(t *testing.T) {
	t.Run("Update Donors", func(t *testing.T) {
		var sut EAD

		want := []FilteredString(nil)
		got := sut.Donors
		assertFilteredStringSlicesEqual(t, want, got, "Initial ead.Donors")

		donors := []FilteredString{"a", "x", "c", "d"}
		sut.Donors = donors
		want = donors
		got = sut.Donors
		assertFilteredStringSlicesEqual(t, want, got, "Post-update ead.Donors")
	})
}

func TestJSONMarshalingWithDonors(t *testing.T) {
	t.Run("JSON Marshaling with Donors", func(t *testing.T) {
		ead := getOmegaEAD(t)

		ead.Donors = []FilteredString{" a", "x ", " Q ", "d"}
		jsonData, err := json.MarshalIndent(ead, "", "    ")
		failOnError(t, err, "Unexpected error marshaling JSON")

		// reference file includes newline at end of file so
		// add newline to jsonData
		jsonData = append(jsonData, '\n')

		referenceFile := omegaTestFixturePath + "/" + "mos_2021-with-donors.json"
		referenceFileContents, err := os.ReadFile(referenceFile)
		failOnError(t, err, "Unexpected error reading reference file")

		if !bytes.Equal(referenceFileContents, jsonData) {
			jsonFile := "./testdata/tmp/failing-donor-marshal.json"
			err = os.WriteFile(jsonFile, []byte(jsonData), 0644)
			failOnError(t, err, fmt.Sprintf("Unexpected error writing %s", jsonFile))

			errMsg := fmt.Sprintf("JSON Data does not match reference file.\ndiff %s %s", jsonFile, referenceFile)
			t.Errorf(errMsg)
		}
	})
}

func TestJSONMarshalingWithEmptyDAORoles(t *testing.T) {
	t.Run("JSON Marshaling with Empty DAO Roles", func(t *testing.T) {
		ead := getNYHSFoundling(t)

		jsonData, err := json.MarshalIndent(ead, "", "    ")
		failOnError(t, err, "Unexpected error marshaling JSON")

		// reference file includes newline at end of file so
		// add newline to jsonData
		jsonData = append(jsonData, '\n')

		referenceFile := nyhsTestFixturePath + "/" + "nyhs_foundling.json"
		referenceFileContents, err := os.ReadFile(referenceFile)
		failOnError(t, err, "Unexpected error reading reference file")

		if !bytes.Equal(referenceFileContents, jsonData) {
			jsonFile := "./testdata/tmp/failing-empty-role-marshal.json"
			err = os.WriteFile(jsonFile, []byte(jsonData), 0644)
			failOnError(t, err, fmt.Sprintf("Unexpected error writing %s", jsonFile))

			errMsg := fmt.Sprintf("JSON Data does not match reference file.\ndiff %s %s", jsonFile, referenceFile)
			t.Errorf(errMsg)
		}
	})
}

func TestJSONMarshalingWithDonorsAndImageAndImageSets(t *testing.T) {
	t.Run("JSON Marshaling with Donors", func(t *testing.T) {
		ead := getOmegaEAD(t)

		ead.Donors = []FilteredString{" a", "x ", " Q ", "d"}

		testProcessDID(&ead.ArchDesc.DID)
		testProcessCs(ead.ArchDesc.DSC.C)

		jsonData, err := json.MarshalIndent(ead, "", "    ")
		failOnError(t, err, "Unexpected error marshaling JSON")

		// reference file includes newline at end of file so
		// add newline to jsonData
		jsonData = append(jsonData, '\n')

		referenceFile := omegaTestFixturePath + "/" + "mos_2021-with-donors-with-image-counts.json"
		referenceFileContents, err := os.ReadFile(referenceFile)
		failOnError(t, err, "Unexpected error reading reference file")

		if !bytes.Equal(referenceFileContents, jsonData) {
			jsonFile := "./testdata/tmp/failing-donors-with-image-counts-marshal.json"
			err = os.WriteFile(jsonFile, []byte(jsonData), 0644)
			failOnError(t, err, fmt.Sprintf("Unexpected error writing %s", jsonFile))

			errMsg := fmt.Sprintf("JSON Data does not match reference file.\ndiff %s %s", jsonFile, referenceFile)
			t.Errorf(errMsg)
		}
	})
}

func TestGuideTitle(t *testing.T) {
	t.Run("GuideTitle()", func(t *testing.T) {
		sut := getOmegaEAD(t)

		want := "Megan O'Shea's One Resource to Rule Them All"
		got := sut.GuideTitle()
		assertEqual(t, want, got, "TestGuideTitle")
	})
}

func TestTitleProper(t *testing.T) {
	t.Run("TitleProper()", func(t *testing.T) {
		sut := getOmegaEAD(t)

		want := "Guide to Megan O'Shea's \u003cspan class=\"ead-emph ead-emph-italic\"\u003eOne\u003c/span\u003e Resource to \u003cspan class=\"ead-lb\"\u003e\u003c/span\u003e Rule Them All \u003cspan class=\"ead-num\"\u003eMOS.2021\u003c/span\u003e"
		got := sut.TitleProper()
		assertEqual(t, want, got, "TestTitleProper")
	})
}

func TestThemeIDAndRepoIDFunctions(t *testing.T) {
	t.Run("ThemeID()", func(t *testing.T) {
		var want, got string
		sut := getOmegaEAD(t)
		themeid := "cdf80c84-2655-4a01-895d-fbf9a374c1df"
		repoid := "9d396ffa-1b3e-41f0-8bc9-e101a5a828bc"
		sut.PubInfo.SetPubInfo(themeid, repoid)

		want = themeid
		got = sut.ThemeID()
		assertEqual(t, want, got, "TestThemeID()")

		want = repoid
		got = sut.RepoID()
		assertEqual(t, want, got, "TestRepoID()")
	})
}

func TestInitDAOCounts(t *testing.T) {
	t.Run("InitDAOCounts()", func(t *testing.T) {
		sut := getOmegaEAD(t)
		sut.InitDAOCounts()

		assertEqualUint32(t, 3, sut.DAOInfo.AudioCount, "AudioCount")
		assertEqualUint32(t, 2, sut.DAOInfo.VideoCount, "VideoCount")
		assertEqualUint32(t, 4, sut.DAOInfo.ImageCount, "ImageCount")
		assertEqualUint32(t, 2, sut.DAOInfo.ExternalLinkCount, "ExternalLinkCount")
		assertEqualUint32(t, 1, sut.DAOInfo.ElectronicRecordsReadingRoomCount, "ElectronicRecordsReadingRoomCount")
		assertEqualUint32(t, 1, sut.DAOInfo.AudioReadingRoomCount, "AudioReadingRoomCount")
		assertEqualUint32(t, 1, sut.DAOInfo.VideoReadingRoomCount, "VideoReadingRoomCount")
	})
}

func TestEADID(t *testing.T) {
	t.Run("EADID()", func(t *testing.T) {
		sut := getOmegaEAD(t)

		assertEqual(t, "mos_2021", sut.EADID(), "EADID()")
	})
}

func TestDAOCountFunctions(t *testing.T) {
	t.Run("InitDAOCounts()", func(t *testing.T) {
		sut := getOmegaEAD(t)
		sut.InitDAOCounts()

		assertEqualUint32(t, 14, sut.AllDAOCount(), "AllDAOCount")
		assertEqualUint32(t, 3, sut.AudioDAOCount(), "AudioDAOCount")
		assertEqualUint32(t, 2, sut.VideoDAOCount(), "VideoDAOCount")
		assertEqualUint32(t, 4, sut.ImageDAOCount(), "ImageDAOCount")
		assertEqualUint32(t, 2, sut.ExternalLinkDAOCount(), "ExternalLinkDAOCount")
		assertEqualUint32(t, 1, sut.ElectronicRecordsReadingRoomDAOCount(), "ElectronicRecordsReadingRoomDAOCount")
		assertEqualUint32(t, 1, sut.AudioReadingRoomDAOCount(), "AudioReadingRoomDAOCount")
		assertEqualUint32(t, 1, sut.VideoReadingRoomDAOCount(), "VideoReadingRoomDAOCount")

		assertEqualUint32(t, 14, uint32(len(sut.DAOInfo.AllDAOs)), "AllDAOs")
		assertEqualUint32(t, 3, uint32(len(sut.DAOInfo.AudioDAOs)), "AudioDAOs")
		assertEqualUint32(t, 2, uint32(len(sut.DAOInfo.VideoDAOs)), "VideoDAOs")
		assertEqualUint32(t, 4, uint32(len(sut.DAOInfo.ImageDAOs)), "ImageDAOs")
		assertEqualUint32(t, 2, uint32(len(sut.DAOInfo.ExternalLinkDAOs)), "ExternalLinkDAOs")
		assertEqualUint32(t, 1, uint32(len(sut.DAOInfo.ElectronicRecordsReadingRoomDAOs)), "ElectronicRecordsReadingRoomDAOs")
		assertEqualUint32(t, 1, uint32(len(sut.DAOInfo.AudioReadingRoomDAOs)), "AudioReadingRoomDAOs")
		assertEqualUint32(t, 1, uint32(len(sut.DAOInfo.VideoReadingRoomDAOs)), "VideoReadingRoomDAOs")
	})
}

func TestDAOInfoClearFunction(t *testing.T) {
	t.Run("DAOInfo.Clear()", func(t *testing.T) {
		sut := getOmegaEAD(t)
		sut.InitDAOCounts()

		assertEqualUint32(t, 14, sut.AllDAOCount(), "AllDAOCount")
		assertEqualUint32(t, 3, sut.AudioDAOCount(), "AudioDAOCount")
		assertEqualUint32(t, 2, sut.VideoDAOCount(), "VideoDAOCount")
		assertEqualUint32(t, 4, sut.ImageDAOCount(), "ImageDAOCount")
		assertEqualUint32(t, 2, sut.ExternalLinkDAOCount(), "ExternalLinkDAOCount")
		assertEqualUint32(t, 1, sut.ElectronicRecordsReadingRoomDAOCount(), "ElectronicRecordsReadingRoomDAOCount")
		assertEqualUint32(t, 1, sut.AudioReadingRoomDAOCount(), "AudioReadingRoomDAOCount")
		assertEqualUint32(t, 1, sut.VideoReadingRoomDAOCount(), "VideoReadingRoomDAOCount")

		assertEqualUint32(t, 14, uint32(len(sut.DAOInfo.AllDAOs)), "AllDAOs")
		assertEqualUint32(t, 3, uint32(len(sut.DAOInfo.AudioDAOs)), "AudioDAOs")
		assertEqualUint32(t, 2, uint32(len(sut.DAOInfo.VideoDAOs)), "VideoDAOs")
		assertEqualUint32(t, 4, uint32(len(sut.DAOInfo.ImageDAOs)), "ImageDAOs")
		assertEqualUint32(t, 2, uint32(len(sut.DAOInfo.ExternalLinkDAOs)), "ExternalLinkDAOs")
		assertEqualUint32(t, 1, uint32(len(sut.DAOInfo.ElectronicRecordsReadingRoomDAOs)), "ElectronicRecordsReadingRoomDAOs")
		assertEqualUint32(t, 1, uint32(len(sut.DAOInfo.AudioReadingRoomDAOs)), "AudioReadingRoomDAOs")
		assertEqualUint32(t, 1, uint32(len(sut.DAOInfo.VideoReadingRoomDAOs)), "VideoReadingRoomDAOs")

		sut.DAOInfo.Clear()

		assertEqualUint32(t, 0, sut.AllDAOCount(), "AllDAOCount")
		assertEqualUint32(t, 0, sut.AudioDAOCount(), "AudioDAOCount")
		assertEqualUint32(t, 0, sut.VideoDAOCount(), "VideoDAOCount")
		assertEqualUint32(t, 0, sut.ImageDAOCount(), "ImageDAOCount")
		assertEqualUint32(t, 0, sut.ExternalLinkDAOCount(), "ExternalLinkDAOCount")
		assertEqualUint32(t, 0, sut.ElectronicRecordsReadingRoomDAOCount(), "ElectronicRecordsReadingRoomDAOCount")
		assertEqualUint32(t, 0, sut.AudioReadingRoomDAOCount(), "AudioReadingRoomDAOCount")
		assertEqualUint32(t, 0, sut.VideoReadingRoomDAOCount(), "VideoReadingRoomDAOCount")

		assertEqualUint32(t, 0, uint32(len(sut.DAOInfo.AllDAOs)), "AllDAOs")
		assertEqualUint32(t, 0, uint32(len(sut.DAOInfo.AudioDAOs)), "AudioDAOs")
		assertEqualUint32(t, 0, uint32(len(sut.DAOInfo.VideoDAOs)), "VideoDAOs")
		assertEqualUint32(t, 0, uint32(len(sut.DAOInfo.ImageDAOs)), "ImageDAOs")
		assertEqualUint32(t, 0, uint32(len(sut.DAOInfo.ExternalLinkDAOs)), "ExternalLinkDAOs")
		assertEqualUint32(t, 0, uint32(len(sut.DAOInfo.ElectronicRecordsReadingRoomDAOs)), "ElectronicRecordsReadingRoomDAOs")
		assertEqualUint32(t, 0, uint32(len(sut.DAOInfo.AudioReadingRoomDAOs)), "AudioReadingRoomDAOs")
		assertEqualUint32(t, 0, uint32(len(sut.DAOInfo.VideoReadingRoomDAOs)), "VideoReadingRoomDAOs")
	})
}

func TestDAOGrpCountFunction(t *testing.T) {
	t.Run("InitDAOGrpCount()", func(t *testing.T) {
		sut := getOmegaEAD(t)
		sut.InitDAOGrpCount()

		assertEqualUint32(t, 7, sut.AllDAOGrpCount(), "AllDAOGrpCount")
		assertEqualUint32(t, 7, uint32(len(sut.DAOGrpInfo.AllDAOGrps)), "AllDAOGrps length")
	})
}

func TestDAOGrpInfoClearFunction(t *testing.T) {
	t.Run("DAOGrpInfo.Clear()", func(t *testing.T) {
		sut := getOmegaEAD(t)
		sut.InitDAOGrpCount()

		assertEqualUint32(t, 7, sut.AllDAOGrpCount(), "AllDAOGrpCount")
		assertEqualUint32(t, 7, uint32(len(sut.DAOGrpInfo.AllDAOGrps)), "AllDAOGrps length")

		sut.DAOGrpInfo.Clear()

		assertEqualUint32(t, 0, sut.AllDAOGrpCount(), "AllDAOGrpCount")
		assertEqualUint32(t, 0, uint32(len(sut.DAOGrpInfo.AllDAOGrps)), "AllDAOGrps length")
	})
}

func TestDAOParentDID(t *testing.T) {
	t.Run("Test DAO Parent DID", func(t *testing.T) {
		sut := getOmegaEAD(t)
		sut.InitDAOCounts()

		assertEqual(t, "mos_2021_2", string(sut.DAOInfo.AllDAOs[0].ParentDID.UnitID), "Test DAO ParentDID UnitID Level 2")
		assertEqual(t, "mos_2021_3", string(sut.DAOInfo.AllDAOs[1].ParentDID.UnitID), "Test DAO ParentDID UnitID Level 3")
	})
}

func TestFilteredStringStringFunction(t *testing.T) {
	t.Run("Test Filtered String String()", func(t *testing.T) {
		sut := FilteredString("Gilded Youth &#10;| Actors in Image: Nagah Salem")

		assertEqual(t, "Gilded Youth | Actors in Image: Nagah Salem", sut.String(), "Test FilteredString.String()")
	})
}

func TestGetConvertedTextWithTags(t *testing.T) {
	t.Run("Test GetConvertedTextWithTags()", func(t *testing.T) {
		sut, _ := GetConvertedTextWithTags(`The Young Devil
| Actors in Image: Ahmad Ramzy, Amal Farid, Hussein Riad, Zinat Sedki, Ragaa Youssef: 1958-`)

		assertEqual(t, "The Young Devil | Actors in Image: Ahmad Ramzy, Amal Farid, Hussein Riad, Zinat Sedki, Ragaa Youssef: 1958-", string(sut), "Test TestGetConvertedTextWithTags()")
	})
}

func TestGetConvertedTextWithTagsNoLBConversion(t *testing.T) {
	input := `Some materials may be restricted. Permission to publish materials must be obtained in writing from the:<lb/>
New York University Archives<lb/> Elmer Holmes Bobst Library<lb/> 70 Washington Square South<lb/> New York, NY 10012<lb/> Phone: (212) 998-2641<lb/>Fax: (212) 995-4225<lb/>E-mail: university-archives@nyu.edu<lb/>`
	want := `Some materials may be restricted. Permission to publish materials must be obtained in writing from the:<span class="ead-lb"></span> New York University Archives<span class="ead-lb"></span> Elmer Holmes Bobst Library<span class="ead-lb"></span> 70 Washington Square South<span class="ead-lb"></span> New York, NY 10012<span class="ead-lb"></span> Phone: (212) 998-2641<span class="ead-lb"></span>Fax: (212) 995-4225<span class="ead-lb"></span>E-mail: university-archives@nyu.edu<span class="ead-lb"></span>`

	t.Run("Test GetConvertedTextWithTagsNoLBConversion()", func(t *testing.T) {
		sut, _ := GetConvertedTextWithTagsNoLBConversion(input)
		assertEqual(t, want, string(sut), "Test TestGetConvertedTextWithTagsNoLBConversion()")
	})
}

func TestJSONMarshalingInitPresentationContainersNOOP(t *testing.T) {
	t.Run("JSON Marshaling with call to InitPresentationContainers() NOOP", func(t *testing.T) {
		ead := getOmegaEAD(t)

		ead.InitPresentationContainers()

		jsonData, err := json.MarshalIndent(ead, "", "    ")
		failOnError(t, err, "Unexpected error marshaling JSON")

		// reference file includes newline at end of file so
		// add newline to jsonData
		jsonData = append(jsonData, '\n')

		referenceFile := omegaTestFixturePath + "/" + "mos_2021.json"
		referenceFileContents, err := os.ReadFile(referenceFile)
		failOnError(t, err, "Unexpected error reading reference file")

		if !bytes.Equal(referenceFileContents, jsonData) {
			jsonFile := "./testdata/tmp/failing-marshal-with-presentation-containers-noop.json"
			err = os.WriteFile(jsonFile, []byte(jsonData), 0644)
			failOnError(t, err, fmt.Sprintf("Unexpected error writing %s", jsonFile))

			errMsg := fmt.Sprintf("JSON Data does not match reference file.\ndiff %s %s", jsonFile, referenceFile)
			t.Errorf(errMsg)
		}
	})
}

func TestJSONMarshalingInitPresentationContainers(t *testing.T) {
	var params iJSONTestParams

	params.TestName = "JSON Marshaling with call to InitPresentationContainers()"
	params.EADFilePath = filepath.Join(akkasahTestFixturePath, "ad_mc_030_ref184.xml")
	params.JSONReferenceFilePath = filepath.Join(akkasahTestFixturePath, "ad_mc_030_ref184.json")
	params.JSONErrorFilePath = "./testdata/tmp/failing-marshal-with-presentation-containers.json"

	ead := getTestEAD(t, params.EADFilePath)
	ead.InitPresentationContainers()

	params.PrePopulatedEAD = ead
	runiJSONComparisonTest(t, &params)
}

func TestInitPresentationContainersC(t *testing.T) {
	t.Run("InitPresentationContainers() Collapse All Containers", func(t *testing.T) {
		ead := getPresentationContainerEAD(t, "pc-c.xml")

		assertEqual(t, "file-001", string(ead.ArchDesc.DSC.C[0].ID), "initial container ID")
		assertEqual(t, "file-002", string(ead.ArchDesc.DSC.C[1].ID), "initial container ID")
		assertEqual(t, "file-003", string(ead.ArchDesc.DSC.C[2].ID), "initial container ID")
		assertEqual(t, "file-004", string(ead.ArchDesc.DSC.C[3].ID), "initial container ID")
		assertEqual(t, "file-005", string(ead.ArchDesc.DSC.C[4].ID), "initial container ID")
		assertEqual(t, "file-006", string(ead.ArchDesc.DSC.C[5].ID), "initial container ID")

		ead.InitPresentationContainers()

		assertEqual(t, "items001", string(ead.ArchDesc.DSC.C[0].ID), "presentation container ID")
		assertEqual(t, "file-001", string(ead.ArchDesc.DSC.C[0].C[0].ID), "collapsed container ID")
		assertEqual(t, "file-002", string(ead.ArchDesc.DSC.C[0].C[1].ID), "collapsed container ID")
		assertEqual(t, "file-003", string(ead.ArchDesc.DSC.C[0].C[2].ID), "collapsed container ID")
		assertEqual(t, "file-004", string(ead.ArchDesc.DSC.C[0].C[3].ID), "collapsed container ID")
		assertEqual(t, "file-005", string(ead.ArchDesc.DSC.C[0].C[4].ID), "collapsed container ID")
		assertEqual(t, "file-006", string(ead.ArchDesc.DSC.C[0].C[5].ID), "collapsed container ID")

		assertEqual(t, "View Items", string(ead.ArchDesc.DSC.C[0].DID.UnitTitle.Value), "presentation container UnitTitle")
		assertEqual(t, "dl-presentation", string(ead.ArchDesc.DSC.C[0].Level), "presentation container Level")
	})
}

func TestInitPresentationContainersCK(t *testing.T) {
	t.Run("InitPresentationContainers() Collapse First Containers", func(t *testing.T) {
		ead := getPresentationContainerEAD(t, "pc-c-k.xml")

		assertEqual(t, "file-001", string(ead.ArchDesc.DSC.C[0].ID), "initial container ID")
		assertEqual(t, "file-002", string(ead.ArchDesc.DSC.C[1].ID), "initial container ID")
		assertEqual(t, "file-003", string(ead.ArchDesc.DSC.C[2].ID), "initial container ID")
		assertEqual(t, "series-001", string(ead.ArchDesc.DSC.C[3].ID), "initial container ID")
		assertEqual(t, "otherlevel-001", string(ead.ArchDesc.DSC.C[4].ID), "initial container ID")
		assertEqual(t, "recordgrp-001", string(ead.ArchDesc.DSC.C[5].ID), "initial container ID")

		ead.InitPresentationContainers()

		assertEqual(t, "items001", string(ead.ArchDesc.DSC.C[0].ID), "presentation container ID")
		assertEqual(t, "file-001", string(ead.ArchDesc.DSC.C[0].C[0].ID), "collapsed container ID")
		assertEqual(t, "file-002", string(ead.ArchDesc.DSC.C[0].C[1].ID), "collapsed container ID")
		assertEqual(t, "file-003", string(ead.ArchDesc.DSC.C[0].C[2].ID), "collapsed container ID")
		assertEqual(t, "series-001", string(ead.ArchDesc.DSC.C[1].ID), "kept container ID")
		assertEqual(t, "otherlevel-001", string(ead.ArchDesc.DSC.C[2].ID), "kept container ID")
		assertEqual(t, "recordgrp-001", string(ead.ArchDesc.DSC.C[3].ID), "kept container ID")

		assertEqual(t, "View Items", string(ead.ArchDesc.DSC.C[0].DID.UnitTitle.Value), "presentation container UnitTitle")
		assertEqual(t, "dl-presentation", string(ead.ArchDesc.DSC.C[0].Level), "presentation container Level")
	})
}

func TestInitPresentationContainersKC(t *testing.T) {
	t.Run("InitPresentationContainers() Collapse Last Containers", func(t *testing.T) {
		ead := getPresentationContainerEAD(t, "pc-k-c.xml")

		assertEqual(t, "series-001", string(ead.ArchDesc.DSC.C[0].ID), "initial container ID")
		assertEqual(t, "otherlevel-001", string(ead.ArchDesc.DSC.C[1].ID), "initial container ID")
		assertEqual(t, "recordgrp-001", string(ead.ArchDesc.DSC.C[2].ID), "initial container ID")
		assertEqual(t, "file-001", string(ead.ArchDesc.DSC.C[3].ID), "initial container ID")
		assertEqual(t, "file-002", string(ead.ArchDesc.DSC.C[4].ID), "initial container ID")
		assertEqual(t, "file-003", string(ead.ArchDesc.DSC.C[5].ID), "initial container ID")

		ead.InitPresentationContainers()

		assertEqual(t, "series-001", string(ead.ArchDesc.DSC.C[0].ID), "collapsed container ID")
		assertEqual(t, "otherlevel-001", string(ead.ArchDesc.DSC.C[1].ID), "collapsed container ID")
		assertEqual(t, "recordgrp-001", string(ead.ArchDesc.DSC.C[2].ID), "collapsed container ID")
		assertEqual(t, "items001", string(ead.ArchDesc.DSC.C[3].ID), "presentation container ID")
		assertEqual(t, "file-001", string(ead.ArchDesc.DSC.C[3].C[0].ID), "collapsed container ID")
		assertEqual(t, "file-002", string(ead.ArchDesc.DSC.C[3].C[1].ID), "collapsed container ID")
		assertEqual(t, "file-003", string(ead.ArchDesc.DSC.C[3].C[2].ID), "collapsed container ID")

		assertEqual(t, "View Items", string(ead.ArchDesc.DSC.C[3].DID.UnitTitle.Value), "presentation container UnitTitle")
		assertEqual(t, "dl-presentation", string(ead.ArchDesc.DSC.C[3].Level), "presentation container Level")
	})
}

func TestInitPresentationContainersCKC(t *testing.T) {
	t.Run("InitPresentationContainers() Collapse First and Last Containers", func(t *testing.T) {
		ead := getPresentationContainerEAD(t, "pc-c-k-c.xml")

		assertEqual(t, "file-001", string(ead.ArchDesc.DSC.C[0].ID), "initial container ID")
		assertEqual(t, "file-002", string(ead.ArchDesc.DSC.C[1].ID), "initial container ID")
		assertEqual(t, "file-003", string(ead.ArchDesc.DSC.C[2].ID), "initial container ID")
		assertEqual(t, "series-001", string(ead.ArchDesc.DSC.C[3].ID), "initial container ID")
		assertEqual(t, "otherlevel-001", string(ead.ArchDesc.DSC.C[4].ID), "initial container ID")
		assertEqual(t, "recordgrp-001", string(ead.ArchDesc.DSC.C[5].ID), "initial container ID")
		assertEqual(t, "file-004", string(ead.ArchDesc.DSC.C[6].ID), "initial container ID")
		assertEqual(t, "file-005", string(ead.ArchDesc.DSC.C[7].ID), "initial container ID")
		assertEqual(t, "file-006", string(ead.ArchDesc.DSC.C[8].ID), "initial container ID")

		ead.InitPresentationContainers()

		assertEqual(t, "items001", string(ead.ArchDesc.DSC.C[0].ID), "presentation container ID")
		assertEqual(t, "file-001", string(ead.ArchDesc.DSC.C[0].C[0].ID), "collapsed container ID")
		assertEqual(t, "file-002", string(ead.ArchDesc.DSC.C[0].C[1].ID), "collapsed container ID")
		assertEqual(t, "file-003", string(ead.ArchDesc.DSC.C[0].C[2].ID), "collapsed container ID")
		assertEqual(t, "series-001", string(ead.ArchDesc.DSC.C[1].ID), "collapsed container ID")
		assertEqual(t, "otherlevel-001", string(ead.ArchDesc.DSC.C[2].ID), "collapsed container ID")
		assertEqual(t, "recordgrp-001", string(ead.ArchDesc.DSC.C[3].ID), "collapsed container ID")
		assertEqual(t, "items002", string(ead.ArchDesc.DSC.C[4].ID), "presentation container ID")
		assertEqual(t, "file-004", string(ead.ArchDesc.DSC.C[4].C[0].ID), "collapsed container ID")
		assertEqual(t, "file-005", string(ead.ArchDesc.DSC.C[4].C[1].ID), "collapsed container ID")
		assertEqual(t, "file-006", string(ead.ArchDesc.DSC.C[4].C[2].ID), "collapsed container ID")

		assertEqual(t, "View Items", string(ead.ArchDesc.DSC.C[0].DID.UnitTitle.Value), "presentation container UnitTitle")
		assertEqual(t, "dl-presentation", string(ead.ArchDesc.DSC.C[0].Level), "presentation container Level")

		assertEqual(t, "View Items", string(ead.ArchDesc.DSC.C[4].DID.UnitTitle.Value), "presentation container UnitTitle")
		assertEqual(t, "dl-presentation", string(ead.ArchDesc.DSC.C[4].Level), "presentation container Level")
	})
}

func TestInitPresentationContainersKCK(t *testing.T) {
	t.Run("InitPresentationContainers() Collapse Middle Containers", func(t *testing.T) {
		ead := getPresentationContainerEAD(t, "pc-k-c-k.xml")

		assertEqual(t, "series-001", string(ead.ArchDesc.DSC.C[0].ID), "initial container ID")
		assertEqual(t, "otherlevel-001", string(ead.ArchDesc.DSC.C[1].ID), "initial container ID")
		assertEqual(t, "recordgrp-001", string(ead.ArchDesc.DSC.C[2].ID), "initial container ID")
		assertEqual(t, "file-001", string(ead.ArchDesc.DSC.C[3].ID), "initial container ID")
		assertEqual(t, "file-002", string(ead.ArchDesc.DSC.C[4].ID), "initial container ID")
		assertEqual(t, "file-003", string(ead.ArchDesc.DSC.C[5].ID), "initial container ID")
		assertEqual(t, "series-002", string(ead.ArchDesc.DSC.C[6].ID), "initial container ID")
		assertEqual(t, "otherlevel-002", string(ead.ArchDesc.DSC.C[7].ID), "initial container ID")
		assertEqual(t, "recordgrp-002", string(ead.ArchDesc.DSC.C[8].ID), "initial container ID")

		ead.InitPresentationContainers()

		assertEqual(t, "series-001", string(ead.ArchDesc.DSC.C[0].ID), "kept container ID")
		assertEqual(t, "otherlevel-001", string(ead.ArchDesc.DSC.C[1].ID), "kept container ID")
		assertEqual(t, "recordgrp-001", string(ead.ArchDesc.DSC.C[2].ID), "kept container ID")
		assertEqual(t, "items001", string(ead.ArchDesc.DSC.C[3].ID), "presentation container ID")
		assertEqual(t, "file-001", string(ead.ArchDesc.DSC.C[3].C[0].ID), "collapsed container ID")
		assertEqual(t, "file-002", string(ead.ArchDesc.DSC.C[3].C[1].ID), "collapsed container ID")
		assertEqual(t, "file-003", string(ead.ArchDesc.DSC.C[3].C[2].ID), "collapsed container ID")
		assertEqual(t, "series-002", string(ead.ArchDesc.DSC.C[4].ID), "kept container ID")
		assertEqual(t, "otherlevel-002", string(ead.ArchDesc.DSC.C[5].ID), "kept container ID")
		assertEqual(t, "recordgrp-002", string(ead.ArchDesc.DSC.C[6].ID), "kept container ID")

		assertEqual(t, "View Items", string(ead.ArchDesc.DSC.C[3].DID.UnitTitle.Value), "presentation container UnitTitle")
		assertEqual(t, "dl-presentation", string(ead.ArchDesc.DSC.C[3].Level), "presentation container Level")
	})
}

func TestInitPresentationContainersNoContainers(t *testing.T) {
	t.Run("InitPresentationContainers() Collapse All Containers", func(t *testing.T) {
		ead := getPresentationContainerEAD(t, "pc-no-containers.xml")

		if nil != ead.ArchDesc.DSC.C {
			t.Errorf("expected initial container list to be empty")
		}

		ead.InitPresentationContainers()

		if nil != ead.ArchDesc.DSC.C {
			t.Errorf("expected container list to still be empty after InitPresentationContainers()")
		}
	})
}

func TestJSONMarshalingWithPresentationElementsInTitleStmtChildren(t *testing.T) {
	var params iJSONTestParams

	params.TestName = "JSON Marshaling with Presentation Element In TitleStmt children"
	params.EADFilePath = filepath.Join(omegaTestFixturePath, "mos_2021-with-presentation-elements-in-titlestmt-children.xml")
	params.JSONReferenceFilePath = filepath.Join(omegaTestFixturePath, "mos_2021-with-presentation-elements-in-titlestmt-children.json")
	params.JSONErrorFilePath = "./testdata/tmp/failing-with-presentation-elements-in-titlestmt-children.json"

	runiJSONComparisonTest(t, &params)
}

func runiJSONComparisonTest(t *testing.T, params *iJSONTestParams) {

	var ead *EAD
	t.Run(params.TestName, func(t *testing.T) {
		if params.PrePopulatedEAD == nil {
			ead = getTestEAD(t, params.EADFilePath)
		} else {
			ead = params.PrePopulatedEAD
		}
		jsonData, err := json.MarshalIndent(ead, "", "    ")
		failOnError(t, err, "Unexpected error marshaling JSON")

		// reference file includes newline at end of file so
		// add newline to jsonData
		jsonData = append(jsonData, '\n')

		referenceFile := params.JSONReferenceFilePath
		referenceFileContents, err := os.ReadFile(referenceFile)
		failOnError(t, err, "Unexpected error reading reference file")

		if !bytes.Equal(referenceFileContents, jsonData) {
			jsonErrorFile := params.JSONErrorFilePath
			err = os.WriteFile(jsonErrorFile, []byte(jsonData), 0644)
			failOnError(t, err, fmt.Sprintf("Unexpected error writing %s", jsonErrorFile))

			errMsg := fmt.Sprintf("JSON Data does not match reference file.\ndiff %s %s", jsonErrorFile, referenceFile)
			t.Errorf(errMsg)
		}
	})
}
