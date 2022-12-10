package ead

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

var testFixturePath string = filepath.Join(".", "testdata")
var omegaTestFixturePath string = filepath.Join(testFixturePath, "omega", "v0.1.5")
var falesTestFixturePath string = filepath.Join(testFixturePath, "fales")
var nyhsTestFixturePath string = filepath.Join(testFixturePath, "nyhs")

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
	EADXML, err := ioutil.ReadFile(omegaTestFixturePath + "/" + "Omega-EAD.xml")
	failOnError(t, err, "Unexpected error")

	var ead EAD
	err = xml.Unmarshal([]byte(EADXML), &ead)
	failOnError(t, err, "Unexpected error")

	return ead
}

func getFalesMSS460EAD(t *testing.T) EAD {
	EADXML, err := ioutil.ReadFile(falesTestFixturePath + "/" + "/mss_460.xml")
	failOnError(t, err, "Unexpected error")

	var ead EAD
	err = xml.Unmarshal([]byte(EADXML), &ead)
	failOnError(t, err, "Unexpected error")

	return ead
}

func getNYHSFoundling(t *testing.T) EAD {
	EADXML, err := ioutil.ReadFile(nyhsTestFixturePath + "/" + "nyhs_foundling.xml")
	failOnError(t, err, "Unexpected error")

	var ead EAD
	err = xml.Unmarshal([]byte(EADXML), &ead)
	failOnError(t, err, "Unexpected error")

	return ead
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
	t.Run("JSON Marshaling", func(t *testing.T) {
		ead := getOmegaEAD(t)

		jsonData, err := json.MarshalIndent(ead, "", "    ")
		failOnError(t, err, "Unexpected error marshaling JSON")

		// reference file includes newline at end of file so
		// add newline to jsonData
		jsonData = append(jsonData, '\n')

		referenceFile := omegaTestFixturePath + "/" + "mos_2021.json"
		referenceFileContents, err := ioutil.ReadFile(referenceFile)
		failOnError(t, err, "Unexpected error reading reference file")

		if !bytes.Equal(referenceFileContents, jsonData) {
			jsonFile := "./testdata/tmp/failing-marshal.json"
			err = ioutil.WriteFile(jsonFile, []byte(jsonData), 0644)
			failOnError(t, err, fmt.Sprintf("Unexpected error writing %s", jsonFile))

			errMsg := fmt.Sprintf("JSON Data does not match reference file.\ndiff %s %s", jsonFile, referenceFile)
			t.Errorf(errMsg)
		}
	})
}

func TestUpdateRunInfo(t *testing.T) {
	t.Run("Update RunInfo", func(t *testing.T) {
		var sut EAD

		want := ""
		got := sut.RunInfo.PkgVersion
		assertEqual(t, want, got, "Initial ead.RunInfo.PkgVersion")

		want = "0001-01-01T00:00:00Z"
		got = sut.RunInfo.TimeStamp.Format(time.RFC3339)
		assertEqual(t, want, got, "Initial ead.RunInfo.TimeStamp")

		now := time.Now()
		version := Version // from ead package constant
		sourceFile := "/a/very/fine/path/to/an/ead.xml"

		sut.RunInfo.SetRunInfo(version, now, sourceFile)

		want = version
		got = sut.RunInfo.PkgVersion
		assertEqual(t, want, got, "Post-assignment ead.RunInfo.PkgVersion")

		want = now.Format(time.RFC3339)
		got = sut.RunInfo.TimeStamp.Format(time.RFC3339)
		assertEqual(t, want, got, "Post-assignment ead.RunInfo.TimeStamp")

		want = sourceFile
		got = sut.RunInfo.SourceFile
		assertEqual(t, want, got, "set ead.RunInfo.SourceFile")
	})
}

func TestUpdatePubInfo(t *testing.T) {
	t.Run("Update PubInfo", func(t *testing.T) {
		var sut EAD

		want := ""
		got := sut.PubInfo.ThemeID
		assertEqual(t, want, got, "Initial ead.PubInfo.ThemeID")

		themeid := "cdf80c84-2655-4a01-895d-fbf9a374c1df"
		sut.PubInfo.SetPubInfo(themeid)

		want = themeid
		got = sut.PubInfo.ThemeID
		assertEqual(t, want, got, "Post-assignment ead.PubInfo.ThemeID")

	})
}

func TestBarcodeRemovalFromLabels(t *testing.T) {
	t.Run("Barcode Removal from Labels", func(t *testing.T) {
		ead := getFalesMSS460EAD(t)

		jsonData, err := json.MarshalIndent(ead, "", "    ")
		failOnError(t, err, "Unexpected error marshaling JSON")

		// reference file includes newline at end of file so
		// add newline to jsonData
		jsonData = append(jsonData, '\n')

		referenceFile := falesTestFixturePath + "/mss_460.json"
		referenceFileContents, err := ioutil.ReadFile(referenceFile)
		failOnError(t, err, "Unexpected error reading reference file")

		if !bytes.Equal(referenceFileContents, jsonData) {
			jsonFile := "./testdata/tmp/failing-test-barcode-removal.json"
			err = ioutil.WriteFile(jsonFile, []byte(jsonData), 0644)
			failOnError(t, err, fmt.Sprintf("Unexpected error writing %s", jsonFile))

			errMsg := fmt.Sprintf("JSON Data does not match reference file.\ndiff %s %s", jsonFile, referenceFile)
			t.Errorf(errMsg)
		}
	})
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
		referenceFileContents, err := ioutil.ReadFile(referenceFile)
		failOnError(t, err, "Unexpected error reading reference file")

		if !bytes.Equal(referenceFileContents, jsonData) {
			jsonFile := "./testdata/tmp/failing-donor-marshal.json"
			err = ioutil.WriteFile(jsonFile, []byte(jsonData), 0644)
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
		referenceFileContents, err := ioutil.ReadFile(referenceFile)
		failOnError(t, err, "Unexpected error reading reference file")

		if !bytes.Equal(referenceFileContents, jsonData) {
			jsonFile := "./testdata/tmp/failing-empty-role-marshal.json"
			err = ioutil.WriteFile(jsonFile, []byte(jsonData), 0644)
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
		referenceFileContents, err := ioutil.ReadFile(referenceFile)
		failOnError(t, err, "Unexpected error reading reference file")

		if !bytes.Equal(referenceFileContents, jsonData) {
			jsonFile := "./testdata/tmp/failing-donors-with-image-counts-marshal.json"
			err = ioutil.WriteFile(jsonFile, []byte(jsonData), 0644)
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

func TestThemeID(t *testing.T) {
	t.Run("ThemeID()", func(t *testing.T) {
		sut := getOmegaEAD(t)
		themeid := "cdf80c84-2655-4a01-895d-fbf9a374c1df"
		sut.PubInfo.SetPubInfo(themeid)

		want := themeid
		got := sut.ThemeID()
		assertEqual(t, want, got, "TestThemeID")
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
