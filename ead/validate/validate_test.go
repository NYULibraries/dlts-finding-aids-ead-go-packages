package validate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var fixturesDirPath string
var invalidEadDataFixturePath string
var invalidXMLFixturePath string
var invalidEADFixturePath string
var missingRequiredElementsEADIDAndArchDescFixturePath string
var missingRequiredElementsEADIDAndRepositoryFixturePath string
var missingRequiredElementsEADIDAndRepositoryCorpnameFixturePath string
var validEADFixturePath string
var validEADWithEADIDLeadingAndTrailingWhitespaceFixturePath string

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
	invalidEADFixturePath = filepath.Join(fixturesDirPath, "ad_rg_009_3_2_1.xml")
	missingRequiredElementsEADIDAndArchDescFixturePath = filepath.Join(fixturesDirPath, "mc_100-missing-eadid-and-archdesc.xml")
	missingRequiredElementsEADIDAndRepositoryFixturePath = filepath.Join(fixturesDirPath, "mc_100-missing-eadid-and-repository.xml")
	missingRequiredElementsEADIDAndRepositoryCorpnameFixturePath = filepath.Join(fixturesDirPath, "mc_100-missing-eadid-and-repository-corpname.xml")
	validEADFixturePath = filepath.Join(fixturesDirPath, "mc_100.xml")
	validEADWithEADIDLeadingAndTrailingWhitespaceFixturePath = filepath.Join(fixturesDirPath, "mc_100-valid-eadid-with-leading-and-trailing-spaces.xml")
}

func TestValidEADIDRegexpString(t *testing.T) {
	validEADIDRegexp, err := regexp.Compile(ValidEADIDRegexpString)
	if err != nil {
		t.Fatalf(
			"regexp.Compile error for ValidEADIDRegexpString \"%s\": %s",
			ValidEADIDRegexpString,
			err,
		)
	}

	var validEADIDs = []string{
		"a_b",
		"a_b_c",
		"a_b_c_d",
		"a_b_c_d_e",
		"a_b_c_d_e_f",
		"a_b_c_d_e_f_g",
		"a_b_c_d_e_f_g_h",
		"0_1_2_3_4_5_6_7",
		"abcdefghijklmnopqrstuvwzyz_0123456789_abcdefghijklmnopqrstuvwzyz_0123456789_abcdefghijklmnopqrstuvwzyz_0123456789_abcdefghijklmnopqrstuvwzyz_0123456789",
		"1_abcdefghijklmnopqrstuvwzyz_0123456789_a",
		"mss_417",
		"photos_220",
		"rg_2_2_7",
		"rg_38_0_1_2",
	}

	var invalidEADIDs = []string{
		"a",
		"1",
		"a1",
		"1a",
		"abcdefghijklmnopqrstuvwzyz",
		"0123456789",
		"_",
		"A",
		"à",
		"â",
		"é",
		"è",
		"ê",
		"ë",
		"î",
		"ô",
		"中文",
		"руссиан",
		"aA",
		"A1",
		"1A",
		"A_B_C_D_E_F_G_H",
		"a_",
		"a_b_",
		"_a",
		"_a_b",
		"a__b__c",
		"a-b",
		"a.b",
		"a_b_c_d_e_f_g_h_i",
		"a b c",
		"a,b,c",
		"a,b,c_abc",
		"a|b|c",
		"a|b|c_abc",
		"a&b&c",
		"a&b&c_abc",
		"a-b-c",
		"a-b-c_abc",
		"a.b.c",
		"a.b.c_abc",
		"",
		"mss.417",
		"PHOTOS_220",
		"rg-2-2-7",
		"Rg_38_0_1_2",
	}

	for _, eadid := range validEADIDs {
		match := validEADIDRegexp.Match([]byte(eadid))
		if !match {
			t.Errorf("regexp fails to match valid <eadid> \"%s\"", eadid)
		}
	}

	for _, eadid := range invalidEADIDs {
		match := validEADIDRegexp.Match([]byte(eadid))
		if match {
			t.Errorf("regexp incorrectly matches invalid <eadid> \"%s\"", eadid)
		}
	}
}

func TestValidateEADInvalidData(t *testing.T) {
	var expected = []string{
		makeInvalidEADIDErrorMessage("mc.100", []rune{'.'}),
		makeInvalidRepositoryErrorMessage("NYU Archives"),
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
		"Unable to parse XML file",
		"failed to create parse input: failed to read document from memory: Entity: line 1: parser error : Start tag expected, '<' not found",
	}

	doTest(invalidXMLFixturePath, expected, t)
}

func TestValidateEADInvalidEAD(t *testing.T) {
	var expected = []string{
		makeInvalidXMLErrorMessage(),
        "schema validation failed",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
        "Element '{urn:isbn:1-931666-22-9}container', attribute 'type': '' is not a valid value of the atomic type 'xs:NMTOKEN'.",
	}

	doTest(invalidEADFixturePath, expected, t)
}

func TestValidateEADMissingRequiredElements(t *testing.T) {
	var expected = []string{
		makeInvalidXMLErrorMessage(),
		"schema validation failed",
		"Element '{urn:isbn:1-931666-22-9}filedesc': This element is not expected. Expected is ( {urn:isbn:1-931666-22-9}eadid ).",
		"Element '{urn:isbn:1-931666-22-9}ead': Missing child element(s). Expected is one of ( {urn:isbn:1-931666-22-9}frontmatter, {urn:isbn:1-931666-22-9}archdesc ).",
	}

	doTest(missingRequiredElementsEADIDAndArchDescFixturePath, expected, t)

	expected = []string{
		makeInvalidXMLErrorMessage(),
		"schema validation failed",
		"Element '{urn:isbn:1-931666-22-9}filedesc': This element is not expected. Expected is ( {urn:isbn:1-931666-22-9}eadid ).",
	}

	doTest(missingRequiredElementsEADIDAndRepositoryFixturePath, expected, t)

	expected = []string{
		makeInvalidXMLErrorMessage(),
		"schema validation failed",
		"Element '{urn:isbn:1-931666-22-9}filedesc': This element is not expected. Expected is ( {urn:isbn:1-931666-22-9}eadid ).",
	}

	doTest(missingRequiredElementsEADIDAndRepositoryCorpnameFixturePath, expected, t)
}

func TestValidateEADInvalidEADActual(t *testing.T) {
	doTest(validEADFixturePath, []string{}, t)
	doTest(validEADWithEADIDLeadingAndTrailingWhitespaceFixturePath, []string{}, t)
}

func TestValidateEADValidEADNoErrors(t *testing.T) {
	doTest(validEADFixturePath, []string{}, t)
	doTest(validEADWithEADIDLeadingAndTrailingWhitespaceFixturePath, []string{}, t)
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
