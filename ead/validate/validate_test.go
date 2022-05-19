package validate

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

var fixturesDirPath string
var invalidEadDataFixturePath string
var invalidXMLFixturePath string
var missingRequiredElementsFixturePath string
var validEADFixturePath string

// Source: https://intellij-support.jetbrains.com/hc/en-us/community/posts/360009685279-Go-test-working-directory-keeps-changing-to-dir-of-the-test-file-instead-of-value-in-template?page=1#community_comment_360002183640
func init() {
	_, filename, _, _ := runtime.Caller(0)
	root := path.Join(path.Dir(filename), "..")
	err := os.Chdir(root)
	if err != nil {
		panic(err)
	}

	fixturesDirPath = filepath.Join(root, "validate", "testdata", "fixtures")
	invalidEadDataFixturePath = filepath.Join(fixturesDirPath, "mc_100-invalid-eadid-repository-role-relator-codes-unpublished-material.xml")
	invalidXMLFixturePath = filepath.Join(fixturesDirPath, "invalid-xml.xml")
	missingRequiredElementsFixturePath = filepath.Join(fixturesDirPath, "mc_100-missing-eadid-and-repository-corpname.xml")
	validEADFixturePath = filepath.Join(fixturesDirPath, "mc_100.xml")
}

func TestValidateEADInvalidXML(t *testing.T) {
	const expected = "The XML in this file is not valid.  Please check it using an XML validator."

	var errors = ValidateEAD(getEADXML(invalidXMLFixturePath))
	var numErrors = len(errors)
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", numErrors)
	}

	if errors[0] != "The XML in this file is not valid.  Please check it using an XML validator." {
		t.Errorf(`Expected error \"%s\", got "%s\"`, expected, errors[0])
	}
}

func getEADXML(filepath string) []byte {
	EADXML, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	return EADXML
}
