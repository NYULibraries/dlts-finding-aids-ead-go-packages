package ead

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
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
		if err := ioutil.WriteFile(filepath.Join(testFixtureDataPath, "omega-ead-test-order.json"), jsonData, 0755); err != nil {
			t.Error(err)
		}
	})
}

func getOrderXMLOmega(t *testing.T) EAD {
	EADXML, err := ioutil.ReadFile(filepath.Join(testFixtureDataPath, "Omega-EAD.xml"))
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
