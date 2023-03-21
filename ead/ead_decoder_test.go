package ead

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var testFixtureDataPath = filepath.Join("testdata", "xmlorder")

func TestNoteChildOrder(t *testing.T) {
	var params iJSONTestParams

	params.TestName = "Test Marshal JSON"
	params.EADFilePath = filepath.Join(testFixtureDataPath, "Omega-EAD.xml")
	params.JSONReferenceFilePath = filepath.Join(testFixtureDataPath, "omega-ead-test-order.json")
	params.JSONErrorFilePath = filepath.Join(testFixtureDataPath, "tmp", "failing-test-marshal.json")

	runiJSONComparisonTest(t, &params)
}

func TestMinimalStreamParsingEmbeddedPresentationElements(t *testing.T) {

	t.Run("Test Minimal Stream-Parsing and JSON Marshaling of Embedded Presentation Elements", func(t *testing.T) {
		xmlSnippet, err := os.ReadFile(filepath.Join(testFixtureDataPath, "materialspec.xml"))
		if err != nil {
			t.Error(err)
		}
		var child EADChild
		err = xml.Unmarshal([]byte(xmlSnippet), &child)
		failOnError(t, err, "Unexpected error unmarshaling XML")

		jsonData, err := json.MarshalIndent(child, "", "    ")
		failOnError(t, err, "Unexpected error marshaling JSON")

		// reference file includes newline at end of file so
		// add newline to jsonData
		jsonData = append(jsonData, '\n')

		referenceFile := filepath.Join(testFixtureDataPath, "minimal-stream-parsing-with-embedded-presentation-elements.json")
		referenceFileContents, err := os.ReadFile(referenceFile)
		failOnError(t, err, "Unexpected error reading reference file")

		if !bytes.Equal(referenceFileContents, jsonData) {
			jsonErrorFile := filepath.Join(testFixtureDataPath, "tmp", "failing-minimal-stream-parsing-with-embedded-presentation-elements.json")
			err = os.WriteFile(jsonErrorFile, []byte(jsonData), 0644)
			failOnError(t, err, fmt.Sprintf("Unexpected error writing %s", jsonErrorFile))

			errMsg := fmt.Sprintf("JSON Data does not match reference file.\ndiff %s %s", jsonErrorFile, referenceFile)
			t.Errorf(errMsg)
		}
	})
}

func TestStreamParsingEmbeddedPresentationElements(t *testing.T) {
	var params iJSONTestParams

	params.TestName = "Test Stream-Parsing and JSON Marshaling Embedded Presentation Elements"
	params.EADFilePath = filepath.Join(testFixtureDataPath, "mss_293.xml")
	params.JSONReferenceFilePath = filepath.Join(testFixtureDataPath, "embedded-presentation-elements.json")
	params.JSONErrorFilePath = filepath.Join(testFixtureDataPath, "tmp", "failing-embedded-presentation-elements.json")

	runiJSONComparisonTest(t, &params)
}
