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
