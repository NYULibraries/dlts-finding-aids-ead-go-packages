package validate

import (
	"fmt"
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
var invalidEADTrailingSpaceInEADIDFixturePath string
var invalidEADLeadingSpaceInEADIDFixturePath string
var invalidEADWithEADIDLeadingAndTrailingWhitespaceFixturePath string
var invalidEADSpaceOnlyInEADIDFixturePath string
var validEADFixturePath string
var akkasahRepositoryNameFixturePath string
var eadIDTooLongFixturePath string
var invalidArchDescLevelFixturePath string
var tooBigFileFixturePath string
var invalidEADWithNamespaceErrorsFixturePath string
var cbhValidEADFixturePath string
var bcValidEADFixturePath string
var bhsValidEADFixturePath string
var eadExportedWithASpacePluginFixturePath string
var arabartarchiveValidEADFixturePath string

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
	invalidEADWithNamespaceErrorsFixturePath = filepath.Join(fixturesDirPath, "bcms_0011-namespace-errors.xml")
	missingRequiredElementsEADIDAndArchDescFixturePath = filepath.Join(fixturesDirPath, "mc_100-missing-eadid-and-archdesc.xml")
	missingRequiredElementsEADIDAndRepositoryFixturePath = filepath.Join(fixturesDirPath, "mc_100-missing-eadid-and-repository.xml")
	missingRequiredElementsEADIDAndRepositoryCorpnameFixturePath = filepath.Join(fixturesDirPath, "mc_100-missing-eadid-and-repository-corpname.xml")
	invalidEADLeadingSpaceInEADIDFixturePath = filepath.Join(fixturesDirPath, "mc_100-leading-space-in-eadid.xml")
	invalidEADTrailingSpaceInEADIDFixturePath = filepath.Join(fixturesDirPath, "mc_100-trailing-space-in-eadid.xml")
	invalidEADSpaceOnlyInEADIDFixturePath = filepath.Join(fixturesDirPath, "space-only-eadid.xml")
	invalidEADWithEADIDLeadingAndTrailingWhitespaceFixturePath = filepath.Join(fixturesDirPath, "mc_100-invalid-eadid-with-leading-and-trailing-spaces.xml")
	validEADFixturePath = filepath.Join(fixturesDirPath, "mc_100.xml")
	akkasahRepositoryNameFixturePath = filepath.Join(fixturesDirPath, "ad_mc_030_ref160-corrected-archdesc-level.xml")
	eadIDTooLongFixturePath = filepath.Join(fixturesDirPath, "tam_647-eadid-too-long.xml")
	invalidArchDescLevelFixturePath = filepath.Join(fixturesDirPath, "ad_mc_030_ref160-invalid-archdesc-level.xml")
	tooBigFileFixturePath = filepath.Join(fixturesDirPath, "b44d567a-95c1-4f0d-b16a-d9448cde1aa5.xml")
	cbhValidEADFixturePath = filepath.Join(fixturesDirPath, "cbhm_0002.xml")
	bcValidEADFixturePath = filepath.Join(fixturesDirPath, "bcms_0001.xml")
	bhsValidEADFixturePath = filepath.Join(fixturesDirPath, "arc_061_meeker.xml")
	eadExportedWithASpacePluginFixturePath = filepath.Join(fixturesDirPath, "mc_1.xml")
	arabartarchiveValidEADFixturePath = filepath.Join(fixturesDirPath, "arabartarchive-ad_mc_091.xml")
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

func doTestWithValidateEADFromFilePath(file string, expected []string, t *testing.T) {
	var validationErrors, err = ValidateEADFromFilePath(file)
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
	EADXML, err := os.ReadFile(filepath)
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
		"a_b_c_d_e_f_g_h_i",
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
		"mc_100 ",
		" mc_100",
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
	}

	doTest(invalidEadDataFixturePath, expected, t)
}

func TestValidateEADInvalidXML(t *testing.T) {
	var expected = []string{
		makeInvalidXMLErrorMessage(),
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

func TestValidateEADInvalidEADNamespaceErrors(t *testing.T) {
	var expected = []string{
		makeInvalidXMLErrorMessage(),
		"schema validation failed",
		"Element '{urn:isbn:1-931666-22-9}extref', attribute 'ns2:href': The attribute 'ns2:href' is not allowed.",
		"Element '{urn:isbn:1-931666-22-9}extref', attribute 'ns2:href': The attribute 'ns2:href' is not allowed.",
		"Element '{urn:isbn:1-931666-22-9}extref', attribute 'ns2:href': The attribute 'ns2:href' is not allowed.",
		"Element '{urn:isbn:1-931666-22-9}extref', attribute 'ns2:href': The attribute 'ns2:href' is not allowed.",
		"Element '{urn:isbn:1-931666-22-9}extref', attribute 'ns2:href': The attribute 'ns2:href' is not allowed.",
		"Element '{urn:isbn:1-931666-22-9}extref', attribute 'ns2:href': The attribute 'ns2:href' is not allowed.",
		"Element '{urn:isbn:1-931666-22-9}extref', attribute 'ns2:href': The attribute 'ns2:href' is not allowed.",
		"Element '{urn:isbn:1-931666-22-9}extref', attribute 'ns2:href': The attribute 'ns2:href' is not allowed.",
		"Element '{urn:isbn:1-931666-22-9}extref', attribute 'ns2:href': The attribute 'ns2:href' is not allowed.",
		"Element '{urn:isbn:1-931666-22-9}extref', attribute 'ns2:href': The attribute 'ns2:href' is not allowed.",
		"Element '{urn:isbn:1-931666-22-9}extref', attribute 'ns2:href': The attribute 'ns2:href' is not allowed.",
	}

	doTest(invalidEADWithNamespaceErrorsFixturePath, expected, t)
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

func TestValidateEADIDTooLong(t *testing.T) {
	expected := []string{
		makeEADIDTooLongErrorMessage("iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii_iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii"),
	}

	doTest(eadIDTooLongFixturePath, expected, t)
}
func TestValidateArchDescLevel(t *testing.T) {
	expected := []string{
		makeInvalidArchDescLevelErrorMessage("series"),
	}

	doTest(invalidArchDescLevelFixturePath, expected, t)
}

func TestValidateEADValidEADNoErrors(t *testing.T) {
	doTest(validEADFixturePath, []string{}, t)
}

func TestValidateEADAkkasahTitleEADNoErrors(t *testing.T) {
	doTest(akkasahRepositoryNameFixturePath, []string{}, t)
}

func TestAssertMaxFileSize(t *testing.T) {
	f, err := os.Create(tooBigFileFixturePath)
	if err != nil {
		t.Error("problem creating test file")
	}

	// https://stackoverflow.com/a/16806474
	if err := f.Truncate(MAXIMUM_FILE_SIZE + 1); err != nil {
		t.Error("problem truncating test file")
	}
	defer os.Remove(tooBigFileFixturePath)

	expected := []string{
		makeFileTooBigErrorMessage(tooBigFileFixturePath, MAXIMUM_FILE_SIZE+1),
	}
	doTestWithValidateEADFromFilePath(tooBigFileFixturePath, expected, t)
}

func TestValidateEADValidCBHCenterForBrooklynHistory(t *testing.T) {
	doTest(cbhValidEADFixturePath, []string{}, t)
}

func TestValidateEADValidCBHBrooklynCollection(t *testing.T) {
	doTest(bcValidEADFixturePath, []string{}, t)
}

func TestValidateEADValidCBHBrooklynHistoricalSociety(t *testing.T) {
	doTest(bhsValidEADFixturePath, []string{}, t)
}

func TestValidateEADIDWithTrailingSpace(t *testing.T) {
	expected := []string{
		makeInvalidEADIDErrorMessage("mc_100 ", []rune{' '}),
	}

	doTest(invalidEADTrailingSpaceInEADIDFixturePath, expected, t)
}

func TestValidateEADIDWithLeadingAndTrailingSpace(t *testing.T) {
	invalidRunes := []rune("\n ")
	expected := []string{
		makeInvalidEADIDErrorMessage(" mc_100\n ", invalidRunes),
	}

	doTest(invalidEADWithEADIDLeadingAndTrailingWhitespaceFixturePath, expected, t)
}

func TestValidateEADIDWithLeadingSpace(t *testing.T) {
	expected := []string{
		makeInvalidEADIDErrorMessage(" mc_100", []rune{' '}),
	}

	doTest(invalidEADLeadingSpaceInEADIDFixturePath, expected, t)
}

func TestValidateEADIDWithSpaceOnly(t *testing.T) {
	expected := []string{makeMissingRequiredElementErrorMessage("<eadid>")}

	doTest(invalidEADSpaceOnlyInEADIDFixturePath, expected, t)
}

func TestValidateEADExportedWithASpacePlugin(t *testing.T) {
	expected := []string{makeExportedWithEADPluginErrorMessage()}

	doTest(eadExportedWithASpacePluginFixturePath, expected, t)
}

func TestValidateEADValidArabArtArchivee(t *testing.T) {
        doTest(arabartarchiveValidEADFixturePath, []string{}, t)
}
