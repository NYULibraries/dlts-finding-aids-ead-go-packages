package ead

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func assert(t *testing.T, want string, got string, label string) {
	if want != got {
		t.Errorf("%s Mismatch: want: %s, got: %s", label, want, got)
	}
}

func TestXMLParsing(t *testing.T) {
	t.Run("XML Parsing", func(t *testing.T) {
		EADXML, err := ioutil.ReadFile("./testdata/v0.0.0/Omega-EAD.xml")
		var ead EAD
		err = xml.Unmarshal([]byte(EADXML), &ead)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		want := "collection"
		got := ead.ArchDesc.Level
		assert(t, want, got, "ArchDesc.Level")
	})
}

func TestJSONMarshaling(t *testing.T) {
	t.Run("JSON Marshaling", func(t *testing.T) {
		EADXML, err := ioutil.ReadFile("./testdata/v0.0.0/Omega-EAD.xml")
		var ead EAD
		err = xml.Unmarshal([]byte(EADXML), &ead)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		jsonData, err := json.MarshalIndent(ead, "", "    ")
		if err != nil {
			t.Errorf("Unexpected error marshaling JSON: %s", err)
		}

		// reference file includes newline at end of file so
		// add newline to jsonData
		jsonData = append(jsonData, '\n')

		referenceFileContents, err := ioutil.ReadFile("./testdata/v0.0.0/mos_2021.json")
		if err != nil {
			t.Errorf("Unexpected error reading reference file: %s", err)
		}

		if !bytes.Equal(referenceFileContents, jsonData) {
			jsonFile := "./testdata/tmp/failing-marshal.json"
			errMsg := fmt.Sprintf("JSON Data does not match reference file. Writing marshaled JSON to: %s", jsonFile)
			t.Errorf(errMsg)
			_ = ioutil.WriteFile(jsonFile, []byte(jsonData), 0644)
		}
	})
}

func TestUpdateRunInfo(t *testing.T) {
	t.Run("JSON Marshaling", func(t *testing.T) {
		var sut EAD

		want := ""
		got := sut.RunInfo.PkgVersion
		assert(t, want, got, "Initial ead.RunInfo.PkgVersion")

		want = "0001-01-01T00:00:00Z"
		got = sut.RunInfo.TimeStamp.Format(time.RFC3339)
		assert(t, want, got, "Initial ead.RunInfo.TimeStamp")

		now := time.Now()
		version := Version // from ead package constant

		sut.RunInfo.SetRunInfo(version, now)

		want = version
		got = sut.RunInfo.PkgVersion
		assert(t, want, got, "Post-assignment ead.RunInfo.PkgVersion")

		want = now.Format(time.RFC3339)
		got = sut.RunInfo.TimeStamp.Format(time.RFC3339)
		assert(t, want, got, "Initial ead.RunInfo.TimeStamp")

	})
}
