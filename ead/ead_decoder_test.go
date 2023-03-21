package ead

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"
)

var testFixtureDataPath = filepath.Join("testdata", "xmlorder")

func TestNoteChildOrder(t *testing.T) {

	t.Run("Test Marshal JSON", func(t *testing.T) {
		ead := getOrderXMLOmega(t)
		jsonData, err := json.MarshalIndent(ead, "", "    ")
		if err != nil {
			t.Error(err)
		}

		// reference file includes newline at end of file so
		// add newline to jsonData
		jsonData = append(jsonData, '\n')
		if err := os.WriteFile(filepath.Join(testFixtureDataPath, "omega-ead-test-order.json"), jsonData, 0755); err != nil {
			t.Error(err)
		}
	})
}

func TestMinimalStreamParsingEmbeddedPresentationElements(t *testing.T) {

	t.Run("Test Minimal Stream-Parsing and JSON Marshaling of Embedded Presentation Elements", func(t *testing.T) {
		xmlSnippet, err := os.ReadFile(filepath.Join(testFixtureDataPath, "materialspec.xml"))
		if err != nil {
			t.Error(err)
		}
		var child EADChild
		err = xml.Unmarshal([]byte(xmlSnippet), &child)
		if err != nil {
			t.Error(err)
		}

		jsonData, err := json.MarshalIndent(child, "", "    ")
		if err != nil {
			t.Error(err)
		}

		// reference file includes newline at end of file so
		// add newline to jsonData
		jsonData = append(jsonData, '\n')
		if err := os.WriteFile(filepath.Join(testFixtureDataPath, "minimal-stream-parsing-with-embedded-presentation-elements.json"), jsonData, 0755); err != nil {
			t.Error(err)
		}
	})
}

func TestStreamParsingEmbeddedPresentationElements(t *testing.T) {
	t.Run("Test Stream-Parsing and JSON Marshaling Embedded Presentation Elements", func(t *testing.T) {
		ead := getEmbeddedPresentationElementsEAD(t)
		//ead := getOrderXMLOmega(t)
		jsonData, err := json.MarshalIndent(ead, "", "    ")
		if err != nil {
			t.Error(err)
		}

		// reference file includes newline at end of file so
		// add newline to jsonData
		jsonData = append(jsonData, '\n')
		if err := os.WriteFile(filepath.Join(testFixtureDataPath, "embedded-presentation-elements.json"), jsonData, 0755); err != nil {
			t.Error(err)
		}
	})
}

func getOrderXMLOmega(t *testing.T) EAD {
	EADXML, err := os.ReadFile(filepath.Join(testFixtureDataPath, "Omega-EAD.xml"))
	if err != nil {
		t.Error(err)
	}

	var ead EAD
	err = xml.Unmarshal([]byte(EADXML), &ead)
	if err != nil {
		t.Error(err)
	}
	return ead
}

func getEmbeddedPresentationElementsEAD(t *testing.T) EAD {
	EADXML, err := os.ReadFile(filepath.Join(testFixtureDataPath, "mss_293.xml"))
	if err != nil {
		t.Error(err)
	}

	var ead EAD
	err = xml.Unmarshal([]byte(EADXML), &ead)
	if err != nil {
		t.Error(err)
	}
	return ead
}
