package validate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var fixturesDirPath string
var invalidEadDataFixturePath string
var invalidXMLFixturePath string
var missingRequiredElementsEADIDAndArchDescFixturePath string
var missingRequiredElementsEADIDAndRepositoryFixturePath string
var missingRequiredElementsEADIDAndRepositoryCorpnameFixturePath string
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
	missingRequiredElementsEADIDAndArchDescFixturePath = filepath.Join(fixturesDirPath, "mc_100-missing-eadid-and-archdesc.xml")
	missingRequiredElementsEADIDAndRepositoryFixturePath = filepath.Join(fixturesDirPath, "mc_100-missing-eadid-and-repository.xml")
	missingRequiredElementsEADIDAndRepositoryCorpnameFixturePath = filepath.Join(fixturesDirPath, "mc_100-missing-eadid-and-repository-corpname.xml")
	validEADFixturePath = filepath.Join(fixturesDirPath, "mc_100.xml")
}

func TestValidateEADInvalidData(t *testing.T) {
	var expected = []string{
		makeInvalidRepositoryErrorMessage("NYU Archives"),
		makeInvalidEADIDErrorMessage("mc.100", []rune{'.'}),
		makeAudienceInternalErrorMessage([]string{"<bioghist>", "<processinfo>"}),
		makeUnrecognizedRelatorCodesErrorMessage([][]string{
			{"<controlaccess><corpname>Columbia University</corpname></controlaccess>", "orz"},
			{"<controlaccess><corpname>The New School</corpname></controlaccess>", "cpr"},
			{"<controlaccess><famname>Buell Family</famname></controlaccess>", "cpo"},
			{"<controlaccess><famname>Lanier Family</famname></controlaccess>", "fdr"},
			{"<controlaccess><persname>John Doe, 1800-1900</persname></controlaccess>", "clb"},
			{"<controlaccess><persname>Jane Doe, 1800-1900</persname></controlaccess>", "grt"},
			{"<origination><corpname>Queens College</corpname></origination>", "cpr"},
			{"<origination><corpname>Hunter College</corpname></origination>", "orz"},
			{"<origination><famname>Draper family</famname></origination>", "fro"},
			{"<origination><persname>Daisy, Bert</persname></origination>", "clb"},
			{"<origination><persname>Orchid, Ella</persname></origination>", "grt"},
			{"<repository><corpname>NYU Archives</corpname></repository>", "grt"},
		}),
	}

	doTest(invalidEadDataFixturePath, expected, t)
}

func TestValidateEADInvalidXML(t *testing.T) {
	var expected = []string{
		makeInvalidXMLErrorMessage(),
	}

	doTest(invalidXMLFixturePath, expected, t)
}

func TestValidateEADMissingRequiredElements(t *testing.T) {
	var expected = []string{
		makeMissingRequiredElementErrorMessage("<eadid>"),
		makeMissingRequiredElementErrorMessage("<archdesc>"),
	}

	doTest(missingRequiredElementsEADIDAndArchDescFixturePath, expected, t)

	expected = []string{
		makeMissingRequiredElementErrorMessage("<eadid>"),
		makeMissingRequiredElementErrorMessage("<archdesc>/<did>/<repository>"),
	}

	doTest(missingRequiredElementsEADIDAndRepositoryFixturePath, expected, t)

	expected = []string{
		makeMissingRequiredElementErrorMessage("<eadid>"),
		makeMissingRequiredElementErrorMessage("<archdesc>/<did>/<repository>/<corpname>"),
	}

	doTest(missingRequiredElementsEADIDAndRepositoryCorpnameFixturePath, expected, t)
}

func TestValidateEADValidEADNoErrors(t *testing.T) {
	doTest(validEADFixturePath, []string{}, t)
}

func doTest(file string, expected []string, t *testing.T) {
	var validationErrors, err = ValidateEAD(getEADXML(file))
	if err != nil {
		t.Fatalf(fmt.Sprintf(`Unexpected runtime error: %s`, err))
	}

	if len(validationErrors) != len(expected) {
		var message = getNumErrorsMismatchErrorMessage(expected, validationErrors)
		t.Fatalf(message)
	}

	for idx, err := range validationErrors {
		if err != expected[idx] {
			t.Errorf(`Expected error %d to be "%s", got "%s"`, idx, expected[idx], err)
		}
	}
}

func getEADXML(filepath string) []byte {
	EADXML, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	return EADXML
}

func getNumErrorsMismatchErrorMessage(expected []string, errors []string) string {
	return fmt.Sprintf(`Expected %d error(s):

%s

Got %d error(s):

%s`,
		len(expected),
		strings.Join(expected, "\n"),
		len(errors),
		strings.Join(errors, "\n"),
	)
}
